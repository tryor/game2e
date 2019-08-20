package game2e

import (
	. "tryor/game2e/util"

	. "github.com/tryor/eui"
)

type GifImage struct {
	*Image
	Loop  bool
	ts    *TaskSchedule
	timer *Timer
}

func NewGifImage(ts *TaskSchedule, x, y, w, h int, filename string) *GifImage {
	img := &GifImage{ts: ts, Image: NewImage(x, y, w, h, filename), Loop: true}
	img.Self = img
	return img
}

func (this *GifImage) Destroy() {
	this.Image.Destroy()
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}
}

func (this *GifImage) SetActive(b bool) {
	this.Image.SetActive(b)
	if b {
		if this.timer != nil {
			this.timer.Start()
		}
	} else {
		if this.timer != nil {
			this.timer.Stop()
		}
	}
}

func (this *GifImage) Load() error {
	err := this.Image.Load()
	if err != nil {
		return err
	}

	if this.timer == nil {
		//		this.timer = time.AfterFunc(time.Duration(this.NativeImage.GetFrameDelay())*time.Millisecond, func() {
		//			this.selectNextFrame()
		//		})
		//		this.timer = AfterFunc(this.ts, this.NativeImage.GetFrameDelay(), func() {
		//			this.selectNextFrame()
		//		})

		this.timer = NewTimer(this.ts, func(delayIndex int) {
			this.selectNextFrame(delayIndex)
		}, this.NativeImage.GetFrameDelays()...)
		this.timer.Start()
	} else {
		//this.timer.Reset(time.Duration(this.NativeImage.GetFrameDelay()) * time.Millisecond)
		//		this.timer.Reset(this.NativeImage.GetFrameDelay())
	}

	return nil
}

func (this *GifImage) selectNextFrame(frameIndex int) {
	if !this.Loop && frameIndex == 0 {
		this.timer.Stop()
		return
	}

	this.NativeImage.SelectActiveFrame(uint(frameIndex))
	//	this.NativeImage.SelectNextFrame()
	//	if !this.IsModified() {
	this.Self.(IElement).SetModified(true)
	//	}
	//	this.timer.Reset(time.Duration(this.NativeImage.GetFrameDelay()) * time.Millisecond)
	//	this.timer.Reset(this.NativeImage.GetFrameDelay())
}
