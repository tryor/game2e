package util

import (
	"sync"
)

type runtimeTimer struct {
	when  uint64
	delay int
	f     func(*Timer, uint64)
	arg   interface{}
}

type Timer struct {
	//	C            <-chan uint64
	TaskSchedule *TaskSchedule
	r            runtimeTimer

	//	xx uint64

	active     bool
	delays     []int
	delayIndex int
	task       *timertask
	strategy   *timerstrategy
	locker     sync.RWMutex
}

func when(ts *TaskSchedule, delay int) uint64 {
	return ts.ClockGenerator.Clocks() + uint64(delay)
}

func NewTimer(ts *TaskSchedule, f func(delayIndex int), delay ...int) *Timer {
	if len(delay) < 1 {
		panic("the delay array must be greater than 1")
	}
	timer := &Timer{
		TaskSchedule: ts,
		r: runtimeTimer{
			when:  when(ts, delay[0]),
			delay: delay[0],
			f:     sendTime,
			arg:   f,
		},
		active: false,
		delays: delay,
	}
	timer.task = newTimertask(timer)
	timer.strategy = newTimerstrategy(timer)
	//	ts.Schedule(timer.task, timer.strategy)
	return timer
}

func AfterFunc(ts *TaskSchedule, delay int, f func()) *Timer {
	timer := &Timer{
		TaskSchedule: ts,
		r: runtimeTimer{
			when:  when(ts, delay),
			delay: delay,
			f:     goFunc,
			arg:   f,
		},
		active: true,
	}
	timer.task = newTimertask(timer)
	timer.strategy = newTimerstrategy(timer)
	ts.Schedule(timer.task, timer.strategy)
	return timer
}

func (this *Timer) IsActive() bool {
	return this.active
}

func (this *Timer) Stop() {
	if this.active {
		this.active = false
		this.TaskSchedule.Remove(this.task)
	}
}

//func (this *Timer) stop() {
//	this.active = false
//}

//delayIndex 如果设置此值，指示从delays[delayIndex]中开始
func (this *Timer) Start(delayIndex ...int) {
	if this.active {
		return
	}

	if len(delayIndex) > 0 {
		didx := 0
		if delayIndex[0] < 0 || delayIndex[0] >= len(this.delays) {
			//this.delayIndex = 0
		} else {
			//this.delayIndex = delayIndex[0]
			didx = delayIndex[0]
		}
		this.r.delay = this.delays[didx]
		this.delayIndex = didx
	}

	//	println("2 Timer.Start:", this.active, this.TaskSchedule.TaskCount())

	w := when(this.TaskSchedule, this.r.delay)
	this.locker.Lock()
	this.r.when = w
	this.locker.Unlock()

	//atomic.StoreUint64(&this.r.when, when(this.TaskSchedule, this.r.delay))
	//	w := this.r.when
	//	res := atomic.CompareAndSwapUint64(&this.r.when, w, when(this.TaskSchedule, this.r.delay))

	//	println("Timer.Start:", w, res)
	//	println("this.when:", this.r.when)
	//	atomic.StoreUint64(&this.when, 234)
	//	atomic.StoreUint64(&this.xx, 234)
	//	println("this.xx:", this.xx)
	//	println("this.xx:", this.xx)
	this.TaskSchedule.Schedule(this.task, this.strategy)
	this.active = true
}

func (this *Timer) Reset(delay int) {
	//	this.locker.Lock()
	//	defer this.locker.Unlock()
	//this.r.when = when(this.TaskSchedule, delay)
	w := when(this.TaskSchedule, delay)
	this.locker.Lock()
	this.r.when = w
	this.locker.Unlock()

	//	atomic.StoreUint64(&this.when, mwhen(this.TaskSchedule, delay))
	this.r.delay = delay
	if !this.active {
		this.TaskSchedule.Schedule(this.task, this.strategy)
	}
	this.active = true
}

//手动立即执行
func (this *Timer) Execute() {
	this.task.Run(this.TaskSchedule.ClockGenerator.Clocks())
}

func sendTime(timer *Timer, now uint64) {
	//	timer.locker.Lock()
	//	defer timer.locker.Unlock()

	didx := timer.delayIndex
	didx++
	if didx >= len(timer.delays) {
		didx = 0
	}
	//	println(timer.delayIndex, len(timer.delays))
	timer.delayIndex = didx
	timer.r.delay = timer.delays[didx]

	//	timer.r.when = when(timer.TaskSchedule, timer.r.delay)

	w := when(timer.TaskSchedule, timer.r.delay)
	timer.locker.Lock()
	timer.r.when = w
	timer.locker.Unlock()

	//	atomic.StoreUint64(&timer.r.when, mwhen(timer.TaskSchedule, timer.r.delay))
	timer.r.arg.(func(int))(timer.delayIndex)

}

func goFunc(timer *Timer, now uint64) {
	timer.Stop()
	timer.r.arg.(func())()
}

type timertask struct {
	timer *Timer
}

func newTimertask(timer *Timer) *timertask {
	return &timertask{timer: timer}
}

func (this *timertask) Run(clock uint64) {
	this.timer.r.f(this.timer, clock)
}

type timerstrategy struct {
	timer *Timer
}

func newTimerstrategy(timer *Timer) *timerstrategy {
	return &timerstrategy{timer: timer}
}

func (this *timerstrategy) Executable(t ITask, clock uint64) bool {
	//	this.timer.locker.RLock()
	//	defer this.timer.locker.RUnlock()

	this.timer.locker.RLock()
	w := this.timer.r.when
	this.timer.locker.RUnlock()

	//	println("timerstrategy.Executable", atomic.LoadUint64(&this.timer.when))
	//	if this.timer.active && atomic.LoadUint64(&this.timer.r.when) <= clock {
	if this.timer.active && w <= clock {
		return true
	}
	return false
}
