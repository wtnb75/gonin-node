package protocol

import (
	"fmt"
	"io"
)

func Nodes(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
	fmt.Fprintln(ofp, conn.Srv.Name)
	fmt.Fprintln(ofp, ".")
	return nil
}
