package game2e

import (
	"fmt"

	. "github.com/tryor/eui"
	. "github.com/tryor/eui/graphicsengine"
)

type PathSupport struct {
	Path     IPath
	Rendered bool
	//	PathCreateCount int
}

//func (this *PathSupport) InitPath(visibleRegion IVisibleRegion) {
func (this *PathSupport) InitPath(ge IGraphicsEngine) {
	if this.Path == nil {
		this.Path = NewPath(ge)
		this.Rendered = false
	} else if this.Rendered {
		this.Path.Reset()
		this.Rendered = false
	} else {
		//路径已经被创建或未渲染输出，不能重复创建
		//panic(errors.New("cannot create repeat path!"))
		fmt.Println("cannot create repeat path! ", this.Path)
		this.Path.Reset()
	}
	//this.Path.SetVisibleRegion(visibleRegion)

}

func (this *PathSupport) RenderPath(ge IGraphicsEngine, mode ...RenderMode) {
	//	if this.Path != nil {
	ge.AddPaths(this.Path)
	ge.Render(mode...)
	//	}
	this.Rendered = true
	//	this.PathCreateCount = 0
}

func (this *PathSupport) ReleasePath() {
	if this.Path != nil {
		this.Path.Release()
		this.Path = nil
	}

}

//func InitPath(path IPath, visibleRegion IVisibleRegion) IPath {
func InitPath(path IPath, ge IGraphicsEngine) IPath {

	if path == nil {
		path = NewPath(ge)
	} else {
		path.Reset()
		//path.Release()
		//path = NewPath(ge)
	}
	//	path.SetVisibleRegion(visibleRegion)
	return path
}

func NewPath(ge IGraphicsEngine) IPath {
	switch GraphicsEngineType {
	case Default:
		return NewDefaultPath()
	//case Draw2d:
	//	return NewDraw2dPath()
	case Gdiplus:
		return NewGdiplusPath(ge)
	}
	panic("GraphicsEngineType error")
}
