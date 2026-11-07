Shrimp tank monitoring

1. Hardware interface for measuring data
2. Web server for viewing data and configuring
3. sqlite db for storing data


https://pkg.go.dev/github.com/blackjack/webcam
https://periph.io/project/design/


db:
    sqlite db for storing measurements, configs
instruments/*:
    interface for measuring
scheduler:
    at certain times, queries instruments, saves to db. if outside levels, send alert
server:
    api for accessing data
live:
    live webcam stream
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


Virtual webcam in linux:
v4l2loopback
ffmpeg -loop 1 -re -i ~/Downloads/Lenna_\(test_image\).png -f v4l2 -vcodec rawvideo -pix_fmt yuv420p /dev/video0
sudo modprobe v4l2loopback devices=1 max_buffers=2 exclusive_caps=1 card_label=VirtualCam
