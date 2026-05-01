package cmd

import (
	"fmt"
	"os"

	"github.com/YamasouA/mdview/internal/app"
	"github.com/YamasouA/mdview/internal/render"
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
		Use:   "mdview <file> [file...]",
		Short: "A terminal-native Markdown pager",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pages := make([]render.Page, 0, len(args))
			for _, path := range args {
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
				pages = append(pages, page)
			}

			model := app.NewModelWithPages(pages, render.RenderMarkdown)
			_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
			return err
		},
	}
	return cmd
}
