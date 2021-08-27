package protocol

import (
	"fmt"
	"io"
)

func Version(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
	fmt.Fprintln(ofp, "gonin node on", conn.Srv.Name, "version:", VersionString)
	return nil
}
