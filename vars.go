package main

import (
	"log"
	"strings"
)

type Vars []string

func (v *Vars) String() string {
	return ""
}

func (v *Vars) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func parseVars(vars []string) map[string]string {
	varsMap := make(map[string]string)
	for _, kv := range vars {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("invalid var format: %s, expected key=value", kv)
		}
		key, value := parts[0], parts[1]
		varsMap[key] = value
	}
	return varsMap
}
