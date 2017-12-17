package main

import (
	"fmt"
	"time"

	"github.com/qwantix/qxeye/camera"
	"github.com/qwantix/qxeye/config"
	"github.com/qwantix/qxeye/service/file"
	"github.com/qwantix/qxeye/service/http"
	"github.com/qwantix/qxeye/service/twitter"
	"github.com/qwantix/qxeye/util"
)

type Service interface {
	Push(ev *camera.CameraEvent, t *config.TriggerConfig)
}

var services map[string]Service = make(map[string]Service)
var notifyDelayCache map[string]time.Time = make(map[string]time.Time)

func canNotify(ev *camera.CameraEvent, t *config.TriggerConfig) bool {
	for _, m := range ev.Matches {
		if m.Label == t.On && m.Confidence >= t.Confidence {
			if len(m.Zones) == 0 {
				return true
			}
			for _, z := range m.Zones {
				if t.HasZone(z.Name) && z.Coverage > 0.1 {
					return true
				}
			}
		}
	}
	return false
}

func onEvent(qe *QxEye, ev camera.CameraEvent) {
	defer ev.Release()
	key := func(camera string, idx int) string {
		return fmt.Sprintf("%s:%d", camera, idx)
	}
	applyDelay := func(camera string, idx, delay int) {
		if delay > 0 {
			notifyDelayCache[key(camera, idx)] = time.Now().Add(time.Second * time.Duration(int(delay)))
		}
	}

	now := time.Now()
	for i, t := range qe.config.Triggers {
		if t.Delay > 0 && notifyDelayCache[key(ev.CameraName, i)].After(now) {
			util.Warn("Under delay")
			continue
		}
		if !canNotify(&ev, &t) {
			continue
		}

		if services[t.Service] == nil {
			cfg := qe.config.Services[t.Service]
			var s Service
			switch cfg.Service {
			case "file":
				s = file.New(cfg)
			case "http":
				s = http.New(cfg)
			case "script": // TODO
			case "twitter":
				s = twitter.New(cfg)
			default:
				util.Error("Notifier service '", cfg.Service, "' at '", t.Service, "' not found")
				continue
			}
			services[t.Service] = s
		}
		if services[t.Service] != nil {
			services[t.Service].Push(&ev, &t)
			applyDelay(ev.CameraName, i, t.Delay)
		}
	}
}
