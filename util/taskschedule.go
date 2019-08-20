package util

import (
	//	"os"
	"sync"
	"tryor/game2e/log"
)

type ITask interface {
	Run(clock uint64)
}

type IStrategy interface {
	Executable(t ITask, clock uint64) bool
}

type TaskSchedule struct {
	ClockGenerator    *ClockGenerator
	tasks             map[ITask]*taskStrategy
	eachTaskStrategys []*taskStrategy
	locker            *sync.RWMutex
	isexecuterun      bool

	executecb func(clock uint64)
}

type taskStrategy struct {
	task     ITask
	strategy IStrategy
}

func NewTaskSchedule(clockGenerator *ClockGenerator) *TaskSchedule {
	ts := &TaskSchedule{ClockGenerator: clockGenerator,
		tasks:             make(map[ITask]*taskStrategy),
		eachTaskStrategys: make([]*taskStrategy, 0),
		locker:            new(sync.RWMutex)}

	ts.executecb = ts.execute
	clockGenerator.AddClockCallback(&ts.executecb)
	return ts
}

func (this *TaskSchedule) Destroy() {
	this.ClockGenerator.RemoveClockCallback(&this.executecb)
	//this.Clear()
	this.tasks = nil
}

/**
 * 使用指定的策略安排一个任务
 * @param t 将要被执行的任务
 * @param s 时间策略，根据策略返回状态来执行任务
 */
func (this *TaskSchedule) Schedule(t ITask, s IStrategy) {
	this.locker.Lock()
	defer this.locker.Unlock()
	if this.tasks == nil {
		log.Warn("TaskSchedule is Destroyed")
		return
	}
	this.tasks[t] = &taskStrategy{task: t, strategy: s}
}

func (this *TaskSchedule) TaskCount() int {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return len(this.tasks)
}

func (this *TaskSchedule) Remove(t ITask) {
	this.locker.Lock()
	defer this.locker.Unlock()
	delete(this.tasks, t)
}

func (this *TaskSchedule) Clear() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.tasks = make(map[ITask]*taskStrategy)
}

func (this *TaskSchedule) makeEachTaskStrategys() {
	this.locker.RLock()
	defer this.locker.RUnlock()
	this.eachTaskStrategys = this.eachTaskStrategys[0:0]
	for _, ts := range this.tasks {
		this.eachTaskStrategys = append(this.eachTaskStrategys, ts)
	}
}

func (this *TaskSchedule) execute(clock uint64) {
	//timespender := NewTimespender(fmt.Sprint("TaskSchedule.getTasks,", len(this.tasks)))
	this.makeEachTaskStrategys()
	//timespender.Print()
	//log.Infof("TaskCount:%v, %p, gid:%v, pid:%v, %v, %p", this.TaskCount(), this, os.Getgid(), os.Getpid(), os.Getegid(), this.ClockGenerator)
	for _, ts := range this.eachTaskStrategys {
		if ts.strategy.Executable(ts.task, clock) {
			ts.task.Run(clock)
		}
	}
	//timespender.Print(50)
}

func (this *TaskSchedule) execute33(clock uint64) {
	//timespender := canvas.NewTimespender(fmt.Sprint("TaskSchedule.getTasks,", len(this.tasks)))
	//this.isexecuterun = true

	for _, ts := range this.tasks {
		if ts.strategy.Executable(ts.task, clock) {
			ts.task.Run(clock)
		}
	}
	//this.isexecuterun = false

	//timespender.Print()
}
