package env

import (
	"fmt"

	"github.com/joho/godotenv"
)

var env map[string]string

func Read() (map[string]string, error) {
  return godotenv.Read()
}

func Boot() error {
  e, err := Read()
  if err != nil {
    return fmt.Errorf("Error while booting env: %w", err)
  }

  env = e

  return nil
}

func Get() map[string]string {
  return env
}
