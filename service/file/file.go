package file

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/qwantix/qxeye/camera"
	"github.com/qwantix/qxeye/config"
)

type FileService struct {
	config config.ServiceConfig
	dir    string
	trace  camera.CameraTraceFlag
}

func New(cfg config.ServiceConfig) *FileService {
	fs := new(FileService)
	fs.dir = cfg.String("dir")
	if _, err := os.Stat(fs.dir); os.IsNotExist(err) {
		os.MkdirAll(fs.dir, os.ModeTemporary)
	}
	camera.CameraTraceFlagFromString(cfg.String("trace"))
	return fs
}

func (fs *FileService) CaptureToFile(ev *camera.CameraEvent, t *config.TriggerConfig) string {
	var matches []camera.CameraMatch
	for _, m := range ev.Matches {
		if m.Label == t.On && m.Confidence >= t.Confidence {
			matches = append(matches, m)
		}
	}
	filename := fmt.Sprint("capture", time.Now().UnixNano(), ".jpg")
	filename = path.Join(fs.dir, filename)
	ev.Capture(filename, matches, fs.trace)
	return filename
}

// func RecordToFile() {

// }

func (fs *FileService) RemoveFile(filename string) {
	// TODO ensure is file is in dir
	os.Remove(filename)
}

func (fs *FileService) Push(ev *camera.CameraEvent, t *config.TriggerConfig) {
	fs.CaptureToFile(ev, t)
}
