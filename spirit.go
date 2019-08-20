package game2e

import (
	"fmt"
	"sync/atomic"

	"tryor/game2e/log"
	. "tryor/game2e/util"

	//	"github.com/google/gxui"
	. "github.com/tryor/eui"
)

type ISpirit interface {
	IWidget

	//如果被攻击，攻击方调用此方法
	//attcksp 攻击自己的精灵
	Attacked(attcksp ISpirit)
	IsDeath() bool                                //是否已经死忘
	OnDeath(f func(sp ISpirit)) EventSubscription //当精灵死亡时事件

	GetMHP() int
	SetMHP(mhp int)
	GetHP() int
	SetHP(hp int)
	AddHP(hp int)
	GetMP() int
	GetATN() int
	SetATN(atn int)
	GetDEF() int
	SetDEF(def int)
	GetINT() int
	GetRES() int
	GetSPO() int
	GetEXP() int32
	AddEXP(exp int32) int32
	CalculateLevel() int32
	GetLevel() int32
	AddLevel(l int32) int32
	OnModifyAttrs(f func(ISpirit)) EventSubscription
}

type Spirit struct {
	*spirit

	mhp int //max health point 生命值
	hp  int //health point 当前生命值
	MP  int //megic point 魔法值,法术值，真气值
	atn int //物理攻击力
	def int //物理防御力
	INT int //魔法攻击力
	RES int //魔法防御力
	SPO int //移动增速度，比如: 标准移动速度为10, SPO为20, 最终速度为 10 + 10*(20/100) = 12

	exp   int32 //当前经验，默认
	level int32 //等级

	onModifyAttrsEvent Event //当精灵以上属性被修改时事件
}

func NewSpirit(ts *TaskSchedule, x, y int, spiritid string) *Spirit {
	nsp := newSpirit(ts, x, y, spiritid)
	if nsp == nil {
		return nil
	}
	sp := &Spirit{}
	sp.spirit = nsp

	sp.mhp = 100
	sp.hp = 100
	sp.MP = 30
	sp.atn = 20
	sp.def = 6
	sp.INT = 20
	sp.RES = 6
	sp.SPO = 0
	sp.exp = 100
	sp.level = 1

	return sp
}

func (this *Spirit) HashBytes() []byte {
	return []byte(fmt.Sprint(this))
}

func (this *Spirit) Equals(v2 interface{}) bool {
	return this == v2
}

func (this *Spirit) GetMHP() int {
	return this.mhp
}

func (this *Spirit) SetMHP(mhp int) {
	this.mhp = mhp
	this.fireModifyAttrsEvent()
}

func (this *Spirit) GetHP() int {
	return this.hp
}

func (this *Spirit) SetHP(hp int) {
	this.hp = hp
	this.fireModifyAttrsEvent()
}

func (this *Spirit) AddHP(hp int) {
	this.hp += hp
	this.fireModifyAttrsEvent()
}

func (this *Spirit) GetMP() int {
	return this.MP
}
func (this *Spirit) GetATN() int {
	return this.atn
}
func (this *Spirit) SetATN(atn int) {
	this.atn = atn
	this.fireModifyAttrsEvent()
}

func (this *Spirit) GetDEF() int {
	return this.def
}

func (this *Spirit) SetDEF(def int) {
	this.def = def
	this.fireModifyAttrsEvent()
}

func (this *Spirit) GetINT() int {
	return this.INT
}
func (this *Spirit) GetRES() int {
	return this.RES
}
func (this *Spirit) GetSPO() int {
	return this.SPO
}
func (this *Spirit) GetEXP() int32 {
	return atomic.LoadInt32(&this.exp)
}
func (this *Spirit) AddEXP(exp int32) int32 {
	exp = atomic.AddInt32(&this.exp, exp)
	this.fireModifyAttrsEvent()
	return exp
}

func (this *Spirit) CalculateLevel() int32 {
	return atomic.LoadInt32(&this.level)
}
func (this *Spirit) GetLevel() int32 {
	return atomic.LoadInt32(&this.level)
}
func (this *Spirit) AddLevel(l int32) int32 {
	lv := atomic.AddInt32(&this.level, l)
	this.fireModifyAttrsEvent()
	return lv
}
func (this *Spirit) OnModifyAttrs(f func(ISpirit)) EventSubscription {
	if this.onModifyAttrsEvent == nil {
		this.onModifyAttrsEvent = CreateEvent(func(ISpirit) {})
	}
	return this.onModifyAttrsEvent.Listen(f)
}
func (this *Spirit) fireModifyAttrsEvent() {
	if this.onModifyAttrsEvent != nil {
		this.onModifyAttrsEvent.Fire(this.Self)
	}
}

func (this *Spirit) Attacked(attcksp ISpirit) {
	panic("Not implemented")
}
func (this *Spirit) IsDeath() bool { //是否已经死忘
	panic("Not implemented")
}
func (this *Spirit) OnDeath(f func(sp ISpirit)) EventSubscription { //当精灵死亡时事件
	panic("Not implemented")
}

type spirit struct {
	*Widget
	Name         *Label
	spiritid     string
	movetimer    *Timer
	features     map[string][]*anglesAnimation //key为spirit.features.key
	stands       []*anglesAnimation            //默认动画, 一般为站立时动画
	moveds       []*anglesAnimation
	current      *Animation
	currentAngle float32

	PlayAnimationAtMoved bool
}

type anglesAnimation struct {
	Angles    [][2]float32
	Animation *Animation
}

func newSpirit(ts *TaskSchedule, x, y int, spiritid string) *spirit {
	sp := &spirit{Widget: NewWidget(), spiritid: spiritid}
	sp.Self = sp
	sp.features = make(map[string][]*anglesAnimation)

	spiritInfo := GetSpiritInfo(spiritid)
	if spiritInfo == nil {
		return nil
	}

	w, h := spiritInfo.Width, spiritInfo.Height //85, 87
	sp.SetWidth(w)
	sp.SetHeight(h)
	sp.SetCoordinate(x, y)

	for key, feature := range spiritInfo.Features {
		asas := make([]*anglesAnimation, 0)
		for _, f := range feature {
			asa := &anglesAnimation{Angles: f.Angles}
			asa.Animation = NewAnimation(ts, 0, 0, f.Animation, f.Frameindexs, true)
			if asa.Animation != nil {
				asa.Animation.SetFeature(&f.Feature)
			}
			asas = append(asas, asa)
		}
		sp.features[key] = asas
		if sp.stands == nil {
			sp.stands = asas
		}
	}

	if spiritInfo.Default != "" {
		sp.stands = sp.features[spiritInfo.Default]
	}

	if sp.stands != nil {
		sp.currentAngle = float32(270)
		animn := sp.findAnimationByAngle(sp.stands, sp.currentAngle)
		if animn != nil {
			sp.setAnimation(animn, false)
		}
	}

	if spiritInfo.Moved != "" {
		sp.moveds = sp.features[spiritInfo.Moved]
	}
	if sp.moveds == nil {
		sp.moveds = sp.stands
	}

	nameInfo := spiritInfo.Name
	sp.Name = NewLabel2(0, 0, 0, 0,
		nameInfo.Text, nameInfo.Font, nameInfo.Size, nameInfo.Color)
	sp.Name.FontStyle = nameInfo.Style
	sp.Name.Multiline = nameInfo.Multiline
	sp.Name.TextAlignment = nameInfo.Textalignment
	sp.Name.SetFeature(&nameInfo.Feature)
	sp.AddChild(sp.Name)

	sp.OnMoved(func(me *MovedEvent) {
		sp.onMovedEvent(me)
	})

	return sp
}

func (this *spirit) Destroy() {
	//	if this.current != nil {
	//		this.Self.(IWidget).RemoveChild(this.current)
	//	}
	for _, feature := range this.features {
		for _, f := range feature {
			if f.Animation != nil {
				this.Self.(IWidget).RemoveChild(f.Animation)
				f.Animation.Destroy()
			}
		}
	}
	this.Widget.Destroy()

	this.current = nil

	//	for _, feature := range this.features {
	//		for _, f := range feature {
	//			if f.Animation != nil {
	//				f.Animation.Destroy()
	//			}
	//		}
	//	}

}

func (this *spirit) GetCurrentAngle() float32 {
	return this.currentAngle
}

func (this *spirit) SetCurrentAngle(angle float32) {
	this.currentAngle = angle
}

func (this *spirit) SetMoving(b bool) {
	this.Widget.SetMoving(b)
	if !b && this.PlayAnimationAtMoved {
		this.PlayStand()
	}
}

//根据角度，返回相应动画
func (this *spirit) findAnimationByAngle(anglesAnims []*anglesAnimation, angle float32) *Animation {
	var anims *Animation
	for _, aa := range anglesAnims {
		for _, a := range aa.Angles {
			if angle >= a[0] && angle <= a[1] {
				anims = aa.Animation
				break
			}
		}
	}
	if anims == nil {
		log.Error("animation not find, angle is ", angle, ", spiritid is ", this.spiritid)
	}
	return anims
}

//resetFrame表示是否复位动画帧到开始帧
func (this *spirit) setAnimation(a *Animation, resetFrame bool) {

	current := this.current
	if current == nil {
		current = a
		if !this.Self.ExistChild(current.GetId()) {
			this.Self.(IElement).AddChild(current)
		}
	} else if current != a || resetFrame {
		current.Stop()
		//this.Self.(IElement).RemoveChild(current)
		current.SetVisible(false)
		current = a
		//this.Self.(IElement).AddChild(current)
		if !this.Self.ExistChild(current.GetId()) {
			this.Self.(IElement).AddChild(current)
		}
	}
	a.SetVisible(true)

	if resetFrame {
		current.Start(0)
	} else {
		current.Start()
	}
	current.Self.(IElement).SetModified(true)
	this.current = current
}

func (this *spirit) onMovedEvent(e *MovedEvent) {
	//	this.Widget.OnMovedEvent(e)
	//	println("spirit.OnMovedEvent", e.Dx, e.Dy, e.Angle)

	if this.PlayAnimationAtMoved {
		angle := e.Angle // float64(-1)
		if angle < 0 || angle > 360 {
			x, y := float64(this.ReferencePointX()), float64(this.ReferencePointY())
			angle = float32(Angle(x, y, x+float64(e.Dx), y+float64(e.Dy)))
		}
		if this.moveds != nil {
			walkanimn := this.findAnimationByAngle(this.moveds, angle)
			if walkanimn != nil {
				this.setAnimation(walkanimn, true)
			}
		}
		this.currentAngle = angle
	}

}

//播放站立状态动画
func (this *spirit) PlayStand(angle ...float32) {
	if len(angle) > 0 {
		this.currentAngle = angle[0]
	}
	if this.stands != nil {
		animn := this.findAnimationByAngle(this.stands, this.currentAngle)
		if animn != nil {
			this.setAnimation(animn, false)
		}
	}
}

//播放移动状态动画
//angle 移动角度
func (this *spirit) PlayMoved(angle ...float32) {
	if len(angle) > 0 {
		this.currentAngle = angle[0]
	}
	if this.stands != nil {
		animn := this.findAnimationByAngle(this.moveds, this.currentAngle)
		if animn != nil {
			this.setAnimation(animn, false)
		}
	}
}

//播放指定动画
//featurekey 动画Key
//loop 是否循环播放
//resetFrame 是否从第一帧开始播放
func (this *spirit) PlayAnimation(featurekey string, loop bool, resetFrame bool) *Animation {
	feature := this.features[featurekey]
	if feature != nil {
		animn := this.findAnimationByAngle(feature, this.currentAngle)
		if animn != nil {
			animn.Loop = loop
			this.setAnimation(animn, resetFrame)
			return animn
		}
	}
	return nil
}
