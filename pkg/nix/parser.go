package nix

import (
	"fmt"
	"os"

	"github.com/farcaller/gonix"
)

type ImageInfo struct {
	Image     string
	FollowTag string
	Hash      string
}

func valToString(as map[string]*gonix.Value, key string) (string, error) {
	val, exists := as[key]
	if !exists {
		return "", fmt.Errorf("key %s doesn't exist in the attrset", key)
	}
	return val.GetString()
}

func ImageToPath(image string) (string, error) {
	path := "images/" + image + ".nix"
	if _, err := os.Stat("./" + path); err != nil {
		return "", fmt.Errorf("image not found: %w", err)
	}
	return path, nil
}

func ParseNix(imagePath string) (*ImageInfo, error) {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a store: %w", err)
	}

	state := store.NewState(nil)
	if state == nil {
		return nil, fmt.Errorf("failed to create nix state")
	}

	curdir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get the current directory: %w", err)
	}

	val, err := state.EvalExpr("import ./"+imagePath, curdir)
	if err != nil {
		return nil, fmt.Errorf("failed to eval: %w", err)
	}

	if val.Type() != gonix.NixTypeAttrs {
		return nil, fmt.Errorf("expected an attrset, got %v", val.Type())
	}

	attrs, err := val.GetAttrs()
	if err != nil {
		return nil, fmt.Errorf("failed to eval: %w", err)
	}

	var ii ImageInfo

	ii.Image, err = valToString(attrs, "image")
	if err != nil {
		return nil, err
	}
	ii.FollowTag, err = valToString(attrs, "followTag")
	if err != nil {
		return nil, err
	}
	ii.Hash, err = valToString(attrs, "hash")
	if err != nil {
		return nil, err
	}
	return &ii, nil
}
