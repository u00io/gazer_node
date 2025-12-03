package render

import (
	"image"
)

func Render(script string, data map[string]interface{}) (img image.Image, err error) {
	img = image.NewRGBA(image.Rectangle{Max: image.Point{X: 320, Y: 240}})
	return img, nil
}
