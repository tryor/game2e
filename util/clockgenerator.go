package util

import (
	"runtime"
	"sync/atomic"
	"time"
	"tryor/game2e/log"
)

//type ClockCallback func(clock uint64)

type ClockGenerator struct {
	clockCount uint64
	interval   time.Duration
	//ticker         *time.Ticker
	tickerDone     chan struct{}
	tickerRunning  bool
	clockCallbacks map[*func(clock uint64)]*clockCallbackItem
}

type clockCallbackItem struct {
	running int32 //1 or 0
}

func (this *clockCallbackItem) IsRunning() bool {
	//	running := atomic.LoadInt32(&this.running)
	//	if running > 0 {
	//		println("clockCallbackItem.IsRunning()", atomic.LoadInt32(&this.running))
	//	}
	return atomic.LoadInt32(&this.running) > 0
}

func NewClockGenerator(period time.Duration) *ClockGenerator {
	cg := &ClockGenerator{interval: period, clockCallbacks: make(map[*func(clock uint64)]*clockCallbackItem)}
	return cg
}

func (this *ClockGenerator) SetPeriod(period time.Duration) {
	if this.interval != period {
		this.interval = period
		if this.IsRunning() {
			this.Stop()
			for this.IsRunning() {
				time.Sleep(time.Millisecond)
				log.Info(this.IsRunning())
			}
			this.Start()
		}
	}
}

func (this *ClockGenerator) AddClockCallback(cb *func(clock uint64)) {
	this.clockCallbacks[cb] = &clockCallbackItem{}
}

func (this *ClockGenerator) RemoveClockCallback(cb *func(clock uint64)) {
	//	this.clockCallbacks = append(this.clockCallbacks, cb)
	delete(this.clockCallbacks, cb)
}

func (this *ClockGenerator) IsRunning() bool {
	if this.tickerRunning {
		return true
	}
	//println("ClockGenerator.IsRunning():", this.clockCallbacks)
	for _, item := range this.clockCallbacks {
		if item.IsRunning() {
			return true
		}
	}
	return false
}

func (this *ClockGenerator) Clocks() uint64 {
	return atomic.LoadUint64(&this.clockCount)
}

func (this *ClockGenerator) startTicker(f func()) chan struct{} {
	tickerDone := make(chan struct{}, 1)
	go func() {
		timer := time.NewTicker(this.interval)
		defer timer.Stop()
		for {
			//log.Info("this.ticker.C ", len(timer.C), this.interval)
			select {
			case <-timer.C:
				f()
			case <-tickerDone:
				return
			}
		}
	}()
	return tickerDone
}

func (this *ClockGenerator) Start() {
	if this.tickerRunning {
		return
	}
	this.tickerRunning = true
	this.tickerDone = this.startTicker(func() {
		atomic.AddUint64(&this.clockCount, 1)
		//log.Infof("this.clockCallbacks.size:%v", len(this.clockCallbacks))
		for cb, cbi := range this.clockCallbacks {
			if !cbi.IsRunning() {
				//go this.callCallback(cb, cbi)
				this.callCallback(cb, cbi)
				//} else {
				//	log.Info("ClockGenerator.cbi is running ")
			}
		}
	})
}

//func (this *ClockGenerator) Start222() {
//	if this.ticker != nil && this.tickerRunning {
//		return
//	}
//	this.tickerRunning = true
//	this.ticker = time.NewTicker(this.interval)
//	go func() {
//		for this.tickerRunning {
//			select {
//			case _, ok := <-this.ticker.C:
//				//log.Info("this.ticker.C", len(this.ticker.C), ok)
//				if ok {
//					atomic.AddUint64(&this.clockCount, 1)
//					for cb, cbi := range this.clockCallbacks {
//						if !cbi.IsRunning() {
//							go this.callCallback(cb, cbi)
//						} else {
//							//log.Info("ClockGenerator.cbi is running ")
//						}
//					}
//				} else {
//					log.Info("1 this.tickerRunning:", this.tickerRunning)
//					return
//				}
//			}
//		}
//		log.Info("2 this.tickerRunning:", this.tickerRunning)
//	}()
//}

func (this *ClockGenerator) callCallback(cb *func(clock uint64), cbi *clockCallbackItem) {
	atomic.StoreInt32(&cbi.running, 1)

	defer func() {
		atomic.StoreInt32(&cbi.running, 0)

		if r := recover(); r != nil {
			log.Error("ClockGenerator.callCallback, Runtime error caught: %v", r)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				log.Info(file, line)
			}
		}
	}()
	(*cb)(atomic.LoadUint64(&this.clockCount))
}

func (this *ClockGenerator) Stop() {
	if !this.tickerRunning {
		return
	}
	this.tickerRunning = false
	//	if this.ticker != nil {
	//		this.ticker.Stop()
	//		this.ticker = nil
	//	}
	close(this.tickerDone)
}
