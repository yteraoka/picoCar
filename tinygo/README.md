# picoCar TinyGo 版 (Raspberry Pi Pico 2 W)

リポジトリ直下の MicroPython プログラムを [TinyGo](https://tinygo.org/) に移植し、
Raspberry Pi Pico 2 W (RP2350) で動くようにしたものです。

動作確認環境: TinyGo 0.42.0 / Go 1.27 (全プログラムのビルドのみ確認。実機での走行は未確認)

## ファイル構成

| TinyGo 版 | 元の Python | 内容 |
|---|---|---|
| `config/config.go` | `config.py` | サーボのデューティ値・回転時間の設定 |
| `lib/servo` | `lib_servo.py` | 左右の車輪 (GPIO0, GPIO1)、回転、`Faraway` |
| `lib/sr04` | `lib_SR04.py` | 超音波距離センサー (TRIG GPIO14, ECHO GPIO15) |
| `lib/ir` | `lib_iR.py` | 赤外線センサー (GPIO16) |
| `lib/cds` | `lib_Cds.py` | 明るさセンサー (GPIO26 / ADC0) |
| `lib/mode` | `lib_mode.py` | mode ピン (GPIO21, 20, 19, 18) |
| `lib/led` | `lib_LED.py` | GPIO17 の LED |
| `lib/picoled` | `lib_LED_pico.py` | 基板上の LED |
| `lib/bootsw` | `lib_bootSW.py` | BOOTSEL スイッチ待ち |
| `cmd/main` | `main.py`, `main_01.py` | 障害物を避けながら走る |
| `cmd/main02` | `main_02.py` | main + 暗くなったら停止 |
| `cmd/adjust_run` | `adjust_run.py` | 前進の直進性の調整 |
| `cmd/adjust_rotate_left` | `adjust_rotate_left.py` | 左回転時間の調整 |
| `cmd/adjust_rotate_right` | `adjust_rotate_rgith.py` | 右回転時間の調整 |
| `cmd/adjust_stop` | `adjust_stop.py` | サーボの停止値を求める |
| `cmd/faraway` | `lib_servo.py` の `main()` | 一番遠い方向を向く |
| `cmd/test_*` | 各 `lib_*.py` の `main()` | センサー・LED の動作確認 |

## Python 版との違い

- **書き込むのは 1 つのプログラムだけ**です。MicroPython のように `main.py` と `config.py`
  を別々にアップロードするのではなく、`config/config.go` を編集して再ビルドし、
  できた `.uf2` ファイルを書き込みます。調整用プログラム (`adjust_*`) で調整した後は、
  `main` を書き込み直してください。
- **基板上の LED**: Pico 2 W の LED は GPIO25 ではなく無線チップ (CYW43439) につながっているため、
  [soypat/cyw43439](https://github.com/soypat/cyw43439) ドライバで制御しています。
  無線チップのファームウェアを含むので、LED を使うプログラムは UF2 が約 1.2MB になり、
  起動時に初期化で 1〜2 秒かかります。
- **BOOTSEL 待ちの LED**: Python 版は LED を PWM でふわっと明滅させていましたが、
  無線チップ経由の LED は PWM できないので、同じくらいの周期の点滅にしています。
- **BOOTSEL スイッチの読み取り**: TinyGo には `rp2.bootsel_button()` 相当の関数がないため、
  pico-sdk と同じ方法 (割り込みを止めて RAM 上の関数でフラッシュの CS 端子を読む) で実装しています
  (`lib/bootsw/bootsel_rp2350.go`)。
- デューティ値は Python 版と同じ `duty_u16()` のスケール (0〜65535) なので、
  Python 版で調整した値をそのまま `config/config.go` に書けます。
- `lib_SR04.py` は測定失敗時に GPIO17 の LED が点いたままになっていましたが、TinyGo 版では消灯します。

## ビルド

TinyGo をインストールします (WSL の Ubuntu の場合)。

```sh
# https://github.com/tinygo-org/tinygo/releases から最新版の .deb を入れる
wget https://github.com/tinygo-org/tinygo/releases/download/v0.42.0/tinygo_0.42.0_amd64.deb
sudo dpkg -i tinygo_0.42.0_amd64.deb
tinygo version
```

TinyGo とは別に Go 本体 (1.24 以降) も必要です。

```sh
cd tinygo
make                 # 全プログラムをビルドして build/*.uf2 を作る
make build PROG=main # main だけビルド → build/main.uf2
```

## WSL から Pico 2 W に書き込む方法

WSL2 からは USB 機器が直接見えないため、次のどちらかの方法を使います。

### 方法 A: Windows 側のドライブにコピーする (おすすめ・追加ソフト不要)

Pico 2 W を BOOTSEL モードにすると、Windows に `RP2350` という名前の USB ドライブとして現れます。
そこへ `.uf2` ファイルをコピーすれば書き込みが完了し、自動的に再起動します。

```sh
make flash-wsl PROG=main
```

`scripts/flash-wsl.sh` が WSL から `powershell.exe` を呼び出し、`scripts/flash-uf2.ps1` が次の順に処理します。

1. `RP2350` ドライブがあれば、そこへコピーする
2. なければ、Pico の USB シリアル (COMx) を 1200bps で開閉して BOOTSEL モードに再起動させる
   (TinyGo で書き込んだプログラムが動いているときだけ有効)
3. それでも見つからなければ、BOOTSEL ボタンを押しながら USB を挿し直すよう表示して待つ

**初めて書き込むとき** (MicroPython が入っている、または新品のとき) は、
BOOTSEL ボタンを押しながら USB ケーブルを挿してから実行してください。
2 回目以降は、ボタンを押さなくても 2. の方法で自動的に書き込めます
(ただし Tera Term などでその COM ポートを開いている間はできないので、閉じてから実行してください)。

手動でやる場合は次のとおりです。

```sh
# RP2350 ドライブが E: の場合
powershell.exe -c "Copy-Item '$(wslpath -w build/main.uf2)' E:\\"
# または、エクスプローラーで \\wsl.localhost\Ubuntu\... の build/main.uf2 を RP2350 ドライブへドラッグ
```

`/mnt/e` のような WSL のドライブマウントは WSL 起動後に挿したリムーバブルドライブには自動で作られないので、
`cp build/main.uf2 /mnt/e/` を使う場合は先に `sudo mkdir -p /mnt/e && sudo mount -t drvfs E: /mnt/e` が必要です。

**シリアル出力 (`println` の表示) を見るには**、Windows 側で Tera Term、PuTTY、
Arduino IDE のシリアルモニタなどで Pico の COM ポートを開きます (ボーレートは任意)。
WSL2 からは Windows の COM ポートを直接開けません。

### 方法 B: usbipd-win で USB を WSL に接続する

[usbipd-win](https://github.com/dorssel/usbipd-win) を使うと、USB 機器を WSL に直接接続でき、
`tinygo monitor` などを WSL 内で使えます。

1. Windows にインストール: PowerShell で `winget install usbipd`
2. Pico を挿して、Windows の PowerShell で確認・共有設定 (bind は管理者権限で 1 回だけ)

   ```powershell
   usbipd list                       # 2e8a:xxxx の BUSID を確認 (例: 2-3)
   usbipd bind --busid 2-3           # 管理者の PowerShell で実行
   usbipd attach --wsl --busid 2-3 --auto-attach
   ```

   Pico は書き込みのときに「プログラム実行中 (2e8a:000a, シリアル)」と
   「BOOTSEL モード (2e8a:000f)」で別の USB 機器として再接続されるので、
   `--auto-attach` を付けたままにしておきます。BOOTSEL モードの状態でも `usbipd list` で
   `Not shared` になっていたら、その状態でも `usbipd bind` してください。

3. WSL 側でシリアルを使う準備

   ```sh
   sudo modprobe cdc_acm              # /dev/ttyACM0 が出ない場合
   sudo usermod -aG dialout $USER     # 実行後に WSL を再起動
   make monitor                       # = tinygo monitor -target=pico2-w
   ```

4. 書き込みは [picotool](https://github.com/raspberrypi/picotool) を使うのが確実です
   (`tinygo flash` は BOOTSEL ドライブが自動マウントされることを前提にしているため、
   自動マウントのない WSL ではそのままでは使えません)。

   ```sh
   stty -F /dev/ttyACM0 1200          # TinyGo のプログラムを BOOTSEL モードに再起動させる
   sleep 2
   sudo picotool load -x build/main.uf2   # 書き込んで実行
   ```

   RP2350 には picotool 2.0 以降が必要です (Ubuntu の apt パッケージは古く RP2350 に対応していない場合があります)。
   ビルド済みのものが [raspberrypi/pico-sdk-tools のリリース](https://github.com/raspberrypi/pico-sdk-tools/releases)
   に `picotool-<版>-x86_64-lin.tar.gz` として公開されています。

手軽さでは方法 A、シリアルモニタまで WSL で完結させたい場合は方法 B がおすすめです。

## 調整の流れ

1. `make flash-wsl PROG=adjust_stop` — `cmd/adjust_stop/main.go` の `wheel` を `servo.Right` / `servo.Left`
   に切り替えて、それぞれのサーボの停止値を求め、`config/config.go` の `Stop` に設定
2. `make flash-wsl PROG=adjust_run` — 前進の直進性を見て `Forward` を調整 (説明資料/前進時の直線性の調整.txt)
3. `make flash-wsl PROG=adjust_rotate_left` / `adjust_rotate_right` — 180 度回転する時間を
   `LeftRotateTime` / `RightRotateTime` に設定 (説明資料/左右回転の調整.txt)
4. `make flash-wsl PROG=main` — 本番のプログラムを書き込む
