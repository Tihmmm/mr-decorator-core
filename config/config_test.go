package config

import "testing"

func TestNewConfig(t *testing.T) {
	config, err := NewConfig("../config.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config)
}
