# UF2 ファイルを Raspberry Pi Pico 2 W (RP2350) に書き込みます。
#
#   powershell.exe -NoProfile -ExecutionPolicy Bypass -File flash-uf2.ps1 -Uf2 <UF2ファイル>
#
# 1. RP2350 ドライブ (BOOTSEL モード) が見つかればそこへコピーします。
# 2. 見つからなければ、Pico の USB シリアル (COMx) を 1200bps で開閉して
#    BOOTSEL モードに再起動させます (TinyGo で書き込んだプログラムが動いている場合)。
# 3. それでも見つからなければ、BOOTSEL ボタンを押しながら USB を挿し直すよう案内して待ちます。
param(
    [Parameter(Mandatory = $true)][string]$Uf2,
    [int]$TimeoutSec = 60
)
$ErrorActionPreference = 'Stop'
# WSL から呼んだときに日本語が文字化けしないよう UTF-8 で出力する
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

function Find-PicoDrive {
    Get-CimInstance Win32_LogicalDisk |
        Where-Object { $_.VolumeName -eq 'RP2350' } |
        Select-Object -First 1 -ExpandProperty DeviceID
}

function Wait-PicoDrive([int]$Seconds) {
    $deadline = (Get-Date).AddSeconds($Seconds)
    while ((Get-Date) -lt $deadline) {
        $d = Find-PicoDrive
        if ($d) { return $d }
        Start-Sleep -Milliseconds 500
    }
    return $null
}

if (-not (Test-Path -LiteralPath $Uf2)) {
    throw "UF2 ファイルが見つかりません: $Uf2"
}

$drive = Find-PicoDrive
if (-not $drive) {
    # Raspberry Pi の USB VID は 2E8A
    $ports = Get-CimInstance Win32_PnPEntity |
        Where-Object { $_.PNPDeviceID -like 'USB\VID_2E8A*' -and $_.Name -match '\((COM\d+)\)' }
    foreach ($p in $ports) {
        if ($p.Name -match '\((COM\d+)\)') {
            $com = $Matches[1]
            Write-Host "$com を 1200bps で開閉して BOOTSEL モードにします"
            try {
                $sp = New-Object System.IO.Ports.SerialPort $com, 1200
                $sp.DtrEnable = $false
                $sp.Open()
                Start-Sleep -Milliseconds 200
                $sp.Close()
            } catch {
                # 再起動でポートが消えると例外になることがあるので無視する
            }
        }
    }
    $drive = Wait-PicoDrive 10
}
if (-not $drive) {
    Write-Host "RP2350 ドライブが見つかりません。BOOTSEL ボタンを押しながら USB ケーブルを挿し直してください..."
    $drive = Wait-PicoDrive $TimeoutSec
}
if (-not $drive) {
    throw "RP2350 ドライブが見つかりませんでした"
}

Write-Host "$Uf2 を $drive に書き込みます"
Copy-Item -LiteralPath $Uf2 -Destination "$drive\"
Write-Host "書き込み完了 (Pico は自動的に再起動します)"
