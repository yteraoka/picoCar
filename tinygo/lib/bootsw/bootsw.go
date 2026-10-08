// Package bootsw は基板上 LED を点滅させながら BOOTSEL スイッチを監視します
// (lib_bootSW.py の移植)。
//
// 元の Python 版は GPIO25 の LED を PWM で明滅させていましたが、Pico 2 W の
// LED は無線チップ経由で PWM できないため、ほぼ同じ周期の点滅にしています。
package bootsw

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/cds"
	"github.com/momorara/picoCar/tinygo/lib/picoled"
)

const (
	pollInterval = 5 * time.Millisecond
	onTime       = 410 * time.Millisecond // Python 版で明るくしていく時間
	offTime      = 330 * time.Millisecond // Python 版で暗くしていく時間
)

// Wait は BOOTSEL スイッチが押されるまで待ちます。
func Wait() {
	wait(false)
}

// WaitOrBright は BOOTSEL スイッチが押されるか、呼び出し時より Cds が
// 30% 以上明るくなるまで待ちます (Python 版の SW(1) 相当)。
func WaitOrBright() {
	wait(true)
}

func wait(useCds bool) {
	var cdsLevel int
	if useCds {
		cdsLevel = min(cds.Read(true)*13/10, 970)
	}
	for {
		if blinkUntilPressed(true, onTime) || blinkUntilPressed(false, offTime) {
			break
		}
		// Cds が 30% 以上明るくなった
		if useCds && cdsLevel < cds.Read(true) {
			break
		}
	}
	picoled.Off()
}

// blinkUntilPressed は LED を on にして d の間 BOOTSEL を監視し、押されたら true を返します。
func blinkUntilPressed(on bool, d time.Duration) bool {
	picoled.Set(on)
	for start := time.Now(); time.Since(start) < d; {
		if Pressed() {
			return true
		}
		time.Sleep(pollInterval)
	}
	return false
}
