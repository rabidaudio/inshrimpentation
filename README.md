Shrimp tank monitoring

1. Hardware interface for measuring data
2. Web server for viewing data and configuring
3. sqlite db for storing data


https://pkg.go.dev/github.com/blackjack/webcam
https://periph.io/project/design/
Live stream is handled out-of-band using mediamtx
https://mediamtx.org/docs/read/web-browsers#embed-in-a-website

build raspberry pi image using cloud-init and firstrun.sh
https://www.raspberrypi.com/news/cloud-init-on-raspberry-pi-os/


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
