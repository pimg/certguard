package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/adapter"
	"github.com/pimg/certguard/internal/ports/models"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enables debug logging to a file located in ~/.local/certguard/debug.log")
}

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "crl",
	Long:    "Crl Inspector (crl) can download and inspect x.509 Certificate Revocation Lists",
	Example: "crl",
	RunE:    runInteractiveCertGuard,
}

func runInteractiveCertGuard(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Flags().GetBool("debug")

	if debug {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}

		logDir := filepath.Join(homeDir, ".local", "share", "certguard")
		err = os.MkdirAll(logDir, 0o777)
		if err != nil {
			return err
		}

		f, err := tea.LogToFile(filepath.Join(logDir, "debug.log"), "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	cacheDir, err := adapter.NewFileCache()
	if err != nil {
		return err
	}
	log.Printf("file cache initialized at: %s", cacheDir)

	if _, err := tea.NewProgram(models.NewBaseModel(), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
