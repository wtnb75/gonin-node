package protocol

import (
	"bytes"
	"strings"
	"testing"
)

func TestCap(t *testing.T) {
	conn := Conn{}
	out := bytes.Buffer{}
	if conn.Multigraph {
		t.Error("multigraph(init)")
	}
	if err := Cap(&conn, "cap", []string{"multigraph"}, &out); err != nil {
		t.Error("cap error", err)
	}
	if !conn.Multigraph {
		t.Error("multigraph")
	}
	if !strings.Contains(out.String(), "cap") {
		t.Error("output(cap)")
	}
	if !strings.Contains(out.String(), "multigraph") {
		t.Error("output(multigraph)")
	}
}
