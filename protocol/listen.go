package protocol

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/wtnb75/gonin-node/config"
)

func ListenAndServe(addr string, conf config.Config) error {
	name, err := os.Hostname()
	if err != nil {
		name = "gonin-node"
	}
	srv := Server{Addr: addr, Config: conf, Name: name}
	return srv.ListenAndServe()
}

type Server struct {
	Addr   string
	Name   string
	Config config.Config
}

type Conn struct {
	RemoteAddr  net.Addr
	Srv         *Server
	Conn        net.Conn
	Multigraph  bool
	DirtyConfig bool
}

func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		log.Printf("connect from %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
		connstruct := Conn{Conn: conn, RemoteAddr: conn.RemoteAddr(), Multigraph: false, DirtyConfig: false, Srv: srv}
		go connstruct.Handle(conn, conn)
	}
}

func (conn *Conn) Handle(inp io.Reader, outp io.Writer) {
	defer conn.Close()
	fmt.Fprintln(conn.Conn, "# munin node at", conn.Srv.Name)
	if err := conn.Repl(conn.Conn, conn.Conn); err != nil {
		if err != io.EOF {
			log.Fatalln("repl error", err)
		}
	}
}

func (conn *Conn) Close() {
	log.Println("close connection", conn.RemoteAddr)
	conn.Conn.Close()
}
