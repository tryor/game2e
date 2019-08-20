package win

import (
	"runtime"
	//	"time"
	. "tryor/game2e"
	//	"tryor/game2e/log"

	. "github.com/tryor/winapi"
)

func Run(appRoutine func(driver IDriver)) {
	var driverFuncs DriverFuncs
	driverFuncs.Init = win32Init
	driverFuncs.Terminate = win32Terminate
	driverFuncs.WaitEvents = win32WaitEvents
	driverFuncs.Wake = win32Wake
	driverFuncs.CreateViewport = createWin32Viewport
	startDriver(&driverFuncs, appRoutine)

}

type driver struct {
	driverFuncs *DriverFuncs
}

func (d *driver) CreateViewport(width, height int, name string) (IViewport, error) {
	return d.driverFuncs.CreateViewport(width, height, name)
}

func (d *driver) Terminate() {
}

func startDriver(driverFuncs *DriverFuncs, appRoutine func(driver IDriver)) {
	if runtime.GOMAXPROCS(-1) < 2 {
		runtime.GOMAXPROCS(2)
	}

	if err := driverFuncs.Init(); err != nil {
		panic(err)
	}
	defer driverFuncs.Terminate()

	driver := &driver{driverFuncs: driverFuncs}
	appRoutine(driver)

	var m Msg
	for {
		r, e := GetMessage(&m, 0, 0, 0)
		if e != nil {
			panic(e)
		}
		if r == 0 {
			break
		}
		TranslateMessage(&m)
		DispatchMessageW(&m)
	}

}

func win32Init() error {
	//log.Debug("win32Init")
	return nil
}

func win32Terminate() {
	//log.Debug("win32Terminate")
}

func win32WaitEvents() {
	//log.Debug("win32WaitEvents")
	////time.Sleep(time.Second)
}

func win32Wake() {
	//log.Debug("win32Wake")
}
