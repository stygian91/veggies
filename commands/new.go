package commands

import (
	"os"
	"os/exec"

	"github.com/stygian91/veggies/internal/templates"
)

type fileEntry struct{ Path, Content string }

var mapping []fileEntry

func init() {
	m := []fileEntry{
		{Path: "/main.go", Content: templates.Main},
		{Path: "/go.mod", Content: templates.Gomod},
		{Path: "/go.sum", Content: templates.Gosum},
		{Path: "/sqlc.yaml", Content: templates.Sqlc},
		{Path: "/query.sql", Content: templates.Query},
		{Path: "/schema.sql", Content: templates.Schema},
		{Path: "/.env", Content: templates.EnvExample},
		{Path: "/handlers/greet.go", Content: templates.Greet},
		{Path: "/routes/routes.go", Content: templates.Routes},
	}
	mapping = m
}

func writeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func New(name string) error {
	cmd := exec.Command("mkdir", name)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("mkdir", name+"/handlers")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("mkdir", name+"/routes")
	if err := cmd.Run(); err != nil {
		return err
	}

	for _, entry := range mapping {
		if err := writeFile(name+entry.Path, entry.Content); err != nil {
			return err
		}
	}

	return nil
}
