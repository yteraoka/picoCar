// Package wifi は Pico 2 W の無線チップ (CYW43439) で無線 LAN に接続し、
// TCP/IP スタック (soypat/lneto) を動かします。
//
// cyw43439 リポジトリの examples/cywnet を元に、必要な部分だけにしたものです。
package wifi

import (
	"errors"
	"log/slog"
	"machine"
	"net/netip"
	"time"

	"github.com/soypat/cyw43439"
	"github.com/soypat/lneto/x/xnet"
)

// Config は接続設定です。
type Config struct {
	SSID     string
	Password string // 空ならパスワードなしのネットワークに接続
	Hostname string // DHCP でルーターに伝えるホスト名
	// RequestedAddr は DHCP で希望する IP アドレスです (ルーターが応じるとは限りません)。
	RequestedAddr netip.Addr
	// MaxTCPPorts は同時に使う TCP ポート (リスナー) の数です。
	MaxTCPPorts uint16
}

// Network は接続済みの無線 LAN です。
type Network struct {
	Dev   *cyw43439.Device
	Stack *xnet.StackAsync
	Addr  netip.Addr // 割り当てられた IP アドレス

	sendbuf []byte
}

const pollSleep = 5 * time.Millisecond

// Connect は無線チップを初期化して cfg.SSID に接続し、DHCP でアドレスを取得します。
// パケットの送受信は内部のゴルーチンで行います。
func Connect(cfg Config) (*Network, error) {
	if cfg.SSID == "" {
		return nil, errors.New("SSID が設定されていません")
	}
	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	devcfg := cyw43439.DefaultWifiConfig()
	devcfg.Logger = logger

	start := time.Now()
	dev := cyw43439.NewPicoWDevice()
	if err := dev.Init(devcfg); err != nil {
		return nil, errors.New("CYW43439 の初期化に失敗: " + err.Error())
	}
	if err := dev.Join(cfg.SSID, cyw43439.JoinOptions{Passphrase: cfg.Password}); err != nil {
		return nil, errors.New("無線 LAN への接続に失敗: " + err.Error())
	}
	mac, err := dev.HardwareAddr6()
	if err != nil {
		return nil, err
	}

	n := &Network{
		Dev:     dev,
		Stack:   new(xnet.StackAsync),
		sendbuf: make([]byte, cyw43439.MaxFrameSize),
	}
	err = n.Stack.Reset(xnet.StackConfig{
		Hostname:          cfg.Hostname,
		MaxActiveTCPPorts: cfg.MaxTCPPorts,
		RandSeed:          time.Since(start).Nanoseconds(),
		HardwareAddress:   mac,
		MTU:               1500,
	})
	if err != nil {
		return nil, err
	}
	dev.RecvEthHandle(n.Stack.IngressEthernet)
	go n.loop()

	// DHCP でアドレスを取得し、ルーターの MAC アドレスを調べる。
	var reqaddr [4]byte
	if cfg.RequestedAddr.Is4() {
		reqaddr = cfg.RequestedAddr.As4()
	}
	rstack := n.Stack.StackRetrying(backoff)
	results, err := rstack.DoDHCPv4(reqaddr, 3*time.Second, 3)
	if err != nil {
		return nil, errors.New("DHCP に失敗: " + err.Error())
	}
	if err := n.Stack.AssimilateDHCPResults(results); err != nil {
		return nil, err
	}
	n.Addr = results.AssignedAddr
	if results.Router.IsValid() {
		gw, err := rstack.DoResolveHardwareAddress6(results.Router, 500*time.Millisecond, 4)
		if err != nil {
			println("wifi: ルーターの MAC アドレスを取得できません:", err.Error())
		} else {
			n.Stack.SetGateway6(gw)
		}
	}
	return n, nil
}

// loop はパケットの受信と送信を繰り返します。
func (n *Network) loop() {
	for {
		recv, err := n.Dev.PollOne()
		if err != nil {
			println("wifi: 受信エラー:", err.Error())
		}
		send, err := n.Stack.EgressEthernet(n.sendbuf)
		if err != nil {
			println("wifi: 送信データ作成エラー:", err.Error())
		} else if send > 0 {
			if err := n.Dev.SendEth(n.sendbuf[:send]); err != nil {
				println("wifi: 送信エラー:", err.Error())
			}
		}
		if !recv && send == 0 {
			time.Sleep(pollSleep)
		}
	}
}

// backoff は DHCP などの再試行の待ち時間です (examples/cywnet と同じ)。
func backoff(consecutive uint) time.Duration {
	const (
		minWait = 100 * time.Microsecond
		maxWait = 20 * time.Millisecond
	)
	return min(minWait<<min(consecutive, 15), maxWait)
}
