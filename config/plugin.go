package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

type Plugin struct {
	CommandPath string
	MultiGraph  bool
	Env         map[string]string
	Uid         uint32
	Gid         uint32
}

func username2uid(name string) (uint32, error) {
	if user, err := user.Lookup(name); err != nil {
		return 0, err
	} else {
		if uid, err := strconv.ParseUint(user.Uid, 8, 32); err != nil {
			return 0, err
		} else {
			return uint32(uid), nil
		}
	}
}

func groupname2gid(name string) (uint32, error) {
	if grp, err := user.LookupGroup(name); err != nil {
		return 0, err
	} else {
		if uid, err := strconv.ParseUint(grp.Gid, 8, 32); err != nil {
			return 0, err
		} else {
			return uint32(uid), nil
		}
	}
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
			plugin := Plugin{CommandPath: fn}
			username := cfg.GetWithPart(file.Name(), "user", "munin")
			groupname := cfg.GetWithPart(file.Name(), "group", "munin")
			if uid, err := username2uid(username); err != nil {
				log.Println("user name error:", username, err)
			} else {
				plugin.Uid = uid
			}
			if gid, err := groupname2gid(groupname); err != nil {
				log.Println("group name error:", groupname, err)
			} else {
				plugin.Gid = gid
			}
			plugin.Env = map[string]string{}
			for k, v := range cfg.ListWithPart(file.Name(), "env.") {
				plugin.Env[k[4:]] = v
			}
			var stdout bytes.Buffer
			if err := plugin.ShowConfig(false, &stdout); err != nil {
				log.Println("cmd error", err)
				log.Println("output", strconv.Quote(stdout.String()))
				continue
			}
			confstr := strings.TrimSpace(stdout.String())
			plugin.MultiGraph = strings.HasPrefix(confstr, "multigraph")
			log.Printf("plugin: %s %+v", file.Name(), plugin)
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
	for k, v := range p.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdin = nil
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: p.Uid, Gid: p.Gid}
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
