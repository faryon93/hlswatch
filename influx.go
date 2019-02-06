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
    "os"
    "time"

    "github.com/influxdata/influxdb1-client/v2"

    "github.com/faryon93/hlswatch/state"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    // influx connection options
    INFLUX_TIMEOUT      = 500 * time.Millisecond
    INFLUX_FAIL_COUNT   = 15

    // influx measurement description
    POINT_PRECISION    = "s"
    STREAM_MEASUREMENT = "streams"
    TAG_NODE = "node"
    TAG_STREAM = "stream"
    VALUE_VIEWER = "viewers"

    // task options
    INFLUX_CYCLE_TIME = 1 * time.Second
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func InfluxMetrics(ctx *state.State) {
    // connect to influxdb
    influx, err := client.NewHTTPClient(client.HTTPConfig{
        Addr: ctx.Conf.Influx.Address,
        Username: ctx.Conf.Influx.User,
        Password: ctx.Conf.Influx.Password,
        Timeout: INFLUX_TIMEOUT,
    })
    if err != nil {
        log.Println("[influx] failed to connect to influxdb:", err.Error())
        return
    }

    // check the connectivity
    _, _, err = influx.Ping(INFLUX_TIMEOUT)
    if err != nil {
        log.Println("[influx] failed to connect to influxdb:", err.Error())
        return
    }

    // some data needed for the measurements
    hostname, err := os.Hostname()
    if err != nil {
        log.Println("[influx] failed to query hostname:", err.Error())
        return
    }
    viewerTimeout := time.Duration(ctx.Conf.Common.ViewerTimeout) * time.Second

    // variables for runtime
    failCount := 0

    for {
        // A new batch of points
        bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
            Database: ctx.Conf.Influx.Database,
            Precision: POINT_PRECISION,
        })

        startTime := time.Now()

        // save the numbers for all streams to influx db
        ctx.StreamsMutex.Lock()
        for streamName, stream := range ctx.Streams {
            count := stream.GetCurrentViewers(viewerTimeout)
            pt, _ := client.NewPoint(
                STREAM_MEASUREMENT,
                map[string]string{
                    TAG_NODE: hostname,
                    TAG_STREAM: streamName,
                },
                map[string]interface{}{
                    VALUE_VIEWER: count,
                },
                startTime,
            )
            bp.AddPoint(pt)
        }
        ctx.StreamsMutex.Unlock()

        // write the datapoints to influx
        if len(bp.Points()) > 0 {
            err = influx.Write(bp)
            if err != nil {
                log.Println("[influx] failed to write datapoint:", err.Error())

                // check the fail count and disable this module if necessary
                failCount++
                if failCount >= INFLUX_FAIL_COUNT {
                    log.Println("[influx] reached fail count, disabling module")
                    return
                }

            // reset the failcount
            } else {
                failCount = 0
            }
        }

        time.Sleep(INFLUX_CYCLE_TIME)
    }
}
