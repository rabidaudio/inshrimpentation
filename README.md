Shrimp tank monitoring

1. Hardware interface for measuring data
2. Web server for viewing data and configuring
3. sqlite db for storing data


https://pkg.go.dev/github.com/blackjack/webcam
https://periph.io/project/design/
Live stream is handled out-of-band using mediamtx
https://mediamtx.org/docs/read/web-browsers#embed-in-a-website

Because of mediamtx's non-standard build process, it can't be installed with `go install`,
and therefore can't be added directly to gokrazy. Instead, we have a git submodule that
we compile locally, and add to gokrazy using the local directory. This also allows easy
forks and version updates.

Additionally, it's rPi camera interface isn't statically compiled and thus cannot run on gokrazy.
Instead, we manually load ffmpeg onto the image and use that for accessing the camera.

TODO: could we compile our own simple v4l2 -> rtsp program in go??
Could run ffmpeg in docker via podman
Could try compiling mtxrpi in the gokrazy docker image
Could load rpicam-vid manually

https://github.com/AlexxIT/go2rtc

db:
    sqlite db for storing measurements, configs
instruments/*:
    interface for measuring
scheduler:
    at certain times, queries instruments, saves to db. if outside levels, send alert
server:
    api for accessing data
frontend:
    svelte spi for displaying data


WaterLevel
    measurements: inches
    alarms:
        less than 7.5 inches
        greater than 8.5 inches
Temperature
    measurements: celcius
    alarms:
        less than 20 C (68F)
        more than 23 C (73.5F)
Brightness:
    measurements: (arbitrary - relative measure, range TBD)
