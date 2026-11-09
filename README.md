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


formats: map[842093913:Planar YUV 4:2:0]
sizes: [{512 512 0 512 512 0}]
frame rates: [1/30]


macos:
ffmpeg -f avfoundation -framerate 30 -video_size 640x480 -i "0" out.mp4

linux: -f v4l2 -i /dev/video0



		"-f", "mjpeg",
		"-i", inpath,
		"-c:v", "libvpx-vp9", // -c:v libvpx-vp9: Encode video using VP9
		"-b:v", "800k", // -b:v 800k: Target video bitrate
		"-deadline", "realtime", // -deadline realtime: Optimize encoding speed for live streaming (low latency)
		"-cpu-used", "4", // -cpu-used 4: CPU preset for libvpx-vp9 (range 0-8; higher means faster encoding, lower means better compression)
		"-g", "60", // -g 60: Keyframe interval (every 60 frames)
		"-pix_fmt", "yuv420p", // -pix_fmt yuv420p: Ensure pixel format compatibility for web browsers
		"-an", // -an: Strip audio (MJPEG is video-only; see note below for Opus audio)
		"-y",  // Overwrite output file if it exists




ffmpeg -f avfoundation -framerate 30 -video_size 640x480 -i "0" -b:v 800k -deadline realtime -an -y out.mp4