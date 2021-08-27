package protocol

import "io"

func Quit(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
	return io.EOF
}
