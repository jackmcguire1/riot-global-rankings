package utils

import (
	"os"
	"strconv"
)

func EnvironmentWithDefaultInt(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}

	res, _ := strconv.Atoi(v)
	return res
}
