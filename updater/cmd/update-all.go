package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	updateAllCommitFlag bool
)

func init() {
	updateAllCmd.Flags().BoolVar(&updateAllCommitFlag, "commit", false, "create a commit")
	rootCmd.AddCommand(updateAllCmd)
}

var updateAllCmd = &cobra.Command{
	Use:   "update-all",
	Short: "Update all images",
	RunE: func(cmd *cobra.Command, args []string) error {
		images, err := filepath.Glob("images/*.nix")
		if err != nil {
			return fmt.Errorf("failed to glob images: %w", err)
		}

		for _, imagePath := range images {
			err = updateOneImage(imagePath, updateAllCommitFlag)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
