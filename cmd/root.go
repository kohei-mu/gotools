package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config Config

var globalFlags struct {
	language string
}

type Config struct {
	AppID string
}

var rootCmd = &cobra.Command{
	Use:   "gotool",
	Short: "gotool command",
	Long:  "gotool command.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalFlags.language, "lang", "l", "en", "Language : en, ja")
}
