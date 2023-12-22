package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version.",
		Long:  "print version.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			Printer()
		},
	}

	workerCmd = &cobra.Command{
		Use:   "worker",
		Short: "run worker service.",
		Long:  "run worker service.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			RunWorker()
		},
	}
)

func Execute() {
	rootCmd := &cobra.Command{Use: "crawler"}
	rootCmd.AddCommand(masterCmd, workerCmd, versionCmd)
	err := rootCmd.Execute()
	if err != nil {
		println(fmt.Errorf("root cmd execute %w", err))
	}
}
