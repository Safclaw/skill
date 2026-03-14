package main

import (
	"fmt"

	"github.com/Safclaw/skill/pkg/cache"
	"github.com/Safclaw/skill/pkg/config"
	"github.com/spf13/cobra"
)

var (
	cacheAllFlag bool
)

func initCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage skill cache",
		Long:  `Manage skill cache.`,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "clean",
			Short: "Clean cache",
			Long:  `Clean the skill cache.`,
			RunE:  runCacheClean,
		},
		&cobra.Command{
			Use:   "verify",
			Short: "Verify cache integrity",
			Long:  `Verify the integrity of cached skills.`,
			RunE:  runCacheVerify,
		},
	)

	cmd.PersistentFlags().BoolVar(&cacheAllFlag, "all", false, "Clean all cache")

	return cmd
}

func runCacheClean(cmd *cobra.Command, args []string) error {
	cfg := config.LoadConfig()
	c := cache.NewManager(cfg.CacheDir)

	if err := c.Clean(cacheAllFlag); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	fmt.Println("Cache cleaned successfully")
	return nil
}

func runCacheVerify(cmd *cobra.Command, args []string) error {
	cfg := config.LoadConfig()

	// TODO: 实现缓存验证逻辑
	fmt.Println("Cache verification not yet implemented")
	fmt.Printf("Cache directory: %s\n", cfg.CacheDir)

	return nil
}
