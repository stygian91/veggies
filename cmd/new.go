package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stygian91/veggies/internal/templates"
)

type fileEntry struct{ Path, Content string }

var mapping []fileEntry = []fileEntry{
	{Path: "/main.go", Content: templates.Main},
	{Path: "/go.mod", Content: templates.Gomod},
	{Path: "/sqlc.yaml", Content: templates.Sqlc},
	{Path: "/query.sql", Content: templates.Query},
	{Path: "/schema.sql", Content: templates.Schema},
	{Path: "/.env", Content: templates.EnvExample},
	{Path: "/handlers/greet.go", Content: templates.Greet},
	{Path: "/routes/routes.go", Content: templates.Routes},
}
var subdirs []string = []string{
	"/handlers",
	"/routes",
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new project with veggies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(args[0]); err != nil {
			fmt.Printf("Error while running `new` command: %s\n", err)
		}
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
}

var moduleName *string

func init() {
	moduleName = newCmd.Flags().String("module-path", "veggie-app", "Module path for the project")
	rootCmd.AddCommand(newCmd)
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

func createDirs(name string) error {
	if err := os.Mkdir(name, 0755); err != nil {
		return fmt.Errorf("Error while creating main project directory: %w", err)
	}

	for _, subpath := range subdirs {
		dir := name + subpath
		if err := os.Mkdir(dir, 0755); err != nil {
			return fmt.Errorf("Error while creating subdirectory '%s': %w", dir, err)
		}
	}

	return nil
}

func runTidy(name string) (string, error) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = fmt.Sprintf("./%s", name)

	out, err := cmd.CombinedOutput()

	return string(out), err
}

// TODO:
// add info output
// maybe add options for db drivers
func run(name string) error {
	if err := createDirs(name); err != nil {
		return err
	}

	for _, entry := range mapping {
		path := name + entry.Path
		content := strings.ReplaceAll(entry.Content, "{{module}}", strings.TrimRight(*moduleName, "/"))
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("Error while writing template file '%s': %w", path, err)
		}
	}

	out, err := runTidy(name);
	fmt.Println(out)
	if  err != nil {
		return fmt.Errorf("Error while running `go mod tidy`: %w", err)
	}

	fmt.Printf("Successfully created new project '%s'\n", name)

	return nil
}
