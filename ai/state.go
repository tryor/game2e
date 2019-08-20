package ai

import (
	"fmt"
)

const (
	STATE_STAND  = iota //站立
	STATE_WALK          //行走
	STATE_JUMP          //跳
	STATE_FALL          //落下
	STATE_TRACE         //追踪
	STATE_AWAY          //逃跑
	STATE_ATTACK        //攻击
	STATE_DEATH         //死亡
)

type State struct {
	v int
}

func (this *State) Equals(s IState) bool {
	return this == s || this.v == s.(*State).v
}

func (this *State) String() string {
	text := ""
	switch this.v {
	case STATE_STAND:
		text = "Stand"
	case STATE_WALK:
		text = "Walk"
	case STATE_JUMP:
		text = "Jump"
	case STATE_FALL:
		text = "Fall"
	case STATE_TRACE:
		text = "Trace"
	case STATE_AWAY:
		text = "Away"
	case STATE_ATTACK:
		text = "Attack"
	case STATE_DEATH:
		text = "Death"
	}
	return fmt.Sprintf("State:%v(%v)", text, this.v)
}

//状态
var Stand = &State{STATE_STAND}   //站立
var Walk = &State{STATE_WALK}     //游走
var Jump = &State{STATE_JUMP}     //跳
var Fall = &State{STATE_FALL}     //落下
var Trace = &State{STATE_TRACE}   //追
var Away = &State{STATE_AWAY}     //逃跑
var Attack = &State{STATE_ATTACK} //攻击
var Death = &State{STATE_DEATH}   //死亡

//全局精灵状态
var SpiritStates map[string]IState //Key为状态别名

func init() {
	SpiritStates = map[string]IState{
		"stand":  Stand,
		"walk":   Walk,
		"Jump":   Jump,
		"Fall":   Fall,
		"trace":  Trace,
		"away":   Away,
		"attack": Attack,
		"death":  Death,
	}

}
