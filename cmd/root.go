package cmd

import (
	"fmt"
	"os"

	"github.com/YamasouA/mdless/internal/app"
	"github.com/YamasouA/mdless/internal/render"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mdless <file>",
		Short: "A terminal-native Markdown pager",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			info, err := os.Stat(path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return fmt.Errorf("%s is a directory", path)
			}

			page, err := render.RenderMarkdown(path)
			if err != nil {
				return err
			}

			model := app.NewModel(page, render.RenderMarkdown)
			_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
			return err
		},
	}
	return cmd
}
