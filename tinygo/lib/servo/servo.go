// Package servo は左右の車輪 (連続回転サーボ SG90-360) を制御します
// (lib_servo.py の移植)。
//
// 右車輪: GPIO0 (PWM0 A), 左車輪: GPIO1 (PWM0 B)。
package servo

import (
	"machine"
	"time"

	"github.com/momorara/picoCar/tinygo/config"
	"github.com/momorara/picoCar/tinygo/lib/sr04"
)

// Wheel は左右どちらの車輪かを表します。
type Wheel int

const (
	Right Wheel = iota // GPIO0
	Left               // GPIO1
)

const period = 20 * time.Millisecond // 50Hz

var (
	pwm      = machine.PWM0
	channels [2]uint8
)

func init() {
	if err := pwm.Configure(machine.PWMConfig{Period: uint64(period)}); err != nil {
		println("servo: PWM の設定に失敗:", err.Error())
	}
	channels[Right], _ = pwm.Channel(machine.GP0)
	channels[Left], _ = pwm.Channel(machine.GP1)
}

// SetDuty は MicroPython の duty_u16() と同じ 0〜65535 のスケールでデューティを設定します。
func SetDuty(w Wheel, duty uint16) {
	pwm.Set(channels[w], uint32(uint64(pwm.Top())*uint64(duty)/65535))
}

// SetPulse はパルス幅を指定します (MicroPython の duty_ns() 相当)。
func SetPulse(w Wheel, width time.Duration) {
	pwm.Set(channels[w], uint32(uint64(pwm.Top())*uint64(width)/uint64(period)))
}

func SetRight(duty uint16) { SetDuty(Right, duty) }
func SetLeft(duty uint16)  { SetDuty(Left, duty) }

// Run は前進を始めます。
func Run() {
	println("前進", config.LeftDuty.Forward, config.RightDuty.Forward)
	SetRight(config.RightDuty.Forward)
	SetLeft(config.LeftDuty.Forward)
}

// Back は後退を始めます。
func Back() {
	SetRight(config.RightDuty.Back)
	SetLeft(config.LeftDuty.Back)
}

// Stop は停止します。
func Stop() {
	println("停止")
	SetRight(config.RightDuty.Stop)
	SetLeft(config.LeftDuty.Stop)
}

// RotateLeft はその場で deg 度 (最大 180) 左回転します。
func RotateLeft(deg int) {
	deg = min(deg, 180)
	SetRight(config.RightDuty.Forward)
	SetLeft(config.LeftDuty.Back)
	time.Sleep(config.LeftRotateTime * time.Duration(deg) / 180)
	stopQuiet()
}

// RotateRight はその場で deg 度 (最大 180) 右回転します。
func RotateRight(deg int) {
	deg = min(deg, 180)
	SetRight(config.RightDuty.Back)
	SetLeft(config.LeftDuty.Forward)
	time.Sleep(config.RightRotateTime * time.Duration(deg) / 180)
	stopQuiet()
}

func stopQuiet() {
	SetRight(config.RightDuty.Stop)
	SetLeft(config.LeftDuty.Stop)
}

// Faraway は deg 度ずつ (20 度以上) stopTime 停止しながら 1 周し、
// 距離を測定して一番遠い方向を向きます。回る向きはランダムに決めます。
func Faraway(deg int, stopTime time.Duration) {
	left := randomBit() == 0
	direction := "rigth"
	if left {
		direction = "left"
	}
	rotate := func(left bool, d int) {
		if left {
			RotateLeft(d)
			println(d, "left_lotate")
		} else {
			RotateRight(d)
			println(d, "rigth_lotate")
		}
	}

	var dist []int
	for a := 0; a < 360; a += deg {
		println(a)
		if left {
			RotateLeft(deg)
		} else {
			RotateRight(deg)
		}
		dist = append(dist, sr04.Read())
		time.Sleep(stopTime)
	}
	// 一番大きい数字の位置を取得
	pos := 0
	for i, d := range dist {
		if d > dist[pos] {
			pos = i
		}
	}
	// 一番遠い角度 (元位置の一歩手前からなので、その分多い)
	farawayDeg := deg*pos + deg

	println(direction)
	print("[")
	for i, d := range dist {
		if i > 0 {
			print(", ")
		}
		print(d)
	}
	println("]")
	println(pos)
	println(farawayDeg)

	switch {
	case farawayDeg >= 360: // 今向いている方向が最大なので向きを変えない
		println(0)
	case farawayDeg <= 180: // 180 度以内なら同じ方向へ
		rotate(left, farawayDeg)
	default:
		rotate(!left, 360-farawayDeg)
	}
}

func randomBit() uint32 {
	r, err := machine.GetRNG()
	if err != nil {
		r = uint32(time.Now().UnixNano())
	}
	return r & 1
}
