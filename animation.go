package game2e

import (
	"time"

	"tryor/game2e/log"
	. "tryor/game2e/util"

	//"github.com/google/gxui"
	. "github.com/tryor/eui"
)

type Animation struct {
	*Widget
	Cached        bool
	Loop          bool
	AnimationName string
	AnimationInfo *AnimationInfo
	cachedImages  []*ImageByResource
	current       *ImageByResource
	frameIndex    int
	playEnd       bool //仅当Loop==false时，此值才有可能被设置
	timer         *Timer

	timer2 *time.Timer

	onPlayStopEvent Event
}

func NewAnimation(ts *TaskSchedule, x, y int, animationName string, frameindexs []int, cached bool) *Animation {
	ainfo := AnimationInfos[animationName]
	if ainfo == nil {
		log.Error("animation " + animationName + " not find!")
		return nil
	}

	animation := &Animation{Widget: NewWidget(), Cached: cached, Loop: true}
	animation.Self = animation
	animation.SetCoordinate(x, y)
	animation.AnimationName = animationName
	animation.AnimationInfo = ainfo
	animation.cachedImages = make([]*ImageByResource, 0) //make([]*CachedImage, len(ainfo.Frames))
	animation.onPlayStopEvent = CreateEvent(func(a *Animation) {})

	w, h := ainfo.Width, ainfo.Height
	animation.SetWidth(w)
	animation.SetHeight(h)

	delays := make([]int, 0)

	//for i, aframe := range ainfo.Frames {
	if frameindexs == nil || len(frameindexs) == 0 {
		size := len(ainfo.Frames)
		frameindexs = make([]int, size)
		for i := 0; i < size; i++ {
			frameindexs[i] = i
		}
	}
	//for i := 0; i < len(ainfo.Frames); i++ {
	for _, idx := range frameindexs {
		aframe := ainfo.Frames[idx]
		cimg := NewImageByResource(0, 0, aframe.Name, aframe.Id)
		if cimg != nil {
			//println("aframe.Feature.OrderZ:", aframe.Feature.OrderZ)
			cimg.SetFeature(&aframe.Feature)
			//			animation.cachedImages[i] = cimg
			cimg.Cached = cached
			if !cached {
				cimg.SetWidth(w)
				cimg.SetHeight(h)
			}
			animation.cachedImages = append(animation.cachedImages, cimg)
			delays = append(delays, aframe.Delay)
			//			cimg.SetCoordinate(cimg.GetImageWidth()/2, cimg.GetImageHeight()/2)
		}
	}

	if len(animation.cachedImages) > 0 {
		cimg := animation.cachedImages[0]
		if cimg != nil {

			if cached {
				if w <= 0 {
					animation.SetWidth(cimg.GetImageWidth())
				}
				if h <= 0 {
					animation.SetHeight(cimg.GetImageHeight())
				}
			}
			animation.setFrameImage(cimg)

		}
	}

	if len(animation.cachedImages) > 1 {

		//		delays := make([]int, len(ainfo.Frames))
		//		for i, frame := range ainfo.Frames {
		//			delays[i] = frame.Delay
		//		}

		animation.timer = NewTimer(ts, func(delayIndex int) {
			animation.selectNextFrame(delayIndex)
		}, delays...)

		//		animation.timer2 = time.AfterFunc(time.Duration(delays[0])*time.Millisecond*10, func() {
		//			animation.selectNextFrame2()
		//		})

	}

	return animation
}

func (this *Animation) Destroy() {
	//println("Animation.Destroy:", len(this.cachedImages))
	//	if this.timer != nil {
	//		this.timer.Stop()
	//		this.timer = nil
	//	}
	this.Stop()
	for _, cachedImage := range this.cachedImages {
		if cachedImage != nil {
			this.RemoveChild(cachedImage)
			cachedImage.Destroy()
		}
	}
	this.cachedImages = this.cachedImages[0:0]

	//	this.ClearChildren()
	this.Widget.Destroy()
}

func (this *Animation) IsPlayEnd() bool {
	return this.playEnd
}

func (this *Animation) OnPlayStop(f func(a *Animation)) EventSubscription {
	return this.onPlayStopEvent.Listen(f)
}

func (this *Animation) Load() error {
	return nil
}

func (this *Animation) Start(delayIndex ...int) {
	if this.timer != nil {
		this.playEnd = false
		this.timer.Start(delayIndex...)
	}
}

func (this *Animation) Stop() {
	if this.timer != nil && this.timer.IsActive() {
		this.timer.Stop()
		//log.Info("playEnd:", this.playEnd, len(this.onPlayStopEvent.ParameterTypes()), this.Self.GetId())
		this.onPlayStopEvent.Fire(this)
	}
}

func (this *Animation) SetActive(b bool) {
	this.Widget.SetActive(b)
	if b {
		if this.Loop {
			this.Start()
		}
	} else {
		this.Stop()
	}
}

func (this *Animation) setFrameImage(fimg *ImageByResource) {
	current := this.current
	if current != nil {
		//this.Self.(IElement).RemoveChild(current)
		current.SetVisible(false)
	}
	this.current = fimg

	if !this.ExistChild(fimg.GetId()) {
		this.Self.(IElement).AddChild(fimg)
	}
	fimg.SetVisible(true)

	if this.GetParent() != nil {
		this.Self.(IElement).SetOrderZ(fimg.GetOrderZ())
	}
	this.Self.(IElement).SetModified(true)
}

func (this *Animation) selectNextFrame(frameIndex int) {
	this.frameIndex = frameIndex
	if !this.Loop && frameIndex == 0 {
		//		this.timer.Stop()
		this.playEnd = true
		this.Stop()
		return
	}

	if this.IsDestroyed() {
		log.Warn("IsDestroyed")
		this.Stop()
		return
	}

	cimg := this.cachedImages[frameIndex]
	if cimg != nil {
		this.setFrameImage(cimg)
	}
}

//func (this *Animation) selectNextFrame2() {

//	//	timespender := NewTimespender("Animation.selectNextFrame 1")
//	//	defer timespender.Print(0)
//	this.FrameIndex++
//	if !this.Loop && this.FrameIndex >= len(this.cachedImages) {
//		this.timer.Stop()
//		return
//	}

//	if this.FrameIndex >= len(this.cachedImages) {
//		this.FrameIndex = 0
//	}

//	cimg := this.cachedImages[this.FrameIndex]
//	if cimg != nil {
//		this.ClearChildren()
//		this.AddChild(cimg)
//		this.Self.(IElement).SetModified(true)
//	}

//	//	println("Animation.selectNextFrame")
//	this.timer2.Reset(time.Duration(this.AnimationInfo.Frames[this.FrameIndex].Delay) * time.Millisecond * 10)
//	//	this.timer2.Reset(this.AnimationInfo.Frames[this.FrameIndex].Delay)
//}
