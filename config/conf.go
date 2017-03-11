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
    "reflect"

    "github.com/BurntSushi/toml"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    ENV_PREFIX = "HLS"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Conf struct {
    // common configuration options
    Common struct {
        Listen string `toml:"listen"`
        SslCertificate string `toml:"ssl_certificate"`
        SslPrivateKey string `toml:"ssl_privatekey"`
        HlsPath string `toml:"hls_path"`
        ViewerTimeout int `toml:"viewer_timeout"`
    } `toml:"common"`

    // influx database
    Influx struct {
        Address string `toml:"address" env:"ADDR"`          // HLS_INFLUX_ADDR
        User string `toml:"user" env:"USER"`                // HLS_INFLUX_USER
        Password string `toml:"password" env:"PASSWORD"`    // HLS_INFLUX_PASSWORD
        Database string `toml:"database" env:"DB"`          // HLS_INFLUX_DB
    } `toml:"influx" env:"INFLUX"`
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func Load(path string) (*Conf, error) {
    // decode the conf file to struct
    var conf Conf
    if _, err := toml.DecodeFile(path, &conf); err != nil {
        return nil, err
    }

    // apply all configuration overrides from environment variables
    err := applyEnvOverrides(ENV_PREFIX, reflect.ValueOf(&conf))
    if err != nil {
        return nil, err
    }

    return &conf, nil
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (c *Conf) IsSslEnabled() (bool) {
    return len(c.Common.SslCertificate) > 0 &&
           len(c.Common.SslPrivateKey) > 0
}

