package kecpsignal_test

import (
	"log"
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-signal"
)

func TestSetLogger(t *testing.T) {
	SetLogger(log.Default())
}
