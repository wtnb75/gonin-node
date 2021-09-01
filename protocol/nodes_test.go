package protocol

import (
	"bytes"
	"strings"
	"testing"
)

func TestNodes(t *testing.T) {
	conn := Conn{Srv: &Server{Name: "hello"}}
	buf := bytes.Buffer{}
	if err := Nodes(&conn, "nodes", []string{}, &buf); err != nil {
		t.Error("error(nodes)", err)
	}
	if !strings.HasPrefix(buf.String(), "hello") {
		t.Error("nodes(name)")
	}
}
