package config

import (
	"testing"
	"time"
)

func fakeEnv(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoad_NothingSet_UsesDefaults(t *testing.T) {
	cfg, err := Load(fakeEnv(map[string]string{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 10s", cfg.ShutdownTimeout)
	}
}

func TestLoad_ValuesSet_UsesThem(t *testing.T) {
	cfg, err := Load(fakeEnv(map[string]string{
		"TASKBOARD_HTTP_ADDR":                "127.0.0.1:9000",
		"TASKBOARD_SHUTDOWN_TIMEOUT_SECONDS": "30",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPAddr != "127.0.0.1:9000" {
		t.Errorf("HTTPAddr = %q", cfg.HTTPAddr)
	}
	if cfg.ShutdownTimeout != 30*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 30s", cfg.ShutdownTimeout)
	}
}

func TestLoad_BadShutdownTimeout_ReturnsError(t *testing.T) {
	cases := map[string]string{
		"not a number": "abc",
		"zero":         "0",
		"negative":     "-5",
		"too big":      "301",
		"overflow":     "9999999999999",
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Load(fakeEnv(map[string]string{
				"TASKBOARD_SHUTDOWN_TIMEOUT_SECONDS": raw,
			}))
			if err == nil {
				t.Fatalf("expected error for %q, got nil", raw)
			}
		})
	}
}
