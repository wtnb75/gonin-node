package config

import (
	"bufio"
	"io"
	"log"
	"path/filepath"

	"github.com/google/shlex"
)

type Config struct {
	kv      map[string]string
	parts   map[string]map[string]string
	plugins map[string]Plugin
}

func (cfg *Config) ReadConfig(ifp io.Reader) error {
	if cfg.kv == nil {
		cfg.kv = map[string]string{}
	}
	if cfg.parts == nil {
		cfg.parts = map[string]map[string]string{}
	}
	input := bufio.NewReader(ifp)
	for {
		line, _, err := input.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		tokens, err := shlex.Split(string(line))
		if err != nil {
			return err
		}
		if len(tokens) == 2 {
			cfg.kv[tokens[0]] = tokens[1]
		} else if len(tokens) != 0 {
			log.Println("ignore line: ", tokens)
		}
	}
}

func (cfg *Config) Get(key, def string) string {
	if v, ok := cfg.kv[key]; ok {
		return v
	} else {
		return def
	}
}

func (cfg *Config) GetWithPart(part, key, def string) string {
	for m := range cfg.parts {
		if ok, err := filepath.Match(m, part); ok {
			if v, ok := cfg.parts[m][key]; ok {
				return v
			} else {
				return def
			}
		} else {
			log.Println("error", err)
		}
	}
	return def
}
