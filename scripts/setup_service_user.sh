#!/usr/bin/env bash
set -euo pipefail

if ! id assignment &>/dev/null; then
  sudo useradd --system --no-create-home --shell /usr/sbin/nologin assignment
fi

sudo install -d -o assignment -g assignment -m 0755 /var/log/assignment-sensor

echo "[✓] Service user 'assignment' and /var/log/assignment-sensor prepared."
