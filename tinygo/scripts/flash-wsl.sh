#!/usr/bin/env bash
# WSL から Windows の PowerShell を使って UF2 を Pico 2 W に書き込みます。
# usbipd-win は不要です。
#
#   scripts/flash-wsl.sh build/main.uf2
set -euo pipefail

if [ $# -ne 1 ]; then
	echo "usage: $0 <file.uf2>" >&2
	exit 1
fi

script_dir=$(cd "$(dirname "$0")" && pwd)
uf2=$(wslpath -w "$(realpath "$1")")
ps1=$(wslpath -w "$script_dir/flash-uf2.ps1")

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$ps1" -Uf2 "$uf2"
