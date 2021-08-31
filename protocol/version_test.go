package protocol

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	conn := Conn{Srv: &Server{Name: "hello"}}
	out := bytes.Buffer{}
	if err := Version(&conn, "version", []string{}, &out); err != nil {
		t.Error("cap error", err)
	}
	if !strings.Contains(out.String(), "node on") {
		t.Error("output(node string)")
	}
	if !strings.Contains(out.String(), "version") {
		t.Error("output(version string)")
	}
}
