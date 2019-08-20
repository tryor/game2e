package util

import (
	"math"

	. "github.com/tryor/eui"
)

//两点之间的距离
func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func DistanceI(x1, y1, x2, y2 int) int {
	return int(math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))))
}

//返回值为 0 - 359.999999
func Angle(x1, y1, x2, y2 float64) float64 {
	x := x2 - x1
	y := y2 - y1
	if x == 0 && y == 0 {
		return 0
	}
	d := math.Sqrt(x*x + y*y)
	angle := 180 / (math.Pi / math.Acos(x/d))
	if y > 0 {
		angle = 360 - angle
	}
	if math.IsNaN(angle) {
		//		println("Angle:", x, y)
		angle = 0
	}
	return angle
}

func AngleI(x1, y1, x2, y2 int) float64 {
	return Angle(float64(x1), float64(y1), float64(x2), float64(y2))
}

//参数（起点坐标，角度，斜边长（距离））
//逆时针角度为正, 顺时针角度为负，
func CalculatePoint(x1, y1 float64, angle float64, distance float64) (x2, y2 float64) {
	radian := -angle * math.Pi / 180
	xMargin := math.Cos(radian) * distance
	yMargin := math.Sin(radian) * distance
	return x1 + xMargin, y1 + yMargin
}

func CalculatePointI(x1, y1 int, angle float64, distance int) (x2, y2 int) {
	fx2, fy2 := CalculatePoint(float64(x1), float64(y1), angle, float64(distance))
	x2, y2 = int(fx2), int(fy2)
	return
}

//判断点(ri_x, ri_y)是否在椭圆内
//(ri_centreX, ri_centreY)椭圆心坐标
//ri_hradius, ri_vradius 横向半径，纵向半径
func PointInEllipse(ri_x, ri_y, ri_centreX, ri_centreY, ri_hradius, ri_vradius int) bool {

	var i_X, i_Y int
	var i_SX, i_SY int

	i_X = ri_x - ri_centreX
	i_Y = ri_y - ri_centreY

	i_SX = ((i_X << 8) / ri_hradius)
	i_SY = ((i_Y << 8) / ri_vradius)
	if i_SX > (1<<8) || i_SY > (1<<8) {
		return false
	} else {
		return (i_SX*i_SX)+(i_SY*i_SY) <= (1 << 16)
	}
}

func PointInPolygon(x, y int, points []Point) bool {
	nCross := 0
	//for _, p := range points {
	xf := float32(x)
	//	yf := float32(y)
	nCount := len(points)
	for i := 0; i < nCount; i++ {
		p1 := points[i]
		p2 := points[(i+1)%nCount]
		// 求解 y=p.y 与 p1p2 的交点
		if p1.Y == p2.Y { // p1p2 与 y=p0.y平行
			continue
		}
		if y < Min(p1.Y, p2.Y) { // 交点在p1p2延长线上
			continue
		}
		if y >= Max(p1.Y, p2.Y) { // 交点在p1p2延长线上
			continue
		}
		// 求交点的 X 坐标 --------------------------------------------------------------
		xc := float32(y-p1.Y)*float32(p2.X-p1.X)/float32(p2.Y-p1.Y) + float32(p1.X)

		if xc > xf {
			nCross++ // 只统计单边交点
		}
	}
	return nCross%2 == 1
}
