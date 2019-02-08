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
    "strings"
    "time"
    "errors"
    "net/http"
    "crypto/rand"
    "encoding/base64"
    "log"
    "net"
    "bytes"

    "github.com/faryon93/hlswatch/state"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    PLAYLIST_FILE_EXTENSION = "m3u8"
    TOKEN_URL_PARAMETER     = "token"
)

//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
    start net.IP
    end net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
    // strcmp type byte comparison
    if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
        return true
    }
    return false
}

var privateRanges = []ipRange{
    ipRange{
        start: net.ParseIP("10.0.0.0"),
        end:   net.ParseIP("10.255.255.255"),
    },
    ipRange{
        start: net.ParseIP("100.64.0.0"),
        end:   net.ParseIP("100.127.255.255"),
    },
    ipRange{
        start: net.ParseIP("172.16.0.0"),
        end:   net.ParseIP("172.31.255.255"),
    },
    ipRange{
        start: net.ParseIP("192.0.0.0"),
        end:   net.ParseIP("192.0.0.255"),
    },
    ipRange{
        start: net.ParseIP("192.168.0.0"),
        end:   net.ParseIP("192.168.255.255"),
    },
    ipRange{
        start: net.ParseIP("198.18.0.0"),
        end:   net.ParseIP("198.19.255.255"),
    },
}


// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
    // my use case is only concerned with ipv4 atm
    if ipCheck := ipAddress.To4(); ipCheck != nil {
        // iterate over all our ranges
        for _, r := range privateRanges {
            // check if this ip is in a private range
            if inRange(r, ipAddress){
                return true
            }
        }
    }
    return false
}


func getIPAdress(r *http.Request) string {
    for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
        addresses := strings.Split(r.Header.Get(h), ",")
        // march from right to left until we get a public address
        // that will be the address right before our proxy.
        for i := len(addresses) -1 ; i >= 0; i-- {
            ip := strings.TrimSpace(addresses[i])
            // header can contain spaces too, strip those out.
            realIP := net.ParseIP(ip)
            if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
                // bad address, go to next
                continue
            }
            return ip
        }
    }
    return ""
}


// --------------------------------------------------------------------------------------
//  http handler
// --------------------------------------------------------------------------------------

func Hls(ctx *state.State, h http.Handler) http.Handler {
    f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // headers applicable for all files
        w.Header().Set("Access-Control-Allow-Origin", "*")

        // We keep track of the concurrent users by watching the
        // m3u8 accesses. Each clients gets a unique token assigned
        // wich is used to identify it. In order to get the new video fragments
        // the client is supposed to reload the playlistfile periodically, we use
        // this circumstance to count the viewers based on a timeout value.
        if strings.Contains(r.URL.String(), PLAYLIST_FILE_EXTENSION) {
            // some meta information needed to process the request properly
            streamName, err := getStreamName(r.URL.Path)
            if err != nil {
                log.Println("invalid hls url requested:", r.URL.String())
                http.Error(w, "invalid streaming url", http.StatusNotAcceptable)
                return
            }
            stream := ctx.GetStream(streamName)
            if stream == nil {
                http.Error(w, "invalid hls stream", http.StatusNotFound)
                return
            }

            // get the viewer by its token
            token := r.URL.Query().Get(TOKEN_URL_PARAMETER)
            stream.Lock()
            viewer := stream.Viewers[token]
            stream.Unlock()

            // we do not want caching for the playlist, because it changes
            // everytime a new video fragment is created
            w.Header().Set("Cache-Control", "no-cache")
            w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

            // generate a new token if the client does not supply one
            if token == "" {
                token = nextToken()
                
                actualIP := getIPAdress(r)
                
                stream.Lock()
                stream.Viewers[token] = &state.Viewer{
                    FirstSeen: time.Now(),
                    LastSeen: time.Now(),
                    Ip: actualIP,
                }
                stream.Unlock()

                // assemble the url with the token appended
                query := r.URL.Query()
                query.Set(TOKEN_URL_PARAMETER, token)
                r.URL.RawQuery = query.Encode()

                // redirect the client to the url with the appended token
                http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
                return

                // the client supplied an invalid token -> wo do not handle this request
            } else if viewer == nil {
                http.Error(w, "invalid hls token", http.StatusNotAcceptable)
                return
            }

            // update the lastseen time, so we can keep track of the viewers
            viewer.LastSeen = time.Now()
        }

        // serve the requested file from fs
        h.ServeHTTP(w, r)
    })

    return f
}


// --------------------------------------------------------------------------------------
//  helper functions
// --------------------------------------------------------------------------------------

func getStreamName(url string) (string, error) {
    // the first "subdirectory" is the name of the stream
    parts := strings.SplitN(strings.Trim(url, "/"), "/", 2)
    if len(parts) < 2 {
        return "", errors.New("url to short to get stream name")
    }

    return parts[0], nil
}

func nextToken() string {
    key := make([]byte, 32)
    rand.Read(key)

    return base64.StdEncoding.EncodeToString(key)
}
