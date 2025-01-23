package main_test

import (
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd"
)

var execute = cmd.Execute // Injected dependency for testing

func TestMain(m *testing.M) {
	execute()
}
