package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/internal/adapter/db"
	"github.com/pimg/certguard/internal/ports/models"
	cmds "github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/pkg/domain/crl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enables debug logging to a file located in ~/.local/certguard/debug.log")
}

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "certguard",
	Long:    "Certguard can download, store and inspect x.509 Certificate Revocation Lists",
	Example: "certguard",
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

	cacheDir, err := determineCacheDir()
	if err != nil {
		return err
	}

	dbConnection, err := db.NewDBConnection(cacheDir)
	if err != nil {
		return err
	}

	libsqlStorage := db.NewLibSqlStorage(dbConnection)
	defer func() {
		err := libsqlStorage.CloseDB()
		if err != nil {
			log.Printf("could not close database: %v", err)
		}
	}()

	err = libsqlStorage.InitDB(context.Background())
	if err != nil {
		return err
	}

	storage, err := crl.NewStorage(libsqlStorage, cacheDir)
	if err != nil {
		return err
	}

	log.Printf("cache initialized at: %s", cacheDir)

	commands := cmds.NewCommands(storage)

	if _, err := tea.NewProgram(models.NewBaseModel(commands), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func determineCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not create file path to User home dir, Cache will not be enabled")
	}

	return filepath.Join(homeDir, ".cache", "certguard"), nil
}
