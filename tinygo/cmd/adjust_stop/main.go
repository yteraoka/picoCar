// adjust_stop はサーボの停止値を求めるプログラムです (adjust_stop.py の移植)。
//
// サーボには個体差があり config.go の設定値を調整する必要があります。
// その基本となるのがサーボが停止する値で、これがズレていると直進性を出すのに苦労します。
//
// 左右のサーボそれぞれについて、パルス幅を 1450〜1550us まで 5us 刻みで変えるので、
// 止まった値と動き出した値を記録して、その中間値を停止値としてください。
// (デューティ値 = パルス幅[us] * 65535 / 20000。例: 1495us → 4899)
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/servo"
)

// 調べる車輪: 右車輪は servo.Right (GPIO0)、左車輪は servo.Left (GPIO1) を設定してください。
const wheel = servo.Right

func main() {
	servo.SetPulse(wheel, 1300*time.Microsecond)
	time.Sleep(time.Second)

	// 1450〜1550us を 5us 刻みで確認
	for us := 1450; us <= 1550; us += 5 {
		servo.SetPulse(wheel, time.Duration(us)*time.Microsecond)
		println("PWM =", us, "us  duty_u16 =", us*65535/20000)
		println("この値でサーボが止まっているか確認してください")
		time.Sleep(2 * time.Second)
	}

	// 最後は停止
	servo.SetPulse(wheel, 1500*time.Microsecond)
	println("終了")
	println("止まっていた、中間値を停止値としてください。")
	for {
		time.Sleep(time.Hour)
	}
}
