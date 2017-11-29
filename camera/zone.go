package camera

import (
	"math"
	"strings"
)

type cameraZoneRegion struct {
	top    float64
	left   float64
	bottom float64
	right  float64
}

type cameraZone struct {
	name    string
	ignore  bool
	color   string
	regions []cameraZoneRegion
}

func newCameraZone(name string, mask []string) cameraZone {
	var z cameraZone
	z.name = name
	z.setMask(mask)
	return z
}

func (cz *cameraZone) SetColor(color string) {
	cz.color = color
}
func (cz *cameraZone) SetIgnore(ignore bool) {
	cz.ignore = ignore
}

func (cz *cameraZone) setMask(mask []string) {
	numRows := float64(len(mask))
	for i, row := range mask {
		numCols := float64(len(row))
		var z *cameraZoneRegion
		for j, col := range strings.Split(row, "") {
			if col != " " {
				if z != nil {
					continue
				}
				z = new(cameraZoneRegion)
				z.top = float64(i) / numRows
				z.left = float64(j) / numCols
				z.bottom = z.top + (1.0 / numRows)
			} else if z != nil {
				z.right = float64(j) / numCols
				cz.regions = append(cz.regions, *z)
				z = nil
			}
		}
		if z != nil {
			z.right = 1.0
			cz.regions = append(cz.regions, *z)
		}
	}
}

func (cz *cameraZone) intersect(imgWidth, imgHeight, left, top, width, height int) float64 {
	coverage := 0.0
	for _, z := range cz.regions {
		x1 := math.Max(float64(left), z.left*float64(imgWidth))
		x2 := math.Min(float64(left+width), z.right*float64(imgWidth))
		y1 := math.Max(float64(top), z.top*float64(imgHeight))
		y2 := math.Min(float64(top+height), z.bottom*float64(imgHeight))
		if x1 < x2 && y1 < y2 {
			coverage += (x2 - x1) * (y2 - y1) / float64(width*height)
		}
	}
	return coverage
}
