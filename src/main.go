package main

import (
    "bufio"
    "context"
    "encoding/binary"
    "flag"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"
)

var (
    flagInterval = flag.Duration("interval", 5*time.Second, "Sampling interval (e.g. 5s, 1s, 250ms)")
    flagLogFile  = flag.String("logfile", "/tmp/assignment_sensor.log", "Log file path (default in /tmp with documented fallback)")
    flagDevice   = flag.String("device", "/dev/urandom", "Entropy source (use /dev/urandom, /dev/random, or 'internal')")
)

func main() {
    flag.Parse()
    if *flagInterval <= 0 {
        fmt.Fprintln(os.Stderr, "interval must be > 0")
        os.Exit(1)
    }

    logFile, usedFallback, err := openWithFallback(*flagLogFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open log file: %v
", err)
        os.Exit(1)
    }
    defer logFile.Close()
    writer := bufio.NewWriterSize(logFile, 64*1024)
    defer writer.Flush()

    now := time.Now().UTC().Format(time.RFC3339Nano)
    first := fmt.Sprintf("%s | START interval=%s device=%s logfile=%s", now, flagInterval.String(), *flagDevice, logFile.Name())
    if usedFallback {
        first += " (fallback_log_path_used)"
    }
    fmt.Fprintln(writer, first)
    writer.Flush()

    rdr, closer, err := openSource(*flagDevice)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open device/source: %v
", err)
        os.Exit(1)
    }
    if closer != nil {
        defer closer.Close()
    }

    ctx, cancel := signalContext()
    defer cancel()
    ticker := time.NewTicker(*flagInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            ts := time.Now().UTC().Format(time.RFC3339Nano)
            fmt.Fprintf(writer, "%s | STOP
", ts)
            writer.Flush()
            return
        case <-ticker.C:
            val, rerr := sampleUint64(rdr)
            if rerr != nil {
                ts := time.Now().UTC().Format(time.RFC3339Nano)
                fmt.Fprintf(writer, "%s | ERROR read: %v
", ts, rerr)
                writer.Flush()
                os.Exit(2)
            }
            ts := time.Now().UTC().Format(time.RFC3339Nano)
            fmt.Fprintf(writer, "%s | 0x%016x (%d)
", ts, val, val)
            writer.Flush()
        }
    }
}

func signalContext() (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(context.Background())
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
    go func() {
        <-ch
        cancel()
    }()
    return ctx, cancel
}

func openWithFallback(path string) (*os.File, bool, error) {
    f, err := openForAppend(path)
    if err == nil {
        return f, false, nil
    }
    base := filepath.Base(path)
    fallback := filepath.Join("/var/tmp", base)
    ff, ferr := openForAppend(fallback)
    if ferr == nil {
        return ff, true, nil
    }
    return nil, false, fmt.Errorf("open %s: %w; fallback %s also failed: %v", path, err, fallback, ferr)
}

func openForAppend(path string) (*os.File, error) {
    _ = os.MkdirAll(filepath.Dir(path), 0o755)
    return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
}

func openSource(device string) (io.Reader, io.Closer, error) {
    if device == "internal" {
        src := rand.New(rand.NewSource(time.Now().UnixNano()))
        return src, nil, nil
    }
    f, err := os.Open(device)
    if err != nil {
        return nil, nil, err
    }
    return f, f, nil
}

func sampleUint64(r io.Reader) (uint64, error) {
    var b [8]byte
    _, err := io.ReadFull(r, b[:])
    if err != nil {
        return 0, err
    }
    return binary.LittleEndian.Uint64(b[:]), nil
}
