package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestServeCommand_RequiresInternalToken(t *testing.T) {
	t.Setenv("RESCALER_INTERNAL_TOKEN", "")
	t.Setenv("RESCALER_HTTP_ADDR", "")
	t.Setenv("RESCALER_DB_PATH", "")
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"serve"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error when RESCALER_INTERNAL_TOKEN is unset, got nil (output=%q)", buf.String())
	}
	if !strings.Contains(err.Error(), "RESCALER_INTERNAL_TOKEN") {
		t.Fatalf("expected error to mention RESCALER_INTERNAL_TOKEN, got %q", err.Error())
	}
}

func TestServeCommand_HelpText(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"serve", "--help"})
	_ = rootCmd.Execute()
	if !strings.Contains(buf.String(), "RESCALER_INTERNAL_TOKEN") {
		t.Fatalf("expected help text to mention RESCALER_INTERNAL_TOKEN, got %q", buf.String())
	}
}