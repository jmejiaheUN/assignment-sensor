# AI Reflection (<=500 words)
Goal: Go daemon as a systemd service that samples a mock sensor, logs ISO-8601 timestamp + value, defaults to /tmp with fallback, and shuts down gracefully on SIGTERM.

Process: Used AI to draft structure and code. Fixed Makefile tab/target issues, removed rebuild during install, corrected io.NopCloser misuse by returning a nil closer for the internal PRNG, and addressed sudo PATH problems by separating build from install.

Reliability & permissions: Running the service as my user created permission errors when appending to /tmp logs. I added a dedicated system user and a service-owned log directory under /var/log. The foreground default remains /tmp (with fallback), meeting the assignment requirement; the service uses /var/log for safety. README documents both behaviors.

Validation: Manual runs produced ISO-8601 lines and a STOP line on Ctrl+C. Test scripts verify SIGTERM behavior and the /var/tmp fallback. Under systemd, status shows the service running and logs update continuously; stopping the service appends STOP. Fatal init errors (bad device) return non-zero.

Outcome: The repo builds with `make`, installs cleanly, runs at multi-user.target, logs correctly, includes tests and AI provenance.
