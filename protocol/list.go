package protocol

import (
	"fmt"
	"io"
	"strings"
)

func List(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
	fmt.Fprintln(ofp, strings.Join(conn.Srv.Config.ListPlugin(conn.Multigraph), " "))
	return nil
}
