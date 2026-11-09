package live

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Stream(ctx context.Context) (io.ReadCloser, func() error, error) {
	path := "/tmp/stream.mp4"

	if err := syscall.Mkfifo(path, 0666); err != nil && !os.IsExist(err) {
		return nil, nil, err
	}
	// defer os.Remove(path)

	args := []string{
		"-f", "avfoundation",
		"-framerate", "30",
		"-video_size", "640x480",
		"-i", "0",
		"-pix_fmt", "uyvy422",
		"-b:v", "800k",
		"-an", "-y",
		path,
	}

	log.Printf("ffmpeg %v\n", strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	go func() { io.Copy(os.Stdout, stderrPipe) }()

	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return f, cmd.Wait, nil

	// conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Net: "unix", Name: "/tmp/stream.mp4"})
	// if err != nil {
	// 	return nil, err
	// }

	// return conn, nil
}
