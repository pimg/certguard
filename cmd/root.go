package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/internal/ports/models"
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "crl",
		Long:    "Crl Inspector (crl) can download and inspect x.509 Certificate Revocation Lists",
		Example: "crl",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := adapter.NewFileCache()
			if err != nil {
				return err
			}

			if _, err := tea.NewProgram(models.NewBaseModel(), tea.WithAltScreen()).Run(); err != nil {
				return err
			}
			return nil
		},
	}

	return rootCmd.Execute()
}
