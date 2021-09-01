package protocol

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wtnb75/gonin-node/config"
)

func TestListEmpty(t *testing.T) {
	conn := Conn{Srv: &Server{Name: "hello"}}
	ofp := bytes.Buffer{}
	if err := List(&conn, "list", []string{}, &ofp); err != nil {
		t.Error("list error", err)
	}
}

func make_plugin() string {
	dir, err := ioutil.TempDir("", "plugins")
	if err != nil {
		log.Panic("tmpdir", err)
	}
	plugin1 := filepath.Join(dir, "plugin1")
	content1 := `#! /bin/sh
echo "hello world"
`
	if err := ioutil.WriteFile(plugin1, []byte(content1), 0755); err != nil {
		log.Panic("output", err)
	}
	plugin2 := filepath.Join(dir, "plugin2")
	content2 := `#! /bin/sh
echo "multigraph hello world"
`
	if err := ioutil.WriteFile(plugin2, []byte(content2), 0755); err != nil {
		log.Panic("output", err)
	}
	return dir
}

func TestList(t *testing.T) {
	conf := config.Config{}
	dir := make_plugin()
	if err := conf.ReadPlugins(dir); err != nil {
		t.Error("readplugins", err)
	}
	conn := Conn{Srv: &Server{Name: "hello", Config: conf}}
	ofp := bytes.Buffer{}
	if err := List(&conn, "list", []string{}, &ofp); err != nil {
		t.Error("list error", err)
	}
	if !strings.Contains(ofp.String(), "plugin1") {
		t.Error("not contain: plugin1")
	}
	if strings.Contains(ofp.String(), "plugin2") {
		t.Error("contain: plugin2")
	}
	os.RemoveAll(dir)
}

func TestListMultigraph(t *testing.T) {
	conf := config.Config{}
	dir := make_plugin()
	if err := conf.ReadPlugins(dir); err != nil {
		t.Error("readplugins", err)
	}
	conn := Conn{Srv: &Server{Name: "hello", Config: conf}}
	ofp := bytes.Buffer{}
	if err := Cap(&conn, "cap", []string{"multigraph"}, &ofp); err != nil {
		t.Error("cap error", err)
	}
	if err := List(&conn, "list", []string{}, &ofp); err != nil {
		t.Error("list error", err)
	}
	if !strings.Contains(ofp.String(), "plugin1") {
		t.Error("not contain: plugin1", ofp.String())
	}
	if !strings.Contains(ofp.String(), "plugin2") {
		t.Error("not contain: plugin2", ofp.String())
	}
	os.RemoveAll(dir)
}
