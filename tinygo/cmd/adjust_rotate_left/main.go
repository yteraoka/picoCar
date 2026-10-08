// adjust_rotate_left は左回転で 180 度回転する時間を調整するプログラムです
// (adjust_rotate_left.py の移植)。
//
// PWM は直進時の値を流用し、左右それぞれ前進・後退の値で回転します。
// BOOTSEL スイッチを押すと、mode ピン (1〜4) に対応した
// config.AdjustmentRotateTime の時間だけ左回転します。
// ちょうど 180 度回転した時間を config.LeftRotateTime に設定してください。
//
// 詳しくは 説明資料/左右回転の調整.txt を参照。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/config"
	"github.com/momorara/picoCar/tinygo/lib/bootsw"
	"github.com/momorara/picoCar/tinygo/lib/mode"
	"github.com/momorara/picoCar/tinygo/lib/picoled"
	"github.com/momorara/picoCar/tinygo/lib/servo"
)

func main() {
	r, l := config.RightDuty, config.LeftDuty
	println(r.Forward, r.Back, r.Stop)
	println(l.Forward, l.Back, l.Stop)

	adj := config.AdjustmentRotateTime
	for _, d := range adj {
		print(d.String(), " ")
	}
	println()

	for {
		time.Sleep(500 * time.Millisecond)
		bootsw.Wait() // BOOTSEL スイッチが押されたら実行
		picoled.Blink(2)

		servo.Stop()

		// mode ターミナルブロックの位置に対応した時間回転する
		// (mode プラグなしの場合は Python 版と同じく最後の値を使う)
		i := (mode.Pin() + len(adj) - 1) % len(adj)
		println("左回転", l.Back, r.Forward)
		servo.SetRight(r.Forward)
		servo.SetLeft(l.Back)
		println(adj[i].String())
		time.Sleep(adj[i])

		servo.Stop()
	}
}
