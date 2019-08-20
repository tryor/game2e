package game2e

import (
	//"log"
	"time"

	. "github.com/tryor/eui"
)

type IDirector interface {
	Start()
	Shutdown()
	Destroy()
	RunScene(scene IScene)
	GetViewport() IViewport
	GetGraphicsEngine() IGraphicsEngine
}

type Director struct {
	viewport       IViewport
	graphicsEngine IGraphicsEngine

	currentScene IScene
	//	currentSceneOnMouseSub gxui.EventSubscription

	renderClockGeneratorPeriod time.Duration
	clockGeneratorPeriod       time.Duration
}

func NewDirector(viewport IViewport, graphicsEngine IGraphicsEngine) *Director {
	director := &Director{viewport: viewport, graphicsEngine: graphicsEngine}
	director.clockGeneratorPeriod = time.Millisecond * DefaultClockGeneratorPeriod
	director.renderClockGeneratorPeriod = time.Millisecond * DefaultRenderClockGeneratorPeriod

	director.viewport.OnMouse(func(me IMouseEvent) {
		if director.currentScene != nil {
			director.currentScene.TrackEvent(me)
		}
	})

	director.viewport.OnKey(func(ke *KeyEvent) {
		//log.Printf("OnKey. key:%v, %c\n", ke.Key, ke.Char)
		if director.currentScene != nil {
			director.currentScene.TrackEvent(ke)
		}
	})

	director.viewport.OnKeyChar(func(ke *KeyEvent) {
		//log.Printf("OnKeyChar. key:%v, %c\n", ke.Key, ke.Char)
		if director.currentScene != nil {
			director.currentScene.TrackEvent(ke)
		}
	})

	return director
}

func (this *Director) SetRenderClockGeneratorPeriod(period time.Duration) {
	if this.renderClockGeneratorPeriod != period {
		this.renderClockGeneratorPeriod = period
		if this.currentScene != nil {
			this.currentScene.SetRenderClockGeneratorPeriod(period)
		}
	}
}

func (this *Director) SetClockGeneratorPeriod(period time.Duration) {
	if this.clockGeneratorPeriod != period {
		this.clockGeneratorPeriod = period
		if this.currentScene != nil {
			this.currentScene.SetClockGeneratorPeriod(period)
		}
	}
}

func (this *Director) Start() {

}

func (this *Director) Shutdown() {
	this.viewport.Close()
}

func (this *Director) Destroy() {

}

func (this *Director) RunScene(scene IScene) {
	if this.currentScene != nil {
		this.currentScene.Stop()
	}
	this.currentScene = scene
	scene.SetClockGeneratorPeriod(this.clockGeneratorPeriod)
	scene.SetRenderClockGeneratorPeriod(this.renderClockGeneratorPeriod)

	scene.Start()
}

func (this *Director) GetViewport() IViewport {
	return this.viewport
}
func (this *Director) GetGraphicsEngine() IGraphicsEngine {
	return this.graphicsEngine
}
