package snapshot

import (
	"bytes"
	"image"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"

	_ "image/jpeg"
)

func Capture() (image.Image, error) {
	dev, err := device.Open("/dev/video0",
		device.WithBufferSize(1),
		device.WithPixFormat(v4l2.PixFormat{
			Width:  640,
			Height: 480,
			// PixelFormat: v4l2.PixelFmtYUYV,
			PixelFormat: v4l2.PixelFmtMJPEG, // hopefully supported
		}),
	)
	if err != nil {
		return nil, err
	}
	defer dev.Close()

	frame := <-dev.GetFrames()

	var buf bytes.Buffer
	_, err = buf.Read(frame.Data)
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(&buf)
	return i, err
}
