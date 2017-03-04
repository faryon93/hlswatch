# hlswatch - keep track of hls viewer stats
hlswatch is a simple program to keep track of the concurrent viewer count of a [HLS](https://tools.ietf.org/html/draft-pantos-http-live-streaming-20) live stream. This peace of software is intended to be placed in front of a NGINX server with the [nginx-rtmp-module](https://github.com/arut/nginx-rtmp-module) installed and HLS encoding enabled.
The NGINX server is used as a transcoding instance to translate an incoming RTMP stream into HLS compatible video fragments. Delivering the fragments and playlists to the clients is handled by hlswatch.

To count the number of concurrent viewers hlswatch monitors the m3u8 playlist accesses per client. If the client does not issue a playlist reload in a certain amount if time it is considered as "not watching anymore".

As database backend [InfluxDB](https://www.influxdata.com/) is supported only.

## Configuration
Per default hlswatch uses ```/etc/hlswatch/hlswatch.conf``` as configuration file. If you want to change this path, just call hlswatch with your configuration file as the first argument.

```
[common]
listen = ":3000"
hls_path = "/tmp/hls/"
viewer_timeout = 15

[influx]
address = "http://localhost:8086"
database = "hlswatch"
user = "hlswatch"
password = "hlswatch"
```

## Tested Players
The application was tested with the following web players:

Player                                     | Working |
-------------------------------------------|---------|
[clappr](https://github.com/clappr/clappr) |    ✔   |

## ToDo
- HTTPS support to ensure SSL termination
- Caching of m3u8 playlist and video fragments in RAM
- REST interface for statistics
