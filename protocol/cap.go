package protocol

import (
	"fmt"
	"io"
	"log"
)

func Cap(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
	fmt.Fprintln(ofp, "cap multigraph dirtyconfig")
	for _, v := range arg {
		if v == "multigraph" {
			// enable multigraph mode
			log.Println("enable multigraph mode")
			conn.Multigraph = true
		} else if v == "dirtyconfig" {
			log.Println("enable dirtyconfig mode")
			conn.DirtyConfig = true
		}
	}
	return nil
}
