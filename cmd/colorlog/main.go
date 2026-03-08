package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lfsgroup/colorlog"
)

func main() {
	handler := colorlog.NewHandler(os.Stdout, nil)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var raw map[string]any
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			fmt.Println(line)
			continue
		}

		// Extract timestamp
		t := time.Now()
		if ts, ok := raw["time"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339Nano, ts); err == nil {
				t = parsed
			}
			delete(raw, "time")
		}

		// Extract level
		level := slog.LevelInfo
		if lvl, ok := raw["level"].(string); ok {
			_ = level.UnmarshalText([]byte(lvl))
			delete(raw, "level")
		}

		// Extract message
		msg := ""
		if m, ok := raw["msg"].(string); ok {
			msg = m
			delete(raw, "msg")
		} else if m, ok := raw["message"].(string); ok {
			msg = m
			delete(raw, "message")
		}

		r := slog.NewRecord(t, level, msg, 0)
		for k, v := range raw {
			r.AddAttrs(slog.Any(k, v))
		}

		_ = handler.Handle(context.Background(), r)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin:", err)
		os.Exit(1)
	}
}
