package main

import "testing"

func TestGetInputDir(t *testing.T) {
	var d string
	d = getInputDir("/home/gustav/test/dsd.yml")
	if d != "/home/gustav/test" {
		t.Error("Expected /home/gustav/test, got ", d)
	}
	d = getInputDir("test/dsd.yml")
	if d != "/home/gustav/code/go/src/github.com/gntech/gmb/test" {
		t.Error("Expected /home/gustav/test, got ", d)
	}
}
