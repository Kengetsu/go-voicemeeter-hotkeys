package main

import (
	"github.com/MakeNowJust/hotkey"
	"github.com/electricbubble/go-toast"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/onyx-and-iris/voicemeeter/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

var volumeUpIncrement = 3.0
var volumeDownIncrement = -3.0
var voicemeeterWaitTime = 30 * time.Second
var vm *voicemeeter.Remote

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	var err error
	vm, err = vmConnect()
	if err != nil {
		_ = toast.Push("Failed to connect to voice meeter. Voicemeeter Hotkey Closing.",
			toast.WithTitle("Voicemeeter Hotkey"),
			toast.WithAppID("Voicemeeter Hotkey"),
			toast.WithAudio(toast.Default),
			toast.WithLongDuration())
		panic("Failed to connect to voice meeter.")
	}

	systray.Run(onReady, onExit)
}

func vmConnect() (*voicemeeter.Remote, error) {
	for retry := 0; retry <= 3; retry++ {
		var err error
		if retry == 3 {
			return nil, err
		}
		vm, err = voicemeeter.NewRemote("banana", 5)
		err = vm.Login()
		if err != nil {
			_ = toast.Push("Failed to connect to voice meeter. Please make sure Voice Meeter is running.",
				toast.WithTitle("Voicemeeter Hotkey"),
				toast.WithAppID("Voicemeeter Hotkey"),
				toast.WithAudio(toast.Default),
				toast.WithShortDuration())
			time.Sleep(voicemeeterWaitTime)
			continue
		}
		break
	}
	return vm, nil
}
func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Voicemeeter Hotkey")
	systray.SetTooltip("Voicemeeter Hotkey")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	mQuit.SetIcon(icon.Data)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	//fmt.Printf(vm.Type())
	go registerHotkeys()
}
func onExit() {
	_ = vm.Logout()
}

func registerHotkeys() {
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
