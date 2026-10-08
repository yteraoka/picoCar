// Package led は GPIO17 に接続された LED を点滅させます (lib_LED.py の移植)。
package led

import (
	"machine"
	"time"
)

const pin = machine.GP17

func init() {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func On()  { pin.High() }
func Off() { pin.Low() }

// Blink は num 回点滅させます。
func Blink(num int) {
	for range num {
		On()
		time.Sleep(200 * time.Millisecond)
		Off()
		time.Sleep(500 * time.Millisecond)
	}
}

// End は終了時の点滅パターンです。
func End() {
	for range 4 {
		On()
		time.Sleep(500 * time.Millisecond)
		Off()
		time.Sleep(100 * time.Millisecond)
	}
	On()
	time.Sleep(3 * time.Second)
	Off()
}
