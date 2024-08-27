package cmd

import (
	"fmt"

	"github.com/farcaller/gonix"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

func nixToGo(val *gonix.Value) (interface{}, error) {
	switch val.Type() {
	case gonix.NixTypeThunk:
		return "<thunk>", nil
	case gonix.NixTypeInt:
		return val.GetInt()
	case gonix.NixTypeFloat:
		return val.GetFloat()
	case gonix.NixTypeBool:
		return val.GetBool()
	case gonix.NixTypeString:
		return val.GetString()
	case gonix.NixTypePath:
		return val.GetPath()
	case gonix.NixTypeNull:
		return "<null>", nil
	case gonix.NixTypeAttrs:
		a, err := val.GetAttrs()
		if err != nil {
			return nil, err
		}
		b := map[string]interface{}{}
		for k, v := range a {
			vv, err := nixToGo(v)
			if err != nil {
				return nil, err
			}
			b[k] = vv
		}
		return b, nil
	case gonix.NixTypeList:
		l, err := val.GetList()
		if err != nil {
			return nil, err
		}
		ll := []interface{}{}
		for _, v := range l {
			vv, err := nixToGo(v)
			if err != nil {
				return nil, err
			}
			ll = append(ll, vv)
		}
		return ll, nil
	case gonix.NixTypeFunction:
		return "<func>", nil
	case gonix.NixTypeExternal:
		return "<external>", nil
	}
	return nil, fmt.Errorf("unknown type")
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show info about a specific nixdockertag image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := gonix.NewContext()
		store, err := gonix.NewStore(ctx, "dummy", nil)
		if err != nil {
			return fmt.Errorf("failed to create a store: %w", err)
		}
		state := store.NewState(nil)

		val, err := state.EvalExpr("import "+args[0], ".")
		if err != nil {
			return fmt.Errorf("failed to eval: %w", err)
		}

		gval, err := nixToGo(val)
		if err != nil {
			return fmt.Errorf("failed to eval: %w", err)
		}

		fmt.Println(gval)
		return nil
	},
}
