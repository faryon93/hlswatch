package state
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
    "time"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Stream struct {
    Viewers map[string]*Viewer
    StartTime time.Time
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func NewStream() *Stream {
    return &Stream{
        Viewers: make(map[string]*Viewer),
        StartTime: time.Now(),
    }
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (s *Stream) GetCurrentViewers(timeout time.Duration) int {
    count := 0
    for _, viewer := range s.Viewers {
        // only count those who did not time out already
        if viewer.LastSeen.Add(timeout).After(time.Now()) {
            count++
        }
    }

    return count
}

func (s *Stream) GetUptime() time.Duration {
    return time.Now().Sub(s.StartTime)
}
