package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pimg/certguard/config"
	"github.com/pimg/certguard/internal/adapter/db"
	"github.com/pimg/certguard/internal/ports/models"
	cmds "github.com/pimg/certguard/internal/ports/models/commands"
	"github.com/pimg/certguard/internal/ports/models/styles"
	"github.com/pimg/certguard/pkg/domain/crl"
	"github.com/spf13/cobra"
)

var v *config.ViperConfig

func init() {
	v = config.NewViperConfig()
	rootCmd.PersistentFlags().BoolVarP(&v.Config().Log.Debug, "debug", "d", false, "enables debug logging to a file located in ~/.local/certguard/debug.log")
	rootCmd.PersistentFlags().StringVarP(&v.Config().Theme.Name, "theme", "t", "dracula", "set the theme of the application. Allowed values: 'dracula', 'gruvbox'")

	// bind Cobra flags to viper config
	_ = v.BindPFlag("config.theme.name", rootCmd.PersistentFlags().Lookup("theme"))
	_ = v.BindPFlag("config.log.debug", rootCmd.PersistentFlags().Lookup("debug"))
}

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "certguard",
	Long:    "Certguard can download, store and inspect x.509 Certificate Revocation Lists",
	Example: "certguard",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return v.InitializeConfig()
	},
	RunE: runInteractiveCertGuard,
}

func runInteractiveCertGuard(cmd *cobra.Command, args []string) error {
	debug := v.Config().Log.Debug
	theme := v.Config().Theme.Name

	if debug {
		logDir, err := logDir()
		if err != nil {
			return err
		}

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

	cacheDir, err := cacheDir()
	if err != nil {
		return err
	}

	importDir, err := importDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0o775)
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
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

	storage, err := crl.NewStorage(libsqlStorage, cacheDir, importDir)
	if err != nil {
		return err
	}

	log.Printf("cache initialized at: %s", cacheDir)

	styles.NewStyles(theme)

	commands := cmds.NewCommands(storage)

	if _, err := tea.NewProgram(models.NewBaseModel(commands), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func cacheDir() (string, error) {
	if v.Config().CacheDirectory != "" {
		return v.Config().CacheDirectory, nil
	}

	return defaultDir()
}

func importDir() (string, error) {
	if v.Config().ImportDirectory != "" {
		return v.Config().ImportDirectory, nil
	}

	dir, err := defaultDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "import"), nil
}

func logDir() (string, error) {
	if v.Config().Log.Directory != "" {
		return v.Config().Log.Directory, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not create file path to User home dir, logging will not be enabled")
	}

	return filepath.Join(homeDir, ".local", "share", "certguard"), nil
}

func defaultDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not create file path to User home dir, Cache will not be enabled")
	}

	return filepath.Join(homeDir, ".cache", "certguard"), nil
}
