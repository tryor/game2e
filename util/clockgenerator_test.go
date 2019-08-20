package util

import (
	//"bytes"
	//	"fmt"
	"testing"
	"time"
	//"strings"
	"tryor/game2e/log"
)

func TestClockGenerator(t *testing.T) {
	clockGenerator := NewClockGenerator(time.Millisecond * 10)

	cb := func(clock uint64) {
		//go func() {
		log.Infof("clock:%v", clock)
		time.Sleep(time.Millisecond * 50)
		log.Infof("clock:%v", clock)
		//}()
	}
	clockGenerator.AddClockCallback(&cb)
	clockGenerator.Start()

	time.Sleep(time.Second * 2)
	clockGenerator.Stop()
}
