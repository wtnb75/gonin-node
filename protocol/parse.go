package protocol

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type Cmd func(conn *Conn, cmd string, arg []string, ofp io.Writer) error

var VersionString string = "0.1"

var cmdmap = map[string]Cmd{
	"fetch":   Fetch,
	"cap":     Cap,
	"config":  Config,
	"list":    List,
	"nodes":   Nodes,
	"version": Version,
	"quit":    Quit,
	".":       Quit,
}

func (conn *Conn) Repl(ifp io.Reader, ofp io.Writer) error {
	inp := bufio.NewReader(ifp)
	for {
		if line, _, err := inp.ReadLine(); err != nil {
			log.Println("read error", err)
			return err
		} else {
			l := strings.Fields(string(line))
			if len(l) == 0 {
				continue
			}
			log.Println("got command", l)
			if fn, ok := cmdmap[l[0]]; ok {
				start := time.Now()
				if err := fn(conn, l[0], l[1:], ofp); err != nil {
					if err != io.EOF {
						log.Println("error", err, time.Since(start))
					}
					return err
				}
				log.Println("done", l[0], time.Since(start))
			} else {
				keys := make([]string, 0, len(cmdmap))
				for k := range cmdmap {
					if k != "." {
						keys = append(keys, k)
					}
				}
				fmt.Fprintln(ofp, "# Unknown command. Try", strings.Join(keys, ", "))
			}
		}
	}
}
