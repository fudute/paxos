package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig()
	if err != nil {
		t.Error(err)
	}
	_ = c
}
