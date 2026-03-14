package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "skill",
		Short:   "Skill management tool",
		Long:    `A Go Module-inspired skill management tool for SafeClaw`,
		Version: fmt.Sprintf("%s (%s, built %s)", version, commit, date),
	}

	// 添加子命令
	rootCmd.AddCommand(
		initAddCmd(),
		initRemoveCmd(),
		initListCmd(),
		initInfoCmd(),
		initCacheCmd(),
		initInitCmd(),
	)

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
