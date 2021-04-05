package config

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {

	testconfig := Config{
		App:    "app-name",
		Module: "script",
		URL:    "http://apps.rlxos.dev/1507/",

		Script: []string{
			"ls -al",
		},
	}

	conf, err := Load("sample.yml")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testconfig, *conf) {
		t.Fatalf("%v != %v", testconfig, *conf)
	}
}
