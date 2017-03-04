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

    "github.com/influxdata/influxdb/client/v2"

    "github.com/faryon93/hlswatch/state"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    // influx measurement description
    POINT_PRECISION    = "s"
    STREAM_MEASUREMENT = "streams"
    TAG_NODE = "node"
    TAG_STREAM = "stream"
    VALUE_VIEWER = "viewer"

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
    })
    if err != nil {
        log.Println("[influx] failed to connect to influxdb:", err.Error())
        os.Exit(-1)
    }

    // some data needed for the measurements
    hostname, err := os.Hostname()
    if err != nil {
        log.Println("[influx] failed to query hostname:", err.Error())
        return
    }
    viewerTimeout := time.Duration(ctx.Conf.Common.ViewerTimeout) * time.Second

    for {
        // A new batch of points
        bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
            Database: ctx.Conf.Influx.Database,
            Precision: POINT_PRECISION,
        })

        startTime := time.Now()

        // save the numbers for all streams to influx db
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

        // write the datapoints to influx
        err = influx.Write(bp)
        if err != nil {
            log.Println("[influx] failed to write datapoint:", err.Error())
        }

        time.Sleep(INFLUX_CYCLE_TIME)
    }
}