package main

import (
	"testing"
)

func TestSayHello(t *testing.T) {
	expect := "hello world"
	output := SayHello()

	if expect != output {
		t.Errorf("expected %v not match the %v", expect, output)
	}
}
