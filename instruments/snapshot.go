package instruments

import (
	"fmt"
	"image/jpeg"
	"os"
	"time"

	"github.com/rabidaudio/inshrimpentation/snapshot"
)

// Takes a photo and saves it to local storage.
// The "measurement" returned is the timestamp in seconds
// which is also the name of the file
func Snapshot() (float64, error) {
	t := time.Now().Unix()
	img, err := snapshot.Capture()
	if err != nil {
		return 0, err
	}

	filename := fmt.Sprintf("storage/snapshots/%d.jpg", t)
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	err = jpeg.Encode(f, img, nil)
	if err != nil {
		return 0, err
	}

	return float64(t), nil
}
