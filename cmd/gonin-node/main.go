package main

import (
	"io"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/wtnb75/gonin-node/config"
	"github.com/wtnb75/gonin-node/protocol"
)

type Option struct {
	Verbose    bool   `short:"v" long:"verbose" description:"set verbose"`
	ListenAddr string `short:"l" long:"listen-address" default:":4949"`
	Dir        string `short:"d" long:"plugin-dir" default:"/etc/munin/plugins"`
	ConfFile   string `short:"c" long:"config" default:"/etc/munin/munin-node.conf"`
	Repl       bool   `long:"repl-mode"`
}

var opts Option

func main() {
	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		return
	}
	log.Printf("opts: %+v, args: %+v", opts, args)
	conf := config.Config{}
	cfp, err := os.Open(opts.ConfFile)
	if err != nil {
		log.Fatalln("open", err)
	}
	err = conf.ReadConfig(cfp)
	if err != nil {
		log.Fatalln("read config", err)
	}
	cfp.Close()
	if err := conf.ReadPlugins(opts.Dir); err != nil {
		log.Fatalln("cannot read plugin", err)
	}
	if opts.Repl {
		srv := protocol.Server{Name: "local-repl", Config: conf}
		conn := protocol.Conn{Srv: &srv}
		err := conn.Repl(os.Stdin, os.Stdout)
		if err != nil && err != io.EOF {
			log.Fatalln("error", err)
		}
	} else {
		err := protocol.ListenAndServe(opts.ListenAddr, conf)
		if err != nil {
			log.Fatalln("error", err)
		}
	}
}
