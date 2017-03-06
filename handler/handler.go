package handler
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
    "net/http"
    "encoding/json"

    "github.com/faryon93/hlswatch/state"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Handler func(ctx *state.State, w http.ResponseWriter, r *http.Request)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Writes the JSON representation of v to the supplied http.ResposeWriter.
// If an error occours while marshalling the object the http response
// will be an internal server error.
func Jsonify(w http.ResponseWriter, v interface{}) {
    js, err := json.Marshal(v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}
