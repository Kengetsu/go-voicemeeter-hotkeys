package main

import (
	"fmt"
	"github.com/MakeNowJust/hotkey"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/onyx-and-iris/voicemeeter/v2"
	log "github.com/sirupsen/logrus"
)

var volumeUpIncrement = 3.0
var volumeDownIncrement = -3.0

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	systray.Run(onReady, onExit)
}

func vmConnect() (*voicemeeter.Remote, error) {
	vm, err := voicemeeter.NewRemote("potato", 0)
	if err != nil {
		return nil, err
	}

	err = vm.Login()
	if err != nil {
		return nil, err
	}

	return vm, nil
}
func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Voicemeeter Hotkey")
	systray.SetTooltip("Voicemeeter Hotkey")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	mQuit.SetIcon(icon.Data)

	vm, err := vmConnect()
	if err != nil {
		panic(err)
	}
	go func() {
		<-mQuit.ClickedCh
		err := vm.Logout()
		if err != nil {
			return
		}
		systray.Quit()
	}()
	fmt.Printf(vm.Type())
	go registerHotkeys(vm)
}
func onExit() {
}

func registerHotkeys(vm *voicemeeter.Remote) {
	hkey := hotkey.New()
	//Toggle Audio Input from bus 0 and 1
	_, err := hkey.Register(hotkey.None, hotkey.F23, func() {
		//fmt.Printf("F23 Pressed")
		primaryBusMuted := vm.Bus[0].Mute()
		if primaryBusMuted {
			vm.Bus[0].SetMute(false)
			vm.Bus[1].SetMute(true)
		} else {
			vm.Bus[0].SetMute(true)
			vm.Bus[1].SetMute(false)
		}
	})
	if err != nil {
		return
	}
	_, err = hkey.Register(hotkey.None, hotkey.VOLUME_UP, func() {
		currentBus := vm.Bus[0]
		primaryBusMuted := currentBus.Mute()

		if primaryBusMuted {
			currentBus = vm.Bus[1]
		}
		//fmt.Printf(currentBus.String())
		VolumeChange(currentBus, volumeUpIncrement)
	})
	if err != nil {
		return
	}
	_, err = hkey.Register(hotkey.None, hotkey.VOLUME_DOWN, func() {
		currentBus := vm.Bus[0]
		primaryBusMuted := currentBus.Mute()

		if primaryBusMuted {
			currentBus = vm.Bus[1]
		}
		VolumeChange(currentBus, volumeDownIncrement)
	})
	if err != nil {
		return
	}
	_, err = hkey.Register(hotkey.Ctrl, hotkey.VOLUME_UP, func() {
		currentStrip := vm.Strip[3]
		VolumeChange(currentStrip, volumeUpIncrement)
	})
	if err != nil {
		return
	}

	_, err = hkey.Register(hotkey.Ctrl, hotkey.VOLUME_DOWN, func() {
		currentStrip := vm.Strip[3]
		VolumeChange(currentStrip, volumeDownIncrement)
	})
	if err != nil {
		return
	}
	_, err = hkey.Register(hotkey.Alt, hotkey.VOLUME_UP, func() {
		currentStrip := vm.Strip[4]
		VolumeChange(currentStrip, volumeUpIncrement)
	})
	if err != nil {
		return
	}

	_, err = hkey.Register(hotkey.Alt, hotkey.VOLUME_DOWN, func() {
		currentStrip := vm.Strip[4]
		VolumeChange(currentStrip, volumeDownIncrement)
	})
	if err != nil {
		return
	}
}

type bus interface {
	Gain() float64
	SetGain(v float64)
}

func VolumeChange(device bus, gainIncrement float64) {
	currentGain := device.Gain()
	if currentGain+gainIncrement > 0 {
		device.SetGain(0)
	} else {
		device.SetGain(currentGain + gainIncrement)
	}
}
