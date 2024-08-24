package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/farcaller/nixdockertag/pkg/nix"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"
)

var (
	updateCommitFlag bool
)

func init() {
	updateCmd.Flags().BoolVar(&updateCommitFlag, "commit", false, "create a commit")
	rootCmd.AddCommand(updateCmd)
}

func updateOneImage(imagePath string, commit bool) error {
	info, err := nix.ParseNix(imagePath)
	if err != nil {
		return fmt.Errorf("failed to parse the image info: %w", err)
	}
	manifest, err := crane.Manifest(info.Image + ":" + info.FollowTag)
	if err != nil {
		return fmt.Errorf("failed to fetch the manifest: %w", err)
	}

	hash := sha256.New()
	hash.Write(manifest)
	hashedBytes := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashedBytes)

	if hashHex == info.Hash {
		fmt.Printf("No update for %s[%s]\n", info.Image, info.FollowTag)
		return nil
	}

	err = nix.WriteNix(imagePath, nix.ImageInfo{
		Image:     info.Image,
		FollowTag: info.FollowTag,
		Hash:      hashHex,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Updated %s[%s]: %s -> %s\n", info.Image, info.FollowTag, info.Hash, hashHex)

	if commit {
		r, err := git.PlainOpen(".")
		if err != nil {
			return fmt.Errorf("failed opening the git repository: %w", err)
		}
		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("failed getting the worktree: %w", err)
		}
		_, err = w.Add(imagePath)
		if err != nil {
			return fmt.Errorf("failed adding the file to commit: %w", err)
		}
		_, err = w.Commit(info.Image+": update to "+hashHex, &git.CommitOptions{})
		if err != nil {
			return fmt.Errorf("failed creating a commit: %w", err)
		}
	}
	return nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a single image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		imagePath, err := nix.ImageToPath(args[0])
		if err != nil {
			return err
		}

		return updateOneImage(imagePath, updateCommitFlag)
	},
}
