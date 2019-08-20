package util

import (
	"fmt"
	"testing"
)

func Test_Math(t *testing.T) {
	fmt.Printf("Distance:%f, %f\n", Distance(200, 200, 400, 200), Distance(200, 200, 400, 200))
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 400, 200), Angle(200, 200, 400, 200)) //0度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 400, 100), Angle(200, 200, 400, 100)) //26度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 400, 0), Angle(200, 200, 400, 0))     //45度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 200, 100), Angle(200, 200, 200, 100)) //90度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 50), Angle(200, 200, 100, 50))   //123度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 100), Angle(200, 200, 100, 100)) //135度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 150), Angle(200, 200, 100, 150)) //153度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 190), Angle(200, 200, 100, 190)) //174度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 200), Angle(200, 200, 100, 200)) //180度

	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 250), Angle(200, 200, 100, 250)) //206
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 100, 300), Angle(200, 200, 100, 300)) //225
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 200, 300), Angle(200, 200, 200, 300)) //270度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 250, 300), Angle(200, 200, 250, 300)) //296度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 270, 300), Angle(200, 200, 270, 300)) //304度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 300, 300), Angle(200, 200, 300, 300)) //315度

	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 300, 270), Angle(200, 200, 300, 270)) //325度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 300, 250), Angle(200, 200, 300, 250)) //333度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 300, 210), Angle(200, 200, 300, 210)) //354度
	fmt.Printf("Angle:%f, %f\n ", Angle(200, 200, 300, 200), Angle(200, 200, 300, 200)) //360度
	fmt.Printf("Angle:%f\n ", Angle(200, 200, 200, 300))                                //0度
	fmt.Printf("Angle:%f\n ", Angle(200, 200, 270, 270))                                //270度
	fmt.Printf("Angle:%f\n ", Angle(200, 200, 300, 200))                                //360度

	x, y := CalculatePoint(200, 200, 0, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 300, 200), Angle(200, 200, 300, 200))
	x, y = CalculatePoint(200, 200, 45, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 270.710678, 129.289322), Angle(200, 200, 270.710678, 129.289322))
	x, y = CalculatePoint(200, 200, 60, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 250.000000, 113.397460), Angle(200, 200, 250.000000, 113.397460))
	x, y = CalculatePoint(200, 200, 90, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 200.000000, 100.000000), Angle(200, 200, 200.000000, 100.000000))
	x, y = CalculatePoint(200, 200, 100, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 182.635182, 101.519225), Angle(200, 200, 182.635182, 101.519225))
	x, y = CalculatePoint(200, 200, 150, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 113.397460, 150.000000), Angle(200, 200, 113.397460, 150.000000))
	x, y = CalculatePoint(200, 200, 180, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 100.000000, 200.000000), Angle(200, 200, 100.000000, 200.000000))
	x, y = CalculatePoint(200, 200, 200, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 106.030738, 234.202014), Angle(200, 200, 106.030738, 234.202014))
	x, y = CalculatePoint(200, 200, 260, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 182.635182, 298.480775), Angle(200, 200, 182.635182, 298.480775))
	x, y = CalculatePoint(200, 200, 270, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 200.000000, 300.000000), Angle(200, 200, 200.000000, 300.000000))
	x, y = CalculatePoint(200, 200, 271, 1000)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 217.452406, 1199.847695), Angle(200, 200, 217.452406, 1199.847695))
	x, y = CalculatePoint(200, 200, 272, 1000)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 234.899497, 1199.390827), Angle(200, 200, 234.899497, 1199.390827))
	x, y = CalculatePoint(200, 200, 300, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 250.000000, 286.602540), Angle(200, 200, 250.000000, 286.602540))
	x, y = CalculatePoint(200, 200, 350, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 298.480775, 217.364818), Angle(200, 200, 298.480775, 217.364818))
	x, y = CalculatePoint(200, 200, 360, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 300.000000, 200.000000), Angle(200, 200, 300.000000, 200.000000))
	x, y = CalculatePoint(200, 200, -90, 100)
	fmt.Printf("CalculatePoint:%f, %f, d:%v, a:%v\n ", x, y, Distance(200, 200, 200.000000, 300.000000), Angle(200, 200, 200.000000, 300.000000))

}
