package protocol

import (
	"bytes"
	"io"
	"testing"
)

func TestQuit(t *testing.T) {
	conn := Conn{Srv: &Server{Name: "hello"}}
	out := bytes.Buffer{}
	if err := Quit(&conn, "quit", []string{}, &out); err == nil {
		t.Error("no EOF")
	} else if err != io.EOF {
		t.Error("not EOF", err)
	}
}
