package main
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

// ----------------------------------------------------------------------------
//  variables
// ----------------------------------------------------------------------------

var (
    // release information
    APP_NAME = "hlswatch"
    APP_VERSION = "0.1.1"

    // filled by build tool
    GIT_COMMIT string
    BUILD_TIME string
    BUILD_NUMBER string
)


// ----------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------

func GetAppIdentifier() (string) {
    str := APP_NAME + " " + APP_VERSION
    if len(BUILD_NUMBER) > 0 {
        str += "-" + BUILD_NUMBER
    }

    if len(GIT_COMMIT) > 0 {
        str += " (#" + GIT_COMMIT + ")"
    }

    if len(BUILD_TIME) > 0 {
        str += " " + BUILD_TIME
    }

    return str
}
