package game2e

import (
	"image"
	"time"

	"tryor/game2e/log"
	. "tryor/game2e/util"

	//	"github.com/google/gxui"
	. "github.com/tryor/eui"
)

type IScene interface {
	IDrawPage
	Start()
	Stop()
	LoadMap(id string, vx, vy int) error
	SetRenderClockGeneratorPeriod(period time.Duration)
	SetClockGeneratorPeriod(period time.Duration)
	GetClockGenerator() *ClockGenerator
	GetTaskSchedule() *TaskSchedule
	GetLayersByType(typ LayerType) []ILayer
	GetLayerByType(typ LayerType) ILayer
	GetLayerById(id string) ILayer

	OnEnter(f func()) EventSubscription
	OnLeave(f func()) EventSubscription
}

type Scene struct {
	*DrawPage

	//clockGenerator       *ClockGenerator //默认时钟发生器
	renderClockGenerator *ClockGenerator //用于场景渲染的时钟发生器

	taskSchedule *TaskSchedule //默认TaskSchedule

	onEnterEvent Event //进入场景时事件
	onLeaveEvent Event //离开场景时事件
}

//vrwidth 可视区域宽度
//vrheight 可视区域高度
func NewScene(vrwidth, vrheight int) *Scene {
	scene := &Scene{DrawPage: NewDrawPage(image.Rect(0, 0, vrwidth, vrheight))}
	scene.Self = scene

	scene.renderClockGenerator = NewClockGenerator(time.Millisecond * DefaultRenderClockGeneratorPeriod)
	//scene.clockGenerator = NewClockGenerator(time.Millisecond * DefaultClockGeneratorPeriod)

	//scene.taskSchedule = NewTaskSchedule(scene.clockGenerator)
	scene.taskSchedule = NewTaskSchedule(scene.renderClockGenerator)

	rendercb := func(clock uint64) {
		//go scene.Render()
		scene.Render()
	}
	scene.renderClockGenerator.AddClockCallback(&rendercb)

	scene.onEnterEvent = CreateEvent(func() {})
	scene.onLeaveEvent = CreateEvent(func() {})

	//scene.SetShowStatInfoEnable(true)

	return scene
}

func (this *Scene) Destroy() {
	if this.IsDestroyed() {
		log.Info("Scene is already destroyed")
		return
	}

	this.taskSchedule.Destroy()

	//	this.clockGenerator.Stop()
	//	for this.clockGenerator.IsRunning() {
	//		log.Info("clockGenerator.IsRunning()")
	//		time.Sleep(time.Millisecond * 10)
	//	}

	this.renderClockGenerator.Stop()
	for this.renderClockGenerator.IsRunning() {
		log.Info("renderClockGenerator.IsRunning()")
		time.Sleep(time.Millisecond * 10)
	}

	for this.IsRendering() {
		log.Info("IsRendering()")
		time.Sleep(time.Millisecond * 10)
	}

	this.DrawPage.Destroy()
}

func (this *Scene) SetRenderClockGeneratorPeriod(period time.Duration) {
	this.renderClockGenerator.SetPeriod(period)
}

func (this *Scene) SetClockGeneratorPeriod(period time.Duration) {
	//this.clockGenerator.SetPeriod(period)
	this.renderClockGenerator.SetPeriod(period)
}

func (this *Scene) GetClockGenerator() *ClockGenerator {
	//return this.clockGenerator
	return this.renderClockGenerator
}

func (this *Scene) GetTaskSchedule() *TaskSchedule {
	return this.taskSchedule
}

func (this *Scene) OnEnter(f func()) EventSubscription {
	return this.onEnterEvent.Listen(f)
}

func (this *Scene) OnLeave(f func()) EventSubscription {
	return this.onLeaveEvent.Listen(f)
}

func (this *Scene) Start() {
	this.Self.StartStatTimer()
	this.renderClockGenerator.Start()
	//this.clockGenerator.Start()
	this.onEnterEvent.Fire()
}

func (this *Scene) Stop() {
	this.Self.StopStatTimer()
	//this.clockGenerator.Stop()
	this.renderClockGenerator.Stop()
	this.onLeaveEvent.Fire()
}

//id 地图ID
//vx,vy 可视区域坐标
func (this *Scene) LoadMap(id string, vx, vy int) error {
	m, err := NewMap(this.GetTaskSchedule(), vx, vy, id)
	if err != nil {
		return err
	}
	//x, y := 0, 0
	for _, layer := range m.Layers {
		err := this.AddLayer(layer)
		if err == nil {
			layer.SetCoordinate(layer.X(), layer.Y())
			this.GraphicsEngine.GetLayerMerger().InitLayerGraphicsEngine(layer.(ILayer))
			layer.(ILayer).Init()

			//			if layer.GetLayerType() == LayerTypeActive {
			//				this.mainLayer = layer
			//				this.SetFocusLayer(layer)
			//			}
		} else {
			log.Error(err)
		}
	}
	return nil
}

func (this *Scene) GetLayersByType(typ LayerType) []ILayer {
	layers := make([]ILayer, 0)
	for _, layer := range this.Self.GetLayers() {
		if layer.GetLayerType() == typ {
			layers = append(layers, layer)
		}
	}
	return layers
}

func (this *Scene) GetLayerByType(typ LayerType) ILayer {
	for _, layer := range this.Self.GetLayers() {
		if layer.GetLayerType() == typ {
			return layer
		}
	}
	return nil
}

func (this *Scene) GetLayerById(id string) ILayer {
	for _, layer := range this.Self.GetLayers() {
		if layer.GetId() == id {
			return layer
		}
	}
	return nil
}

//func (this *Scene) AdjustVisibleRegion(dx, dy REAL) {
//	//go func() {
//	//step := 15
//	speed := 10
//	speedf := REAL(speed)
//	dxi, dyi := dx/speedf, dy/speedf
//	//log.Info("dx, dy:", dx, dy)
//	var timer *Timer
//	timer = NewTimer(this.GetTaskSchedule(), func(delayIndex int) {
//		this.DrawPage.AdjustVisibleRegion(dxi, dyi)
//		//log.Info("dxi, dyi:", dxi, dyi)
//		speed -= 1
//		if speed <= 0 {
//			if timer != nil {
//				timer.Stop()
//				timer = nil
//			}
//		}
//	}, 1)
//	timer.Start()
//}
