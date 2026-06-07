package main

import (
	"os"
	"strings"
)

func ParseEnv() map[string]string {
	res := map[string]string{}

	rows := []string{}
	if bytes, err := os.ReadFile(".env"); err == nil {
		lines := strings.Split(string(bytes), "\n")
		rows = append(rows, lines...)
	}
	rows = append(rows, os.Environ()...)

	for _, s := range rows {
		parts := strings.Split(s, "=")
		res[parts[0]] = strings.Join(parts[1:], "=")
	}

	return res
}
