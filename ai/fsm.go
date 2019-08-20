package ai

import (
	. "tryor/game2e/util"

	. "github.com/tryor/eui"

	//	"github.com/google/gxui"
	"github.com/tryor/commons/event"
)

type IState interface {
	Equals(s IState) bool
	String() string
}

type FAction func(from, to IState)
type FRule func(from, to IState) bool

type StateMachine struct {
	*event.Dispatcher

	allFinalStates []IState //所有结束状态
	allStates      []IState //所有状态
	//  typedef hash_map<T, std::vector<ToState*>*> ToStateMap; //用于存储后续状态的Map
	toStateMap map[IState][]*toState //用于存储后续状态的Map
	// typedef hash_map<T, std::vector<T>*> FromStateMap;  /** 定义用于存储前置状态的Map */
	fromStateMap map[IState][]IState

	transformFlag    bool //状态转换标志，如果至少转换过一次，此值为true,否则为false
	initialStateFlag bool //初始状态标志，如果没有设置初始状态，值为flase, 否则为true

	initialState IState /** 初始状态 */
	currState    IState /** 当前状态 */

	//定义事件
	onNofusEvent  Event //*NoFollowUpStateEvent
	onNorfusEvent Event //*NoRuleFollowUpStateEvent
	onFinalEvent  Event //*FinalStateEvent

	timer *Timer //用于定时执行状态机
}

//如果让状态机自动执行，speed为速率，单位为时调周期
func NewStateMachine(ts *TaskSchedule, speed int) *StateMachine {
	sm := &StateMachine{}
	sm.Dispatcher = event.NewDispatcher()
	sm.allStates = make([]IState, 0)
	sm.allFinalStates = make([]IState, 0)
	sm.toStateMap = make(map[IState][]*toState)
	sm.fromStateMap = make(map[IState][]IState)

	//	sm.nofusEvent = NewNoFollowUpStateEvent(sm)
	//	sm.norfusEvent = NewNoRuleFollowUpStateEvent(sm)
	//	sm.finalEvent = NewFinalStateEvent(sm)
	sm.onNofusEvent = CreateEvent(func() {})
	sm.onNorfusEvent = CreateEvent(func() {})
	sm.onFinalEvent = CreateEvent(func() {})

	sm.timer = NewTimer(ts, func(delayIndex int) {
		sm.Execute()
	}, speed)

	return sm
}

/**
 * 向状态机中加入状态
 *
 * 约束：
 *  1. 不能重复加入
 *  2. 前置状态不能为结束状态，不能从结束状态向其它状态转换
 *
 * @param from
 *          起始状态
 * @param to
 *          目标状态
 * @param r
 *          状态转换规则
 * @param a
 *          状态转换成功过后执行的操作
 */
func (this *StateMachine) Add(from, to IState, r FRule, a FAction) {
	/* 约束：2 */
	if this.isFinalState(from) {
		panic("from state is final state")
	}

	tos, ok := this.toStateMap[from]
	if ok {
		if this.existToState(tos, to) {
			panic("The state can not be repeated by adding") //约束：1 不能重复加入
		}
	} else {
		tos = make([]*toState, 0)
	}

	ts := &toState{}
	ts.to = to
	ts.rule = r
	ts.action = a
	tos = append(tos, ts)
	this.toStateMap[from] = tos

	/* 记录前置状态对 */
	fs, ok := this.fromStateMap[to]
	if !ok {
		fs = make([]IState, 0)
	}
	fs = append(fs, from)
	this.fromStateMap[to] = fs

	/* 记录所有状态 */
	if !this.existState(from) {
		this.allStates = append(this.allStates, from)
	}
	if !this.existState(to) {
		this.allStates = append(this.allStates, to)
	}

	return
}

/**
 * 移除状态
 */
func (this *StateMachine) Remove(from, to IState) {
	//移除后续状态
	tos, ok := this.toStateMap[from]
	if ok {
		for i, t := range tos {
			if t.to == to || t.to.Equals(to) {
				this.toStateMap[from] = append(tos[:i], tos[i+1:]...)
				break
			}
		}
	}

	//移除前置状态
	fs, ok := this.fromStateMap[to]
	if ok {
		for i, f := range fs {
			if f == from || f.Equals(from) {
				this.fromStateMap[to] = append(fs[:i], fs[i+1:]...)
				break
			}
		}
	}

}

func (this *StateMachine) Start() {
	if this.timer != nil {
		this.timer.Start()
	}
}

func (this *StateMachine) Stop() {
	if this.timer != nil {
		this.timer.Stop()
	}
}

/**
 * 尝试转换到下一状态,从当前状态开始执行，转换规则执行顺序以加入到状态机中的先后顺序执行，
 * 如果更先的状态的规则执行返回true，后面的状态规则将不再执行。
 *
 *
 * @return 如果成功转换到下一状态，返回true,当前状态变为下一状态， 否则返回false,当前状态不变
 */
func (this *StateMachine) Execute() bool {
	//如果没设置初始状态
	if !this.initialStateFlag {
		panic("Has not set the initial state") //还没有设置初始状态
	}
	//如果是第一次执行，设置当前状态为初始状态
	if !this.transformFlag {
		this.transformFlag = true
		this.currState = this.initialState
	}
	if this.isFinalState(this.currState) {
		//this.FireEvent(this.finalEvent)
		this.onFinalEvent.Fire()
		return false
	}

	tos, ok := this.toStateMap[this.currState]
	if !ok {
		//没有找到后续状态，返回false，并触发NoFollowUpStateEvent事件
		//this.FireEvent(this.nofusEvent)
		this.onNofusEvent.Fire()
		return false
	}

	from := this.currState
	for _, to := range tos {
		if to.rule(from, to.to) {
			this.currState = to.to
			to.action(from, to.to)
			return true
		}
	}

	//没有找到满足条件的后续状态，返回false，并触发NoRuleFollowUpStateEvent事件
	//this.FireEvent(this.norfusEvent)
	this.onNorfusEvent.Fire()
	return false
}

/**
 * 返回当前状态
 */
func (this *StateMachine) GetCurrentState() IState {
	return this.currState
}

/**
 * 设置初始状态
 */
func (this *StateMachine) SetInitialState(s IState) {
	this.initialState = s
	this.initialStateFlag = true
}

/**
 * 复位状态, 设置当前状态为初始状态
 */
func (this *StateMachine) Reset() {
	this.currState = this.initialState
}

/**
 * 设置结束状态
 *
 * 约束：
 *  1. 不能重复加入, 返回false
 *  2. 结束状态不能有后续状态, 返回false
 */
func (this *StateMachine) SetFinalState(fs IState) bool {
	if this.isFinalState(fs) || this.existToStates(fs) {
		return false
	}
	this.allFinalStates = append(this.allFinalStates, fs)
	return true
}

/**
 * 移除结束状态, 如果此结束状态不存在，返回false
 */
func (this *StateMachine) RemoveFinalState(fs IState) bool {
	for i, s := range this.allFinalStates {
		if s == fs || s.Equals(fs) {
			this.allFinalStates = append(this.allFinalStates[:i], this.allFinalStates[i+1:]...)
			return true
		}
	}
	return false
}

/**
 * 返回所有状态
 */
func (this *StateMachine) GetAllStates() []IState {
	return this.allStates
}

/**
 * 返回状态数量
 */
func (this *StateMachine) GetStateCount() int {
	return len(this.allStates)
}

/**
 * 返回后续状态
 */
func (this *StateMachine) GetToStates(from IState) []IState {
	ts := make([]IState, 0)
	tos, ok := this.toStateMap[from]
	if ok {
		for _, to := range tos {
			ts = append(ts, to.to)
		}
	}
	return ts
}

/**
 * 返回前置状态
 */
func (this *StateMachine) GetFromStates(to IState) []IState {
	if fs, ok := this.fromStateMap[to]; ok {
		return fs
	}
	return make([]IState, 0)
}

/**
 * 检查是否是结束状态
 */
func (this *StateMachine) IsFinalState(s IState) bool {
	return this.isFinalState(s)
}
func (this *StateMachine) isFinalState(s IState) bool {
	for _, fs := range this.allFinalStates {
		if s == fs || s.Equals(fs) {
			return true
		}
	}
	return false
}

/**
 * 检查状态是否存在
 */
func (this *StateMachine) ExistState(s IState) bool {
	return this.existState(s)
}
func (this *StateMachine) existState(s IState) bool {
	for _, as := range this.allStates {
		if as == s || as.Equals(s) {
			return true
		}
	}
	return false
	//        return find(allStates.begin(), allStates.end(), s) != allStates.end();
}

/**
 * 检查是否存在后续状态
 */
func (this *StateMachine) ExistToStates(s IState) bool {
	return this.existToStates(s)
}
func (this *StateMachine) existToStates(s IState) bool {
	tos, ok := this.toStateMap[s]
	return ok && this.existToState(tos, s)
}

/**
 * 检查状态是否已经在to状态列表中
 * @param tos 目标状态列表
 * @param s 状态
 * @return 如果状态存在，返回true
 */
func (this *StateMachine) ExistToState(tos []*toState, s IState) bool {
	return this.existToState(tos, s)
}
func (this *StateMachine) existToState(tos []*toState, s IState) bool {
	for _, to := range tos {
		if to.to == s || to.to.Equals(s) {
			return true
		}
	}
	return false
}

/**
 * 清除所有状态
 */
func (this *StateMachine) Clear() {
	//清除后置对象列表
	this.toStateMap = make(map[IState][]*toState)

	//清除前置对象列表
	this.fromStateMap = make(map[IState][]IState)

	//清除allStates
	this.allStates = this.allStates[0:0] //make([]IState, 0)

	//清除allFinalStates
	this.allFinalStates = this.allFinalStates[0:0] //make([]IState, 0)

	this.transformFlag = false
	this.initialStateFlag = false
	this.initialState = nil
	this.currState = nil
}

/** 用于包装后续状态 */
type toState struct {
	to     IState  //状态
	rule   FRule   //规则
	action FAction //行为
}

//func (this *toState) Equals(ts *toState) bool {
//	return (this.to == ts.to || this.to.Equals(ts.to)) && this.rule == ts.rule && this.action == ts.action
//}

////如果当前状态为结束状态，触发此事件
//type FinalStateEvent struct {
//	event.Event
//}

//func NewFinalStateEvent(source interface{}) *FinalStateEvent {
//	return &FinalStateEvent{Event: event.Event{Type: FMS_FINAL_STATE_EVENT_TYPE, Source: source}}
//}

////在进行状态转换时如果没有找到后续状态事件，将触发此事件
//type NoFollowUpStateEvent struct {
//	event.Event
//}

//func NewNoFollowUpStateEvent(source interface{}) *NoFollowUpStateEvent {
//	return &NoFollowUpStateEvent{Event: event.Event{Type: FMS_NO_FOLLOWUP_STATE_EVENT_TYPE, Source: source}}
//}

//type NoRuleFollowUpStateEvent struct {
//	event.Event
//}

//func NewNoRuleFollowUpStateEvent(source interface{}) *NoRuleFollowUpStateEvent {
//	return &NoRuleFollowUpStateEvent{Event: event.Event{Type: FMS_NO_RULEFOLLOWUP_STATE_EVENT_TYPE, Source: source}}
//}

//const FMS_FINAL_STATE_EVENT_TYPE event.Type = 10000
//const FMS_NO_FOLLOWUP_STATE_EVENT_TYPE event.Type = 10001
//const FMS_NO_RULEFOLLOWUP_STATE_EVENT_TYPE event.Type = 10002
