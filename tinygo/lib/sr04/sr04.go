// Package sr04 は超音波センサー HC-SR04 で距離を測定します (lib_SR04.py の移植)。
//
// TRIG: GPIO14, ECHO: GPIO15。測定中は GPIO17 の LED を点灯します。
package sr04

import (
	"machine"
	"time"

	"github.com/momorara/picoCar/tinygo/lib/led"
)

const (
	trig = machine.GP14
	echo = machine.GP15

	timeout = 30 * time.Millisecond
)

func init() {
	trig.Configure(machine.PinConfig{Mode: machine.PinOutput})
	echo.Configure(machine.PinConfig{Mode: machine.PinInput})
}

// Read は距離を cm で返します。測定に失敗した場合は 10cm 未満の値 (9) を返します。
func Read() int {
	time.Sleep(100 * time.Millisecond)
	led.On()
	defer led.Off()

	trig.Low()
	time.Sleep(2 * time.Microsecond)
	trig.High()
	time.Sleep(10 * time.Microsecond)
	trig.Low()

	d, ok := pulseHigh()
	if !ok {
		println("測定失敗1")
		return 9 // 測定失敗したら 10cm 以内の数字を返す
	}
	// 音速 343m/s = 0.0343cm/us を往復分で割る
	return int(d.Microseconds() * 343 / 20000)
}

// pulseHigh は ECHO が High になっている時間を測ります
// (MicroPython の time_pulse_us(ECHO, 1, 30000) 相当)。
func pulseHigh() (time.Duration, bool) {
	start := time.Now()
	for !echo.Get() {
		if time.Since(start) > timeout {
			return 0, false
		}
	}
	start = time.Now()
	for echo.Get() {
		if time.Since(start) > timeout {
			return 0, false
		}
	}
	return time.Since(start), true
}
