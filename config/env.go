package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

var env map[string]string

func read() (map[string]string, error) {
	return godotenv.Read()
}

func BootEnv() error {
	e, err := read()
	if err != nil {
		return fmt.Errorf("Error while booting env: %w", err)
	}

	env = e

	return nil
}

func GetEnv() map[string]string {
	return env
}
