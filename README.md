# assignment-sensor (Go)

Logs an ISO-8601 timestamp and a mock sensor value at a configurable interval.
Default logfile is `/tmp/assignment_sensor.log` with fallback to `/var/tmp/assignment_sensor.log` if `/tmp` is not writable.
Handles `SIGTERM` gracefully (writes a final `STOP` line) and includes a `systemd` unit to start at `multi-user.target`.

## Why `/dev/urandom`?
We use `/dev/urandom` as the default mock source because it's universally available on Linux and safe to read without blocking.
For portability you can also choose `--device=internal` to use a built-in PRNG.

## Cloning the repo
git clone https://github.com/<your-username>/assignment-sensor.git
cd assignment-sensor


## Build & run (foreground)
```bash
make build
./build/assignment-sensor --interval=2s --device=internal
tail -f /tmp/assignment_sensor.log
```

## Service setup (systemd)
> The service runs as a restricted user `assignment` and logs to `/var/log/assignment-sensor/assignment_sensor.log`.
> The program default remains `/tmp` when run manually, satisfying the assignment requirement with documented fallback.

1) Build
```bash
make build
```

2) Prepare service user & log dir (one time)
```bash
./scripts/setup_service_user.sh
```

3) Install + enable
```bash
sudo make install
sudo systemctl enable --now assignment-sensor.service
```

4) Verify & tail logs
```bash
systemctl status assignment-sensor.service
journalctl -u assignment-sensor.service -n 20 --no-pager
tail -n 5 /var/log/assignment-sensor/assignment_sensor.log
```

5) Stop (should write STOP)
```bash
sudo systemctl stop assignment-sensor.service
tail -n 3 /var/log/assignment-sensor/assignment_sensor.log
```

6) Uninstall
```bash
sudo make uninstall
```

## Tests
```bash
chmod +x tests/*.sh
./tests/test_sigterm.sh
./tests/test_fallback.sh
```
## Failure behavior test

./build/assignment-sensor --interval=0s --device=internal
echo "exit code: $?"   # expect non-zero

## ✅ Verification Checklist
- Build: `make build`
- Manual run logs to `/tmp/...` or `/var/tmp/...` (fallback) with ISO-8601 timestamps
- Service starts via systemd at/above `multi-user.target`, logs under `/var/log/assignment-sensor/`
- `systemctl stop` writes a final `STOP`
- Uninstall removes binary and unit

## AI Usage
See `ai/` for prompt log, reflection, and provenance.
