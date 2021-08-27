package protocol

import (
	"fmt"
	"io"
	"log"
)

func Config(conn *Conn, cmd string, arg []string, ofp io.Writer) error {
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
			err := plugin.ShowConfig(conn.DirtyConfig, ofp)
			if err != nil {
				log.Println("error", err)
			}
			fmt.Fprintln(ofp, ".")
		}
	}
	return nil
}
