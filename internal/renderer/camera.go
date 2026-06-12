package renderer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	dragSensitivity = 0.005
	zoomSensitivity = 0.5
	minRadius       = 2.5
	maxRadius       = 30.0
)

func (r *Renderer) handleCameraInput() {
	if wheel := rl.GetMouseWheelMove(); wheel != 0 {
		r.radius = rl.Clamp(r.radius-wheel*zoomSensitivity, minRadius, maxRadius)
		r.syncCameraPosition()
	}
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		delta := rl.GetMouseDelta()
		r.azimuth += delta.X * dragSensitivity
		r.altitude += delta.Y * dragSensitivity
		r.altitude = rl.Clamp(r.altitude, float32(-math.Pi/2)+0.01, float32(math.Pi/2)-0.01)
		r.syncCameraPosition()
	}
}

func (r *Renderer) syncCameraPosition() {
	cosAlt := float32(math.Cos(float64(r.altitude)))
	sinAlt := float32(math.Sin(float64(r.altitude)))
	cosAz := float32(math.Cos(float64(r.azimuth)))
	sinAz := float32(math.Sin(float64(r.azimuth)))
	r.camera.Position = rl.NewVector3(
		r.radius*cosAlt*cosAz,
		r.radius*sinAlt,
		r.radius*cosAlt*sinAz,
	)
}
