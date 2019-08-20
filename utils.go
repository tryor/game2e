package game2e

import (
	"image/color"
	"log"
	"syscall"

	. "github.com/tryor/winapi"
)

func InitCursor(cursorname *uint16) syscall.Handle {
	c, err := LoadCursor(0, cursorname)
	if err != nil {
		log.Println("LoadCursor error, ", err)
	}
	return syscall.Handle(c)
}

func GetColorRGBA(c color.Color) (r, g, b, a uint8) {
	switch rgba := c.(type) {
	case color.NRGBA:
		r, g, b, a = rgba.R, rgba.G, rgba.B, rgba.A
		return
	case color.RGBA:
		r, g, b, a = rgba.R, rgba.G, rgba.B, rgba.A
		return
	}
	r_, g_, b_, a_ := c.RGBA()
	r, g, b, a = uint8(r_), uint8(g_), uint8(b_), uint8(a_)
	return
}
