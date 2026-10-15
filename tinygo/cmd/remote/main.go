// remote はスマホのブラウザから無線 LAN で操作するリモコン版です。
//
// 起動すると無線 LAN に接続し、HTTP サーバーを動かします。スマホなどで
// http://<picoCar の IP アドレス>/ を開くと操作画面が表示され、ボタンを押している間だけ
// 前進・後退・左右回転します。IP アドレスはシリアル出力 (make monitor) に表示されます。
//
// 無線 LAN の SSID とパスワードはビルド時に指定します (README.md 参照)。
//
//	make build PROG=remote WIFI_SSID=xxxx WIFI_PASSWORD=yyyy
//
// 安全のため、画面からの指示が timeout の間届かなければ停止します。
// 接続できると基板上の LED が点灯します。
package main

import (
	"net/netip"
	"strconv"
	"sync"
	"time"

	_ "embed"

	"github.com/momorara/picoCar/tinygo/lib/picoled"
	"github.com/momorara/picoCar/tinygo/lib/servo"
	"github.com/momorara/picoCar/tinygo/lib/wifi"
	"github.com/soypat/lneto"
	"github.com/soypat/lneto/http/httpraw"
	"github.com/soypat/lneto/tcp"
	"github.com/soypat/lneto/x/xnet"
)

// ビルド時に -ldflags "-X main.ssid=... -X main.password=..." で設定します。
var (
	ssid     string
	password string
	hostname = "picocar"
	ipAddr   string // DHCP で希望する IP アドレス (空ならルーター任せ)
)

const (
	listenPort = 80
	maxConns   = 4                      // 同時に受け付ける接続数
	httpBuf    = 512                    // HTTP リクエストヘッダの最大サイズ
	idleTime   = 10 * time.Second       // この間リクエストがなければ接続を閉じる
	timeout    = 800 * time.Millisecond // この間指示がなければ停止 (画面は 250ms ごとに送る)
	tick       = 20 * time.Millisecond  // 走行状態を確認する間隔
)

//go:embed index.html
var webPage []byte

func main() {
	servo.Stop()
	time.Sleep(2 * time.Second) // シリアル出力を見られるように少し待つ

	println("無線 LAN に接続中:", ssid)
	var reqAddr netip.Addr
	if ipAddr != "" {
		reqAddr, _ = netip.ParseAddr(ipAddr)
	}
	network, err := wifi.Connect(wifi.Config{
		SSID:          ssid,
		Password:      password,
		Hostname:      hostname,
		RequestedAddr: reqAddr,
		MaxTCPPorts:   1,
	})
	if err != nil {
		// 接続できなければ停止したまま、エラーを表示し続ける
		for {
			println("remote:", err.Error())
			time.Sleep(5 * time.Second)
		}
	}
	picoled.UseDevice(network.Dev)
	picoled.On()

	pool, err := network.Stack.NewTCPPool(xnet.TCPPoolConfig{
		PoolSize:           maxConns,
		QueueSize:          maxConns,
		TxBufSize:          len(webPage) + 256,
		RxBufSize:          httpBuf,
		EstablishedTimeout: 5 * time.Second,
		ClosingTimeout:     5 * time.Second,
		NewUserData: func() any {
			hdr := new(httpraw.HeaderV1)
			hdr.Reset(make([]byte, 0, httpBuf), 16)
			return hdr
		},
		NewBackoff: func() lneto.BackoffStrategy { return connBackoff },
	})
	if err != nil {
		panic("remote: TCP プール作成に失敗: " + err.Error())
	}
	var listener tcp.Listener
	if err := listener.Reset(listenPort, pool); err != nil {
		panic("remote: listener: " + err.Error())
	}
	if err := network.Stack.RegisterListenerTCP(&listener); err != nil {
		panic("remote: listener 登録: " + err.Error())
	}
	println("操作画面: http://" + network.Addr.String() + "/")

	go drive()

	for {
		if listener.NumberOfReadyToAccept() == 0 {
			time.Sleep(5 * time.Millisecond)
			pool.CheckTimeouts()
			continue
		}
		conn, err := listener.TryAccept()
		if err != nil {
			println("remote: accept:", err.Error())
			continue
		}
		println("接続:", conn.RemoteAddrPort().String())
		go handleConn(conn, conn.UserData().(*httpraw.HeaderV1))
	}
}

// 走行指示。画面から届いた最新の指示と、その時刻を保持します。
var (
	mu      sync.Mutex
	command byte = 's'
	lastCmd time.Time
)

func setCommand(c byte) {
	mu.Lock()
	command = c
	lastCmd = time.Now()
	mu.Unlock()
}

// drive は指示に合わせてサーボを動かします。指示が timeout の間届かなければ停止します。
func drive() {
	current := byte('s')
	for {
		mu.Lock()
		c := command
		if c != 's' && time.Since(lastCmd) > timeout {
			println("指示が途切れたので停止")
			c = 's'
			command = 's'
		}
		mu.Unlock()

		if c != current {
			switch c {
			case 'f':
				servo.Run()
			case 'b':
				servo.Back()
			case 'l':
				servo.SpinLeft()
			case 'r':
				servo.SpinRight()
			default:
				servo.Stop()
			}
			current = c
		}
		time.Sleep(tick)
	}
}

// handleConn は 1 つの接続で届く HTTP リクエストに順に応答します (keep-alive)。
func handleConn(conn tcp.ConnHandle, hdr *httpraw.HeaderV1) {
	defer conn.Close()
	var buf [256]byte
	for {
		hdr.Reset(nil, 0)
		conn.SetReadDeadline(time.Now().Add(idleTime))
		if !readRequest(conn, hdr, buf[:]) {
			return
		}
		status, contentType, body := route(string(hdr.RequestPath()))

		hdr.Reset(nil, 0)
		hdr.SetProtocol("HTTP/1.1")
		hdr.SetStatus(status, statusText(status))
		hdr.Set("Content-Type", contentType)
		hdr.Set("Content-Length", strconv.Itoa(len(body)))
		hdr.Set("Cache-Control", "no-store")
		resp, err := hdr.AppendResponse(buf[:0])
		if err != nil {
			println("remote: レスポンス作成:", err.Error())
			return
		}
		if _, err := conn.Write(resp); err != nil {
			return
		}
		if _, err := conn.Write(body); err != nil {
			return
		}
	}
}

// readRequest はリクエストヘッダを最後まで読みます。接続が閉じられたら false を返します。
func readRequest(conn tcp.ConnHandle, hdr *httpraw.HeaderV1, buf []byte) bool {
	const asResponse = false
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			hdr.ReadFromBytes(buf[:n])
			needMore, perr := hdr.TryParse(asResponse)
			if perr != nil && !needMore {
				println("remote: HTTP の解析に失敗:", perr.Error())
				return false
			}
			if !needMore {
				return true
			}
			if hdr.BufferReceived() >= httpBuf {
				println("remote: リクエストが大きすぎます")
				return false
			}
		}
		if err != nil {
			return false // 相手が閉じた、または idleTime 経過
		}
	}
}

// route は URI に対応する応答を返します。/f /b /l /r /s は走行指示です。
func route(uri string) (status, contentType string, body []byte) {
	switch uri {
	case "/":
		return "200", "text/html; charset=utf-8", webPage
	case "/f", "/b", "/l", "/r", "/s":
		setCommand(uri[1])
		return "200", "text/plain", []byte("ok")
	}
	return "404", "text/plain", []byte("not found")
}

// connBackoff は接続の読み書きで待つときの間隔です。
func connBackoff(consecutive uint) time.Duration {
	return min(time.Millisecond<<min(consecutive, 4), 10*time.Millisecond)
}

func statusText(status string) string {
	if status == "200" {
		return "OK"
	}
	return "Not Found"
}
