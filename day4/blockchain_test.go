package main

import (
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	cli := &CLI{}
	cli.getBalance("zach")
	cli.send("zach", "zhou", 1)
}
