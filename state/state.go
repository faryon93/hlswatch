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
    "sync"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type State struct {
    Streams map[string]*Stream
    StreamMutex sync.Mutex
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func New() *State {
    return &State{
        Streams: make(map[string]*Stream),
    }
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (s *State) GetStream(name string) (*Stream) {
    s.StreamMutex.Lock()
    stream := s.Streams[name]
    if stream == nil {
        stream = NewStream()
        s.Streams[name] = stream
    }
    s.StreamMutex.Unlock()

    return stream
}
