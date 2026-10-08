// Package picoled は Pico 2 W 基板上の LED を点滅させます (lib_LED_pico.py の移植)。
//
// Pico 2 W の基板上 LED は RP2350 の GPIO25 ではなく、無線チップ CYW43439 の
// GPIO0 (WL_GPIO0) に接続されているため、cyw43439 ドライバ経由で制御します。
// 無線チップの初期化には 1〜2 秒かかるので、最初に使うときに初期化します。
package picoled

import (
	"sync"
	"time"

	"github.com/soypat/cyw43439"
)

var (
	once sync.Once
	dev  *cyw43439.Device
)

// Configure は無線チップを初期化します。初期化に失敗しても LED が点かないだけで、
// 車の動作は続けられるようにエラーは表示のみとします。
func Configure() {
	once.Do(func() {
		d := cyw43439.NewPicoWDevice()
		if err := d.Init(cyw43439.DefaultWifiConfig()); err != nil {
			println("picoled: CYW43439 の初期化に失敗:", err.Error())
			return
		}
		dev = d
	})
}

// Set は LED を点灯 (true) / 消灯 (false) します。
func Set(on bool) {
	Configure()
	if dev == nil {
		return
	}
	if err := dev.GPIOSet(0, on); err != nil {
		println("picoled:", err.Error())
	}
}

func On()  { Set(true) }
func Off() { Set(false) }

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
