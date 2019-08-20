package ai

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var monster *BattleRole //怪兽
var player *BattleRole  //玩家
var sm *StateMachine

func init() {
	monster = &BattleRole{name: "怪物", bloodValue: 100, position: rand.Int() % 300}
	player = &BattleRole{name: "英雄", bloodValue: 100, position: rand.Int() % 300}
	sm = NewStateMachine()
}

//状态
var walk = &State{STATE_WALK}     //游走
var trace = &State{STATE_TRACE}   //追逐
var away = &State{STATE_AWAY}     //逃跑
var attack = &State{STATE_ATTACK} //攻击
var death = &State{STATE_DEATH}   //死亡

func Test_StateMachine(t *testing.T) {
	//初始化怪物(monster)状态机
	//加入状态，规则，行为
	sm.Add(walk, walk, ToWalkRule, func(from, to IState) { monster.walk() })   //游走(0)->游走(0)
	sm.Add(trace, walk, ToWalkRule, func(from, to IState) { monster.walk() })  //追逐(1)->游走(0)
	sm.Add(away, walk, ToWalkRule, func(from, to IState) { monster.walk() })   //逃跑(2)->游走(0)
	sm.Add(attack, walk, ToWalkRule, func(from, to IState) { monster.walk() }) //攻击(3)->游走(0)

	sm.Add(walk, trace, ToTraceRule, func(from, to IState) { monster.trace(player) })   //游走(0)=>追逐(1)
	sm.Add(trace, trace, ToTraceRule, func(from, to IState) { monster.trace(player) })  //追逐(1)=>追逐(1)
	sm.Add(away, trace, ToTraceRule, func(from, to IState) { monster.trace(player) })   //逃跑(2)=>追逐(1)
	sm.Add(attack, trace, ToTraceRule, func(from, to IState) { monster.trace(player) }) //攻击(3)=>>追逐(1)

	sm.Add(walk, attack, ToAttackRule, func(from, to IState) { monster.attack(player) })   //游走(0)=>攻击(3)
	sm.Add(trace, attack, ToAttackRule, func(from, to IState) { monster.attack(player) })  //追逐(1)=>攻击(3)
	sm.Add(attack, attack, ToAttackRule, func(from, to IState) { monster.attack(player) }) //攻击(3)=>攻击(3),
	sm.Add(away, attack, ToAttackRule, func(from, to IState) { monster.attack(player) })   //逃跑(2)=>攻击(3),

	sm.Add(walk, away, ToAwayRule, func(from, to IState) { monster.away(player) })   //游走(0)->逃跑(2)
	sm.Add(trace, away, ToAwayRule, func(from, to IState) { monster.away(player) })  //追逐(1)->逃跑(2)
	sm.Add(away, away, ToAwayRule, func(from, to IState) { monster.away(player) })   //逃跑(2)->逃跑(2)
	sm.Add(attack, away, ToAwayRule, func(from, to IState) { monster.away(player) }) //攻击(3)->逃跑(2)

	sm.Add(walk, death, ToDeathRule, func(from, to IState) { fmt.Println("the monster of dead!") })   //游走(0)->死亡(-1)
	sm.Add(trace, death, ToDeathRule, func(from, to IState) { fmt.Println("the monster of dead!") })  //追逐(1)->死亡(-1)
	sm.Add(away, death, ToDeathRule, func(from, to IState) { fmt.Println("the monster of dead!") })   //逃跑(2)->死亡(-1)
	sm.Add(attack, death, ToDeathRule, func(from, to IState) { fmt.Println("the monster of dead!") }) //攻击(3)->死亡(-1)

	fmt.Println("sm.GetStateCount():", sm.GetStateCount())
	fmt.Println("sm.GetCurrentState():", sm.GetCurrentState())
	fmt.Println("len(sm.fromStateMap):", len(sm.fromStateMap))
	fmt.Println("len(sm.toStateMap):", len(sm.toStateMap))

	player.bloodValue += -55
	player.position += 110

	sm.SetInitialState(walk)

	for i := 0; i < 100; i++ {
		currState := sm.GetCurrentState()
		fmt.Printf("currState:%v monster.getBloodValue():%d monster.getPosition():%d player.getBloodValue():%d player.getPosition():%d \n",
			currState, monster.bloodValue, monster.position,
			player.bloodValue, player.position)
		sm.Execute()
		if player.isDeath() {
			fmt.Printf("player is killed!\n")
			break
		}
		time.Sleep(time.Millisecond * 100)

	}

}

/** 状态常量定义 */
//游走(0), 追逐(1)，逃跑(2)，攻击(3)，死亡(-1)
//const STATE_WALK = 0   //游走
//const STATE_TRACE = 1  //追逐
//const STATE_AWAY = 2   //逃跑
//const STATE_ATTACK = 3 //攻击
//const STATE_DEATH = -1 //死亡

//type State struct {
//	v int
//}

//func (this *State) Equals(s IState) bool {
//	return this.v == s.(*State).v
//}

type BattleRole struct {
	name       string
	bloodValue int   /** 血量 */
	position   int   /** 位置 */
	state      State /** 状态: 游走(0), 追逐(1)，逃跑(2)，攻击(3)，死亡(-1) */
}

/** 游走
 */
func (this *BattleRole) walk() {
	//printf("position:%d", position);
	if !this.isDeath() {
		//移动位置
		if rand.Int()%2 == 1 {
			this.position += rand.Int() % 4
		} else {
			this.position -= rand.Int() % 4
		}
		if this.bloodValue < 100 {
			this.bloodValue += 1 //回血
		}
	}
}

/**
 * 追逐
 */
func (this *BattleRole) trace(r *BattleRole) {
	if !this.isDeath() {
		if this.position < r.position {
			this.position += 5
		} else {
			this.position -= 5
		}
	}
}

/**
 * 逃跑(2)
 */
func (this *BattleRole) away(r *BattleRole) {
	if !this.isDeath() {
		if this.position > r.position {
			this.position += 6
		} else {
			this.position -= 6
		}
	}
}

/**
 *  攻击
 */
func (this *BattleRole) attack(r *BattleRole) {
	if !this.isDeath() {
		r.bloodValue -= rand.Int() % 6
	}
}

/**
 * 被攻击
 */
func (this *BattleRole) attacked(val int) {
	if !this.isDeath() {
		this.bloodValue -= val
	}
}

/**
 * 是否已经死亡
 */
func (this *BattleRole) isDeath() bool {
	//return state == STATE_DEATH;
	//return sm->getState() == STATE_DEATH;
	return this.bloodValue <= 0
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

/**
 * 转到->游走(0)
 */
var ToWalkRule = func(from, to IState) bool {
	//血量>0 && 与玩家距离大于100
	distance := abs(monster.position - player.position)
	//printf("distance:%d", distance);
	return monster.bloodValue > 0 && distance > 100
}

/**
 * 转到->追逐(1)
 */
var ToTraceRule = func(from, to IState) bool {
	//血量>50 && 与玩家距离在10到100之间 && 玩家血量小于50
	distance := abs(monster.position - player.position)
	return monster.bloodValue > 50 && (distance >= 10 && distance <= 100) && player.bloodValue < 50
}

/**
 * 转到->攻击(3)
 */
var ToAttackRule = func(from, to IState) bool {
	//血量>50 && 与玩家距离在0到10之间  && 玩家血量小于50
	distance := abs(monster.position - player.position)
	return monster.bloodValue > 50 && (distance >= 0 && distance < 10) && player.bloodValue < 50
}

/**
 * 转到->逃跑(2)
 */
var ToAwayRule = func(from, to IState) bool {
	//血量<50 && 与玩家距离在小于100 && 玩家血量大于50
	distance := abs(monster.position - player.position)
	return monster.bloodValue <= 50 && distance < 100 && player.bloodValue >= 50
}

/**
 * 转到->死亡(-1)
 */
var ToDeathRule = func(from, to IState) bool {
	//血量<=0
	return monster.bloodValue <= 0
}
