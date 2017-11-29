package camera

/*
#cgo LDFLAGS: -s -I/usr/local/include/opencv -I/usr/local/include/opencv2 -L/usr/local/lib/ -lopencv_core -lopencv_video -lopencv_videoio -lopencv_highgui -lopencv_imgproc -lopencv_objdetect  -lopencv_imgcodecs -lopencv_dnn
#include <stdlib.h>
#include <stdbool.h>
#include "camera-binding.hpp"
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Camera struct {
	Name     string
	endpoint string
	cam      *C.struct_CameraHandler
	zones    []*cameraZone
}

type CameraMatch struct {
	Label      string
	Confidence float32
	Top        int
	Left       int
	Width      int
	Height     int
	Zones      []ZoneMatch
}

type CameraEvent struct {
	CameraName string
	Date       time.Time
	Matches    []CameraMatch
	cm         *C.struct_CameraMatch
	cam        *Camera
}

type CameraTraceFlag byte

type ZoneMatch struct {
	Name     string
	Coverage float64
}

const (
	CameraTraceNone             = CameraTraceFlag(0x00)
	CameraTraceMatches          = CameraTraceFlag(0x01)
	CameraTraceIgnoredZonesOnly = CameraTraceFlag(0x02)
	CameraTraceZones            = CameraTraceFlag(0x04)
)

func CameraTraceFlagFromString(flag string) CameraTraceFlag {
	var b CameraTraceFlag
	for _, f := range strings.Split(flag, ",") {
		switch f {
		case "matches":
			b = b + CameraTraceMatches
		case "ignoredZonesOnly":
			b = b + CameraTraceIgnoredZonesOnly
		case "zones":
			b = b + CameraTraceZones
		}
	}
	return b
}

func NewCamera(name, endpoint string) *Camera {
	c := new(Camera)
	c.Name = name
	c.endpoint = endpoint
	return c
}

func (c *Camera) SetPersistence(persistence float32) {
	if c.cam != nil {
		C.camera_setPersistence(c.cam, C.float(persistence))
	}
}

func (c *Camera) SetMatcher(matcher string, matcherType string, params map[string]interface{}) {
	if c.cam != nil {
		// Serialize config
		s := ""
		for k, v := range params {
			s += k + "=" + fmt.Sprint(v) + "\n"
		}
		C.matcher_init(C.CString(matcher), C.CString(matcherType), C.CString(s))
		C.camera_setMatcher(c.cam, C.CString(matcher))
	}
}

func (c *Camera) SetZone(name string, mask []string) *cameraZone {
	z := newCameraZone(name, mask)
	c.zones = append(c.zones, &z)
	return &z
}

func (c *Camera) Open() error {
	c.cam = C.camera_create(C.CString(c.endpoint))
	if c.cam == nil {
		return errors.New("Unable to open camera")
	}
	return nil
}

func (c *Camera) Close() {
	if c.cam != nil {
		C.camera_destroy(c.cam)
	}
}

func (c *Camera) initZoneMask() {
	var node *C.struct_ZoneList
	var head *C.struct_ZoneList
	for _, z := range c.zones {
		if z.ignore == false {
			continue // Ignore
		}
		for _, r := range z.regions {
			item := C.zoneItem_new()
			item.top = C.int(int(r.top * 100.0))
			item.left = C.int(int(r.left * 100.0))
			item.width = C.int(int((r.right - r.left) * 100.0))
			item.height = C.int(int((r.bottom - r.top) * 100.0))
			node = C.zoneList_append(node, item)
			if head == nil {
				head = node
			}
		}
	}
	C.camera_setMask(c.cam, head)
	C.zoneList_destroy(head, true)
}

func (c *Camera) Watch(ch chan CameraEvent) {
	c.initZoneMask()
	go C.camera_start(c.cam)
	defer func() {
		c.Close()
	}()
	for {
		cm := C.cameraMatch_init(c.cam)
		if C.cameraMatch_check(cm) {
			e := CameraEvent{}
			e.CameraName = c.Name
			e.Date = time.Now()
			e.cm = cm
			e.cam = c
			node := C.cameraMatch_getMatches(cm)
			imSize := C.cameraMatch_getImageSize(cm)
			for node != nil {
				item := node.item
				m := CameraMatch{}
				m.Label = C.GoString(item.label)
				m.Confidence = float32(item.confidence)
				m.Top = int(item.top)
				m.Left = int(item.left)
				m.Width = int(item.width)
				m.Height = int(item.height)
				for _, z := range c.zones {
					if z.ignore {
						continue
					}
					cov := z.intersect(int(imSize.width), int(imSize.height), m.Left, m.Top, m.Width, m.Height)
					if cov > 0.0 {
						m.Zones = append(m.Zones, ZoneMatch{z.name, cov})
					}
				}
				e.Matches = append(e.Matches, m)
				node = node.next
			}
			ch <- e
		}
	}
}

func (c *Camera) Record() {
	// TODO
}

func (c *CameraEvent) Capture(filename string, matches []CameraMatch, trace CameraTraceFlag) {
	var node *C.struct_ZoneList
	var head *C.struct_ZoneList
	imSize := C.cameraMatch_getImageSize(c.cm)
	imWidth := float64(imSize.width)
	imHeight := float64(imSize.height)

	if trace&CameraTraceZones != 0 || trace&CameraTraceIgnoredZonesOnly != 0 {
		// Add zones
		for _, z := range c.cam.zones {
			for _, r := range z.regions {
				if trace&CameraTraceZones == 0 && !z.ignore {
					continue
				}
				item := C.zoneItem_new()
				if z.name != "" {
					item.label = C.CString(z.name)
				} else {
					item.label = nil
				}
				item.borderSize = 1
				item.top = C.int(int(r.top * imHeight))
				item.left = C.int(int(r.left * imWidth))
				item.width = C.int(int((r.right - r.left) * imWidth))
				item.height = C.int(int((r.bottom - r.top) * imHeight))
				if z.ignore {
					item.color = C.CString("000000")
					item.fillOpacity = 1
				} else {
					if z.color == "" {
						item.color = C.CString("00FFFF")
					} else {
						item.color = C.CString(z.color)
					}
					item.fillOpacity = 0.1
				}
				node = C.zoneList_append(node, item)
				if head == nil {
					head = node
				}
			}
		}
	}

	if trace&CameraTraceMatches != 0 {
		// Add matches
		for _, m := range matches {
			item := C.zoneItem_new()
			item.label = C.CString(fmt.Sprintf("%s: %.1f%%", m.Label, m.Confidence*100))
			item.top = C.int(m.Top)
			item.left = C.int(m.Left)
			item.width = C.int(m.Width)
			item.height = C.int(m.Height)
			item.color = C.CString("FF0000")
			item.borderSize = 2
			item.fillOpacity = 0
			node = C.zoneList_append(node, item)
			if head == nil {
				head = node
			}
		}
	}
	C.cameraMatch_capture(c.cm, C.CString(filename), head)
	C.zoneList_destroy(head, true)
}

func (c *CameraEvent) Release() {
	C.cameraMatch_destroy(c.cm)
}
