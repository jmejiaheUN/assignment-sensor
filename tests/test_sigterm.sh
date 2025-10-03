#!/usr/bin/env bash
set -euo pipefail
BIN="${1:-./build/assignment-sensor}"
echo "[i] Running SIGTERM test..."
$BIN --interval=1s --device=internal &
PID=$!
sleep 2
kill -TERM $PID
wait $PID || true
LOG=/tmp/assignment_sensor.log
ALT=/var/tmp/assignment_sensor.log
FILE=""
if [[ -f "$LOG" ]]; then FILE="$LOG"; elif [[ -f "$ALT" ]]; then FILE="$ALT"; fi
if [[ -z "$FILE" ]]; then echo "[x] No log file found"; exit 1; fi
tail -n 2 "$FILE"
if tail -n 1 "$FILE" | grep -q "STOP"; then
  echo "[✓] Graceful STOP recorded"
else
  echo "[x] STOP line missing"; exit 1
fi
