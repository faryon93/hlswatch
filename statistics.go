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


// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import (
    "log"
    "time"

    "github.com/faryon93/hlswatch/state"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func StatisticsTask(ctx *state.State) {
    lastCount := 0

    for {
        // echo the concurrent viewer count when it has changed only
        count := ctx.GetStream("futurama").GetCurrentViewers(viewerTimeout)
        if count != lastCount {
            log.Print("[futurama] concurrent viewers: ", count)
            lastCount = count
        }

        // wait for the next computation round
        time.Sleep(cycleTime * time.Second)
    }
}
