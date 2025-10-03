#!/usr/bin/env bash
set -euo pipefail
BIN="${1:-./build/assignment-sensor}"
echo "[i] Running fallback test..."
$BIN --interval=1s --logfile=/root/assignment_sensor.log --device=internal &
PID=$!
sleep 2
kill -TERM $PID || true
sleep 1
if grep -q "fallback_log_path_used" /var/tmp/assignment_sensor.log; then
  echo "[✓] Fallback worked -> /var/tmp/assignment_sensor.log"
  tail -n 3 /var/tmp/assignment_sensor.log
else
  echo "[x] Fallback evidence not found."; exit 1
fi
