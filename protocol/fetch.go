package protocol

import (
	"fmt"
	"io"
	"log"
)

func Fetch(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
	if len(arg) == 0 {
		fmt.Fprintln(ofp, "# Unknown service")
		fmt.Fprintln(ofp, ".")
		return nil
	}
	for _, v := range arg {
		plugin, err := conn.Srv.Config.GetPlugin(v)
		if err != nil {
			log.Println("plugin", err)
			fmt.Fprintln(ofp, "# Unknown service")
			fmt.Fprintln(ofp, ".")
		} else {
			err := plugin.Execute(ofp)
			if err != nil {
				log.Println("exec error", err)
			}
			fmt.Fprintln(ofp, ".")
		}
	}
	return nil
}
