# hlswatch - keep track of hls viewer stats
hlswatch is a simple program to keep track of the concurrent viewer count of a [HLS](https://tools.ietf.org/html/draft-pantos-http-live-streaming-20) live stream. This piece of software is intended to be placed in front of a NGINX server with the [nginx-rtmp-module](https://github.com/arut/nginx-rtmp-module) installed and HLS encoding enabled.
The NGINX server is used as a transcoding instance to translate an incoming RTMP stream into HLS compatible video fragments. Delivering the fragments and playlists to the clients is handled by hlswatch.

To count the number of concurrent viewers hlswatch monitors the m3u8 playlist accesses per client. If the client does not issue a playlist reload in a certain amount of time it is considered as "not watching anymore".

The gathered data can be obtained in two different ways:

1. The viewer counts are written to an [InfluxDB](https://www.influxdata.com/) database every second, if the the ```[influxdb]``` section is present in the config file. The tags ```node``` and ```stream``` are automatically populatet with the hostname and the stream name. The current count of concurrent viewers is stored in the value ```viewers```.

2. A REST endpoint is exposed on the address/port configured in ```common.listen```. A JSON object can be accessed at ```/stats```, which contains the current viewer numbers for each stream.

**Keep in mind that this piece of software hasn't been tested in production!**

## Building
The Go application is built as part of the Docker image build process.
In order to build a fresh copy of *hlswatch* you just have to enter:

    $: docker build -t faryon93/hlswatch .

## Configuration
Per default hlswatch uses ```/etc/hlswatch/hlswatch.conf``` as configuration file. If you want to change this path, just call hlswatch with your configuration file as the first argument.

```
[common]
listen = "0.0.0.0:3000"        # listen port/address for REST interface
hls_path = "/tmp/hls/"         # path where nginx-rtmp stores the HLS fragments
viewer_timeout = 15            # time in seconds when a viewer is considered not watching anymore

[influx]
address = "http://localhost:8086"   # address of the influxdb server
database = "hlswatch"               # database name
user = "hlswatch"                   # database user
password = "hlswatch"               # database password
```

Some configuration parameters can be overriden by environment variables. See ```config/config.go``` for valid variable names. Note: Not all parameters can be replaced by environment variables.

## NGINX Setup
Because this software is responsible for delivering all data to the client it is not necessary to serve the HLS fragments via nginx to the public. If you want all NGINX features like access control, compression, TLS Termination, ... you can reverse proxy incoming requests by NGINX to hlswatch.
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
[clappr](https://github.com/clappr/clappr) |    ✔    |

## Docker
This repository contains a Dockerfile, which builds a container which contains an NGINX webserver compiled with the nginx-rtmp-module.
The current version contained in the master branch is automatically pushed to [Docker Hub](https://hub.docker.com/r/faryon93/hlswatch/).
Some configuration files to enable live streaming via RTMP and pass all HTTP requests to hlswatch is included by default.
For production use you should consider adding SSL termination in NGINX and secure the access to hlswatchs statistics page.

Running the container:
```
$: docker run --rm -t -i \
              --name nginx-hls \
              -p 1935:1935 \
              -p 80:80 \
              -e HLS_INFLUX_ADDR=http://localhost:8086 \
              -e HLS_INFLUX_DB=hlswatch \
              -e HLS_INFLUX_USER=hlswatch \
              -e HLS_INFLUX_PASSWORD=hlswatch \
              faryon93/hlswatch
```

## ToDo
- Caching of m3u8 playlist and video fragments in RAM
- Disable directory listing in hlswatch
- Execute hlswatch as unprivileged user (proper process supervision)
