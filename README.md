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

## NGINX Setup
Because this software is responsible for delivering all data to the client it is not necessary to serve the HLS fragments via nginx to the public. If you want all NGINX features like access control, compression, ... you can reverse proxy incoming requests by NGINX to hlswatch.
This software relies on some configuration option the nginx-rtmp-module offers. The settings `hls_cleanup` and `hls_nested` need to be enabled:

```
rtmp {
    server {
        listen 1935;
        chunk_size 4000;

        application live {
            live on;
            hls on;
            hls_fragment_naming system;
            hls_fragment 5s;
            hls_path /tmp/hls;
            hls_nested on;
        }
    }
}
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
