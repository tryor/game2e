package game2e

import (
	"image/color"
	"sync/atomic"
	"time"

	"github.com/tryor/game2e/log"
	. "github.com/tryor/game2e/util"

	//	"github.com/google/gxui"
	. "github.com/tryor/eui"
)

var AllWidgetCount int32

type IWidget interface {
	IElement
	IsBorder() bool
	SetBorder(b bool)
	GetBorderWidth() float32
	SetBorderWidth(w float32)
	GetBorderColor() color.Color
	SetBorderColor(c color.Color)
	IsFill() bool
	SetFill(b bool)
	GetFillColor() color.Color
	SetFillColor(c color.Color)
	GetUserFillBrush() IBrush
	SetUserFillBrush(brush IBrush)

	//从当前参照点(refPointX, refPointY)，移动到points指定的目标点序列
	//在每次执行移动之前会调用before函数， 返回false,结束移动
	//到达目的地后调用finish函数, 返回false，结束移动
	//step步长，像素点为单位， speed速度，值越小速度越慢，单位:时钟
	//adjustVisibleRegion 可通过此值标示是否调整视口区域，可用于地图卷动效果，一般仅当主角色移动时设置此值
	//renderingSync 当场景正渲染时是否同步。渲染同步可以防止精灵和地图抖动
	MoveTos(ts *TaskSchedule, points []Point, step REAL, speed int,
		adjustVisibleRegion bool, renderingSync bool,
		before func(el IElement, targetX, targetY int) bool,
		after func(el IElement, finish bool) bool) (angle float32)

	SetFeature(feature *ResFeatureInfo)
}

type Widget struct {
	*Element
	PathSupport
	//	UnionSubWidgetsBoundRect bool
	Border          bool
	BorderWidth     float32
	BorderColor     color.Color
	BorderPath      IPath
	borderPen       IPen
	lastBorderColor color.Color
	lastBorderWidth float32

	Fill          bool
	FillColor     color.Color
	fillBrush     IBrush //基于FillColor的Brush
	UserFillBrush IBrush //指定填充的Brush
	lastFillColor color.Color

	timer *Timer

	TestFlag string

	LastError error

	onPaintEventSub EventSubscription
}

var defaultBorderColor color.Color = color.NRGBA{0, 0, 0, 255}
var defaultFillColor color.Color = color.NRGBA{255, 255, 255, 255}

func NewWidget() *Widget {
	w := &Widget{Element: NewElement()}
	w.Self = w
	w.BorderWidth = 1.0
	w.BorderColor = defaultBorderColor
	w.FillColor = defaultFillColor

	atomic.AddInt32(&AllWidgetCount, 1)
	return w
}

func (this *Widget) Destroy() {
	if this.IsDestroyed() {
		log.Warn("object is already destroyed")
		return
	}
	atomic.AddInt32(&AllWidgetCount, -1)

	this.Self.(IElement).RedrawIntersection()
	this.ReleasePath()
	if this.BorderPath != nil {
		this.BorderPath.Release()
		this.BorderPath = nil
	}
	if this.borderPen != nil {
		this.borderPen.Release()
		this.borderPen = nil
	}
	if this.fillBrush != nil {
		this.fillBrush.Release()
		this.fillBrush = nil
	}

	if this.timer != nil {
		this.timer.Stop()
	}

	this.clearonPaintEventSub()

	this.Element.Destroy()

}

func (this *Widget) SetFeature(feature *ResFeatureInfo) {
	this.SetCoordinate(feature.Bounds.X, feature.Bounds.Y)
	if feature.Bounds.W > 0 {
		this.SetWidth(feature.Bounds.W)
	}
	if feature.Bounds.H > 0 {
		this.SetHeight(feature.Bounds.H)
	}
	this.Border = feature.Border
	if feature.BorderColor.A > 0 {
		this.BorderColor = feature.BorderColor
	} else {
		this.BorderColor = defaultBorderColor
	}
	if feature.BorderWidth > 0 {
		this.BorderWidth = feature.BorderWidth
	}
	this.Fill = feature.Fill
	if feature.FillColor.A > 0 {
		this.FillColor = feature.FillColor
	} else {
		this.FillColor = defaultFillColor
	}
	this.GetSelf().SetAlignment(feature.Alignment)
	this.GetSelf().SetAnchorPoint(feature.Anchorpoint.X, feature.Anchorpoint.Y)
	this.GetSelf().SetOrderZ(feature.OrderZ)

}

func (this *Widget) IsBorder() bool {
	return this.Border
}
func (this *Widget) SetBorder(b bool) {
	this.Border = b
}
func (this *Widget) GetBorderWidth() float32 {
	return this.BorderWidth
}
func (this *Widget) SetBorderWidth(w float32) {
	this.BorderWidth = w
}
func (this *Widget) GetBorderColor() color.Color {
	return this.BorderColor
}
func (this *Widget) SetBorderColor(c color.Color) {
	this.BorderColor = c
}

func (this *Widget) IsFill() bool {
	return this.Fill
}
func (this *Widget) SetFill(b bool) {
	this.Fill = b
}
func (this *Widget) GetFillColor() color.Color {
	return this.FillColor
}
func (this *Widget) SetFillColor(c color.Color) {
	this.FillColor = c
}
func (this *Widget) GetUserFillBrush() IBrush {
	return this.UserFillBrush
}
func (this *Widget) SetUserFillBrush(brush IBrush) {
	this.UserFillBrush = brush
}

func (this *Widget) CreatePath() {
	//log.Infof("Id:%v", this.GetId())
	//this.Element.CreatePath()

	if this.Border || this.Fill {
		this.BorderPath = InitPath(this.BorderPath, this.Layer.GetGraphicsEngine())
		brect := this.Self.CreateBoundRect() //GetBoundRect()

		switch this.GetType() {
		case ElementTypeEllipse:
			this.BorderPath.AddEllipse(brect)
			err := this.BorderPath.LastError()
			if err != nil {
				log.Error("AddEllipse error, ", err)
			}
		case ElementTypePolygon:
		default:
			this.BorderPath.AddRectangle(brect)
			err := this.BorderPath.LastError()
			if err != nil {
				log.Error("AddRectangle error, ", err)
			}
		}

	}

	if this.Border {
		this.ClipRegionAdjustValue = int(this.BorderWidth/2) + 2
		ge := this.GetGraphicsEngine()
		if ge != nil {
			if this.borderPen == nil ||
				this.BorderColor != this.lastBorderColor ||
				this.BorderWidth != this.lastBorderWidth {
				if this.borderPen != nil {
					this.borderPen.Release()
				}
				this.borderPen, _ = ge.NewPen(this.BorderColor, this.BorderWidth)
				this.lastBorderColor = this.BorderColor
				this.lastBorderWidth = this.BorderWidth
			}

		}
	}

	if this.Fill && this.UserFillBrush == nil {
		ge := this.GetGraphicsEngine()
		if ge != nil {
			if this.fillBrush == nil {
				this.fillBrush, _ = ge.NewBrush(this.FillColor)
				this.lastFillColor = this.FillColor
			} else if this.FillColor != this.lastFillColor {
				this.fillBrush.Release()
				this.fillBrush, _ = ge.NewBrush(this.FillColor)
				this.lastFillColor = this.FillColor
			}
		}
	}

	this.Element.CreatePath()
}

func (this *Widget) Draw(ge IGraphicsEngine) {

	var rm RenderMode
	if this.Border {
		if this.borderPen == nil {
			ge.SetStrokeColor(this.BorderColor)
			ge.SetLineWidth(this.BorderWidth)
		} else {
			ge.SetRenderPen(this.borderPen)
		}
		rm |= OnlyStroke
	}

	if this.Fill {
		if this.UserFillBrush != nil {
			ge.SetRenderBrush(this.UserFillBrush)
		} else if this.fillBrush != nil {
			ge.SetRenderBrush(this.fillBrush)
		} else {
			ge.SetFillColor(this.FillColor)
		}
		rm |= OnlyFill
	}

	//	if this.GetTag() == "ComboBox" {
	//		log.Infof("2 this.Border:%v, rm:%v, this.BorderPath:%v", this.Border, rm, this.BorderPath)
	//	}

	if rm > 0 {
		ge.AddPaths(this.BorderPath)
		ge.Render(rm)
	}

	this.Element.Draw(ge)

}

func (this *Widget) clearonPaintEventSub() {
	if this.onPaintEventSub != nil {
		this.onPaintEventSub.Unlisten()
		this.onPaintEventSub = nil
	}
}

//根据步长和速率移动到目标位置
//在移动过程中执行 before 和 after
func (this *Widget) MoveTos(ts *TaskSchedule, points []Point, step REAL, speed int, adjustVisibleRegion, renderingSync bool,
	before func(el IElement, x, y int) bool, after func(el IElement, finish bool) bool) (angle float32) {
	if len(points) == 0 {
		return
	}
	if this.IsDestroyed() {
		return
	}
	self := this.GetSelf()
	if self.IsMoving() && this.timer != nil {
		this.timer.Stop()
	}
	self.SetMoving(true)

	pointIdx := 0
	pointCount := len(points)
	firstPoint := points[pointIdx]
	ex, ey := firstPoint.X, firstPoint.Y
	x, y := float64(self.ReferencePointX()), float64(self.ReferencePointY())
	dstc := REAL(Distance(x, y, float64(ex), float64(ey)))
	angle = float32(Angle(x, y, float64(ex), float64(ey)))

	moveddstc := REAL(0)
	stepf := REAL(step)

	layer := self.GetLayer()
	page := layer.GetDrawPage()

	this.timer = NewTimer(ts, func(delayIndex int) {
		if this.IsDestroyed() {
			self.SetMoving(false)
			if this.timer != nil {
				this.timer.Stop()
			}
			return
		}

		//		if renderingSync && page.IsRendering() {
		//			if adjustVisibleRegion {
		//				for page.IsRendering() {
		//					time.Sleep(time.Microsecond * 500)
		//					//this.timer.Execute()
		//					//return
		//				}
		//			} else {
		//				time.Sleep(time.Millisecond * 2)
		//				this.timer.Execute()
		//				return
		//			}
		//		}

		if renderingSync && page.IsDrawing() {
			if adjustVisibleRegion {
				//				this.timer.Stop()
				//				//var onPaintEventSub EventSubscription
				//				this.clearonPaintEventSub()
				//				this.onPaintEventSub = page.OnPaint(func(e *PaintEvent) {
				//					go func() {
				//						//onPaintEventSub.Unlisten()
				//						this.clearonPaintEventSub()
				//						if !this.Self.IsDestroyed() {
				//							this.timer.Start()
				//							this.moveTo(&moveddstc, stepf, step, dstc, &x, &y, angle,
				//								adjustVisibleRegion, &pointIdx, pointCount, points,
				//								before, after)
				//						} else {
				//							log.Info("object is destroyed!")
				//						}
				//					}()
				//				})
				//				return

				page.SetRenderEnable(false)
				defer page.SetRenderEnable(true)
				for page.IsDrawing() {
					time.Sleep(time.Millisecond)
				}

			} else {
				this.timer.Stop()
				go func() {
					for page.IsDrawing() {
						time.Sleep(time.Millisecond * 2)
					}
					this.timer.Start()
					this.moveTo(&moveddstc, stepf, step, &dstc, &x, &y, &angle,
						adjustVisibleRegion, &pointIdx, pointCount, points,
						before, after)

				}()
				return
			}
		}

		this.moveTo(&moveddstc, stepf, step, &dstc, &x, &y, &angle,
			adjustVisibleRegion, &pointIdx, pointCount, points,
			before, after)

		//		moveddstc += stepf
		//		moveend := false
		//		if moveddstc >= dstc {
		//			stepf = stepf - (moveddstc - dstc)
		//			moveend = true
		//		}
		//		if stepf > 0 {
		//			x1f64, y1f64 := CalculatePoint((x), (y), float64(angle), float64(stepf))
		//			x1, y1 := REAL(x1f64), REAL(y1f64)
		//			if before != nil {
		//				if !before(self, int(x1), int(y1)) {
		//					//如果回调函数返回false, 就停止移动
		//					self.SetMoving(false)
		//					if this.timer != nil {
		//						this.timer.Stop()
		//					}
		//					return
		//				}
		//			}

		//			if adjustVisibleRegion {
		//				vr := layer.GetVisibleRegionF()
		//				lfw, lfh := vr.W/2, vr.H/2 //vr.Dx()/2, vr.Dy()/2
		//				verifyx := x1 - vr.X - lfw //int(x1) - float64(vr.Min.X) - lfw
		//				verifyy := y1 - vr.Y - lfh //int(y1) - vr.Min.Y - lfh
		//				avrx := REAL(0)
		//				//			if mx-this.X() > lfw || vr.Min.X > 0 {
		//				if x1 > lfw || vr.X > 0 {
		//					avrx = verifyx
		//				}
		//				avry := REAL(0)
		//				//			if my-this.Y() > lfh || vr.Min.Y > 0 {
		//				if y1 > lfh || vr.Y > 0 {
		//					avry = verifyy
		//				}
		//				//log.Infof("%f, %f", avrx, avry)
		//				page.AdjustVisibleRegion(avrx, avry)

		//			}

		//			x, y = float64(x1), float64(y1)

		//			self.MoveTo(int(x1), int(y1), angle)
		//			self.CreateBoundRect()

		//		}
		//		if moveend {
		//			pointIdx++
		//			if pointIdx < pointCount {
		//				point := points[pointIdx]
		//				ex, ey := point.X, point.Y
		//				x, y = float64(self.ReferencePointX()), float64(self.ReferencePointY())
		//				dstc = REAL(Distance(x, y, float64(ex), float64(ey)))
		//				angle = float32(Angle(x, y, float64(ex), float64(ey)))
		//				moveddstc = REAL(0)
		//				stepf = REAL(step)
		//				moveend = false
		//			}
		//		}

		//		if after != nil {
		//			if !after(self, moveend) {
		//				moveend = true
		//			}
		//		}

		//		if moveend {
		//			self.SetMoving(false)
		//			if this.timer != nil {
		//				this.timer.Stop()
		//			}
		//		}

	}, speed)
	this.timer.Start()
	return
}

func (this *Widget) moveTo(moveddstc *REAL, stepf, step REAL, dstc *REAL, x, y *float64, angle *float32,
	adjustVisibleRegion bool, pointIdx *int, pointCount int, points []Point,
	before func(el IElement, x, y int) bool, after func(el IElement, finish bool) bool) {

	self := this.Self
	if this.IsDestroyed() {
		self.SetMoving(false)
		if this.timer != nil {
			this.timer.Stop()
		}
		return
	}

	layer := self.GetLayer()
	page := layer.GetDrawPage()

	(*moveddstc) += stepf
	moveend := false
	if *moveddstc >= *dstc {
		stepf = stepf - (*moveddstc - *dstc)
		moveend = true
	}
	if stepf > 0 {
		x1f64, y1f64 := CalculatePoint((*x), (*y), float64(*angle), float64(stepf))
		x1, y1 := REAL(x1f64), REAL(y1f64)
		if before != nil {
			if !before(self, int(x1), int(y1)) {
				//如果回调函数返回false, 就停止移动
				self.SetMoving(false)
				if this.timer != nil {
					this.timer.Stop()
				}
				return
			}
		}

		self.MoveTo(int(x1), int(y1), *angle)
		self.CreateBoundRect()

		if adjustVisibleRegion {
			vr := layer.GetVisibleRegionF()
			lfw, lfh := vr.W/2, vr.H/2 //vr.Dx()/2, vr.Dy()/2
			verifyx := x1 - vr.X - lfw //int(x1) - float64(vr.Min.X) - lfw
			verifyy := y1 - vr.Y - lfh //int(y1) - vr.Min.Y - lfh
			avrx := REAL(0)
			//			if mx-this.X() > lfw || vr.Min.X > 0 {
			if x1 > lfw || vr.X > 0 {
				avrx = verifyx
			}
			avry := REAL(0)
			//			if my-this.Y() > lfh || vr.Min.Y > 0 {
			if y1 > lfh || vr.Y > 0 {
				avry = verifyy
			}
			//log.Infof("%f, %f", avrx, avry)
			page.AdjustVisibleRegion(avrx, avry)

		}

		*x, *y = float64(x1), float64(y1)

		//		self.MoveTo(int(x1), int(y1), angle)
		//		self.CreateBoundRect()

	}
	if moveend {
		(*pointIdx)++
		if (*pointIdx) < pointCount {
			point := points[*pointIdx]
			ex, ey := point.X, point.Y
			*x, *y = float64(self.ReferencePointX()), float64(self.ReferencePointY())
			*dstc = REAL(Distance(*x, *y, float64(ex), float64(ey)))
			*angle = float32(Angle(*x, *y, float64(ex), float64(ey)))
			*moveddstc = REAL(0)
			stepf = REAL(step)
			moveend = false
		}
	}

	if after != nil {
		if !after(self, moveend) {
			moveend = true
		}
	}

	if moveend {
		self.SetMoving(false)
		if this.timer != nil {
			this.timer.Stop()
		}
	}

}
