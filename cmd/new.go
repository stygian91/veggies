package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/stygian91/veggies/internal/templates"
)

type rename struct{ Src, Dest string }

// used by the new command for the <project-path> cli argument
var moduleName *string

var renames []rename = []rename{
	{Src: "go.mod.tmpl", Dest: "go.mod"},
}

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

func runTidy(name string) (string, error) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = fmt.Sprintf("./%s", name)

	out, err := cmd.CombinedOutput()

	return string(out), err
}

// TODO:
// maybe add options for db drivers
func run(name string) error {
	fmt.Println("Creating files and directories...")

	subfs, err := fs.Sub(templates.Templates, "data")
	if err != nil {
		return fmt.Errorf("Error while getting template subdir: %w", err)
	}

	if err := os.CopyFS(name, subfs); err != nil {
		return fmt.Errorf("Error while copying templates: %w", err)
	}

	for _, r := range renames {
		src := name + "/" + r.Src
		dest := name + "/" + r.Dest
		if err := os.Rename(src, dest); err != nil {
			return fmt.Errorf("Error while renaming '%s' to '%s': %w", src, dest, err)
		}
	}

	fmt.Printf("Successfully created new project '%s'\n", name)

	return nil
}
