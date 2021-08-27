package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"strconv"
)

type Plugin struct {
	CommandPath string
	MultiGraph  bool
}

func (cfg *Config) ReadPlugins(dirname string) error {
	if cfg.plugins == nil {
		cfg.plugins = map[string]Plugin{}
	}
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if (file.Mode() & 0100) != 0 {
			fn := filepath.Join(dirname, file.Name())
			log.Println("plugin:", fn)
			cmd := exec.Command(fn, "config")
			cmd.Env = baseenv
			cmd.Stdin = nil
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				log.Println("cmd error", err)
				log.Println("output", strconv.Quote(stdout.String()))
				log.Println("error", strconv.Quote(stderr.String()))
				continue
			}
			confstr := strings.TrimSpace(stdout.String())
			mgraph := strings.HasPrefix(confstr, "multigraph")
			plugin := Plugin{CommandPath: fn, MultiGraph: mgraph}
			log.Printf("plugin: %+v", plugin)
			cfg.plugins[file.Name()] = plugin
		}
	}
	return nil
}

func (cfg *Config) ListPlugin(multigraph bool) []string {
	res := []string{}
	for p := range cfg.plugins {
		if !multigraph && cfg.plugins[p].MultiGraph {
			continue
		}
		res = append(res, p)
	}
	sort.Strings(res)
	return res
}

func (cfg *Config) GetPlugin(name string) (res Plugin, err error) {
	res, ok := cfg.plugins[name]
	if !ok {
		err = fmt.Errorf("no such plugin: %s", name)
	}
	return
}

var baseenv = []string{
	"PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin",
	"MUNIN_LIBDIR=/usr/share/munin",
}

func (p *Plugin) Execute(output io.Writer) error {
	cmd := exec.Command(p.CommandPath)
	cmd.Env = baseenv
	cmd.Stdin = nil
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Plugin) ShowConfig(dirtyconfig bool, output io.Writer) error {
	cmd := exec.Command(p.CommandPath, "config")
	cmd.Env = baseenv
	if dirtyconfig {
		cmd.Env = append(cmd.Env, "MUNIN_CAP_DIRTYCONFIG=1")
	}
	log.Printf("run config with dirtyconfig=%v", dirtyconfig)
	cmd.Stdin = nil
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
