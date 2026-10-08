//go:build rp2350

package bootsw

import (
	"device/arm"
	"device/rp"
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"
)

// BOOTSEL スイッチはフラッシュメモリの CS (QSPI_SS) 端子につながっているため、
// 読み取る間はフラッシュにアクセスできません。pico-sdk の get_bootsel_button()
// と同じ手順を、割り込みを止めて RAM 上の関数で実行します
// (MicroPython の rp2.bootsel_button() 相当)。
const (
	qspiSSCtrl = 0x40030000 + 0x1c // IO_QSPI.GPIO_QSPI_SS_CTRL
	sioGPIOHi  = 0xd0000000 + 0x08 // SIO.GPIO_HI_IN
)

// Pressed は BOOTSEL スイッチが押されていれば true を返します。
func Pressed() bool {
	state := interrupt.Disable()
	pressed := readBootsel()
	interrupt.Restore(state)
	return pressed
}

// readBootsel はフラッシュを使わないよう RAM に配置します。
// この中からフラッシュ上の関数を呼んではいけません。
//
//go:section .ramfuncs
//go:noinline
func readBootsel() bool {
	ctrl := (*uint32)(unsafe.Pointer(uintptr(qspiSSCtrl)))
	hiIn := (*uint32)(unsafe.Pointer(uintptr(sioGPIOHi)))
	const mask = rp.IO_QSPI_GPIO_QSPI_SS_CTRL_OEOVER_Msk

	// CS の出力を無効にして、プルアップされた端子の状態を読めるようにする
	v := volatile.LoadUint32(ctrl)&^mask | rp.IO_QSPI_GPIO_QSPI_SS_CTRL_OEOVER_DISABLE<<rp.IO_QSPI_GPIO_QSPI_SS_CTRL_OEOVER_Pos
	volatile.StoreUint32(ctrl, v)
	for range 1000 {
		arm.Asm("nop")
	}
	// スイッチを押すと CS が Low になる
	pressed := volatile.LoadUint32(hiIn)&rp.SIO_GPIO_HI_IN_QSPI_CSN == 0

	v = volatile.LoadUint32(ctrl)&^mask | rp.IO_QSPI_GPIO_QSPI_SS_CTRL_OEOVER_NORMAL<<rp.IO_QSPI_GPIO_QSPI_SS_CTRL_OEOVER_Pos
	volatile.StoreUint32(ctrl, v)
	return pressed
}
