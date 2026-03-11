package config

import "fmt"

func Required(fields map[string]string) {
	for name, val := range fields {
		if val == "" {
			panic(fmt.Sprintf("config: missing required env var %s", name))
		}
	}
}
