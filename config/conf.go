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
    "github.com/BurntSushi/toml"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Conf struct {
    Common struct {
        Listen string `toml:"listen"`
        HlsPath string `toml:"hls_path"`
        ViewerTimeout int `toml:"viewer_timeout"`
        SslCertificate string `toml:"ssl_certificate"`
        SslPrivateKey string `toml:"ssl_privatekey"`
    } `toml:"common"`

    Influx struct {
        Address string `toml:"address"`
        User string `toml:"user"`
        Password string `toml:"password"`
        Database string `toml:"database"`
    } `toml:"influx"`
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

    return &conf, nil
}