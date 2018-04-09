package main
// hlswatch
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
    "path/filepath"
    "io/ioutil"

    "github.com/fsnotify/fsnotify"

    "github.com/faryon93/hlswatch/util"
    "github.com/faryon93/hlswatch/state"
)



// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func StreamWatcher(ctx *state.State) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Println("[streamwatcher] failed to create fswatcher:", err.Error())
        return
    }
    defer watcher.Close()

    // monitor the hls data directory
    err = watcher.Add(ctx.Conf.Common.HlsPath)
    if err != nil {
        log.Println("[streamwatcher] failed to configure fs watcher:", err.Error())
        return
    }

    // index the currently available streams in fs
    files, _ := ioutil.ReadDir(ctx.Conf.Common.HlsPath)
    for _, f := range files {
       if f.IsDir() {
           log.Println("[streamwatcher] adding stream \"" + f.Name() + "\"")
           ctx.SetStream(f.Name(), state.NewStream())
       }
    }

    // from now on listen for all changes
    for {
        // wat for the next fs event
        event := <-watcher.Events
        streamName := filepath.Base(event.Name)

        // a new stream is created in fs
        if event.Op == fsnotify.Create && util.IsDir(event.Name) {
            log.Println("[streamwatcher] adding new stream \"" + streamName + "\"")
            ctx.SetStream(streamName, state.NewStream())

        // stream is removed from fs
        } else if event.Op == fsnotify.Remove {
            log.Println("[streamwatcher] removing stream \"" + streamName + "\"")
            ctx.RemoveStream(streamName)
        }
    }
}
