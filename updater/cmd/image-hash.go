package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(imageHashCmd)
}

var imageHashCmd = &cobra.Command{
	Use:   "image-hash",
	Short: "Print the sha256 hash of the given image tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := args[0]
		manifest, err := crane.Manifest(image)
		if err != nil {
			return fmt.Errorf("failed to fetch the manifest: %w", err)
		}

		hash := sha256.New()
		hash.Write(manifest)
		hashedBytes := hash.Sum(nil)
		hashHex := hex.EncodeToString(hashedBytes)
		fmt.Println(hashHex)

		return nil
	},
}
