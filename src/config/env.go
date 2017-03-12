package config
// hlswatch - keep track of hls viewer stats
// Copyright (C) 2017 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import (
    "os"
    "reflect"
    "strconv"
    "errors"
)


// --------------------------------------------------------------------------------------
//  private members
// --------------------------------------------------------------------------------------

func applyEnvOverrides(prefix string, obj reflect.Value) (error) {
    // dereference the pointer if necessary
    if obj.Kind() == reflect.Ptr {
        obj = obj.Elem()
    }

    // we only want structs to parse
    if obj.Kind() != reflect.Struct {
        return errors.New("only structs are allowed")
    }

    for fieldIdx := 0; fieldIdx < obj.NumField(); fieldIdx++ {
        // some metadata for the currently inspected
        // field of the struct
        fieldType := obj.Type().Field(fieldIdx)
        fieldValue := obj.Field(fieldIdx)
        fieldTag := fieldType.Tag.Get("env")

        // we are only intrested in fields with "env" tags
        if len(fieldTag) > 0 {
            envVarName := fieldTag
            if len(prefix) > 0 {
                envVarName = prefix + "_" + envVarName
            }

            // another nested struct -> keep on interating
            if fieldType.Type.Kind() == reflect.Struct {
                applyEnvOverrides(envVarName, fieldValue)

            // finally we reached a real fieldValue type -> check env var
            // and override the fieldValue if needed
            } else {
                // an environment variable is specified by the user
                // so try to parse it and set the apropriate fieldValue
                overrideValue := os.Getenv(envVarName)
                if len(overrideValue) > 0 {

                    // String
                    if fieldType.Type.Kind() == reflect.String {
                        fieldValue.SetString(overrideValue)

                    // Integer
                    } else if fieldType.Type.Kind() == reflect.Int {
                        newValue, err := strconv.Atoi(overrideValue)
                        if err != nil {
                            continue
                        }

                        fieldValue.SetInt(int64(newValue))

                    // Boolean
                    } else if fieldType.Type.Kind() == reflect.Bool {
                        newValue := false
                        if overrideValue == "1" || overrideValue == "true" {
                            newValue = true
                        }

                        fieldValue.SetBool(newValue)
                    }
                }
            }
        }
    }

    return nil
}
