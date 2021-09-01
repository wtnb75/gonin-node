package protocol

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRepl(t *testing.T) {
	conn := Conn{Srv: &Server{Name: "hello"}}
	ifp := bytes.NewBufferString("version\n.\n")
	ofp := bytes.Buffer{}
	if err := conn.Repl(ifp, &ofp); err != nil {
		if err != io.EOF {
			t.Error("repl error", err)
		}
	}
	if !strings.Contains(ofp.String(), "gonin node") {
		t.Error("not contains(node)")
	}
	if !strings.Contains(ofp.String(), "version") {
		t.Error("not contains(version)")
	}
}
