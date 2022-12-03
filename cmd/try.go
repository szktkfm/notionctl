package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tryCmd)
}

func someFunc() error {
	return errors.New("this is error")
}

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "Try and possibly fail at something",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := someFunc(); err != nil {
			return err
		}
		return nil
	},
}
