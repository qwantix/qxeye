package main

import (
	"flag"
	"time"

	"github.com/qwantix/qxeye/camera"
	"github.com/qwantix/qxeye/config"
	"github.com/qwantix/qxeye/util"
)

type QxEye struct {
	config           config.Config
	notifyDelayCache map[string]time.Time
	cams             []*camera.Camera
}

func New() *QxEye {
	qe := new(QxEye)
	util.Log("QxEye: create")
	qe.init()
	return qe
}

func (qe *QxEye) init() {
	configFile := flag.String("config", "config.json", "Config file")
	qe.config = config.Load(*configFile)

}

func (qe *QxEye) Start() {
	util.Log("QxEye: opening cameras")
	for _, c := range qe.config.Cameras {
		if !c.Enabled {
			continue // Ignore
		}
		cam := camera.NewCamera(c.Name, c.Endpoint)
		cam.Open()
		if c.Matcher != "" {
			found := false
			for _, m := range qe.config.Matchers {
				if m.Name == c.Matcher {
					cam.SetMatcher(m.Name, m.Type, m.Params)
					found = true
					break
				}
			}
			if !found {
				util.Warn("QxEye: matcher '", c.Matcher, "' not found for camera '", c.Name, "'")
			}
		}
		if c.Persistance > 0 {
			cam.SetPersistence(c.Persistance)
		}
		for _, z := range c.Zones {
			zone := cam.SetZone(z.Name, z.Mask)
			zone.SetColor(z.Color)
			zone.SetIgnore(z.Ignore)
		}

		qe.cams = append(qe.cams, cam)
	}

	util.Log("QxEye: start watching")
	ch := make(chan camera.CameraEvent)
	for _, c := range qe.cams {
		go c.Watch(ch)
	}

	for {
		select {
		case e := <-ch:
			go onEvent(qe, e)
		}
	}
}

func main() {
	qe := New()
	qe.Start()
}
