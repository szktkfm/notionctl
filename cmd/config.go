package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	User string `yaml:secret`
	DB   string `yaml:db`
}

func init() {
	rootCmd.AddCommand(newCmdConfig())

	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "config on .notion-go",
		RunE: func(cmd *cobra.Command, _ []string) error {

			fmt.Println("config on .notion-go")

			return nil
		},
	}

	cmd.AddCommand(newCmdConfigSet())
	cmd.AddCommand(newCmdConfigView())
	return cmd

}

func newCmdConfigSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set config on .notion-go",
		RunE: func(cmd *cobra.Command, _ []string) error {

			value, err := cmd.Flags().GetString("secret")
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			viper.Set("secret", value)

			// write by viper
			fmt.Println(value)

			value2, err := cmd.Flags().GetString("db")

			if err != nil {
				return fmt.Errorf("%v", err)
			}
			// write by viper
			fmt.Println(value2)
			viper.Set("db", value2)

			return viper.WriteConfig()
		},
	}
	cmd.Flags().StringP("secret", "s", "", "notion integration secret")
	cmd.Flags().StringP("db", "", "", "notion db")
	// viper.BindPFlag("secret", cmd.PersistentFlags().Lookup("secret"))
	// viper.BindPFlag("db", cmd.PersistentFlags().Lookup("db"))
	return cmd

}

func newCmdConfigView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "set config on .notion-go",
		RunE: func(cmd *cobra.Command, _ []string) error {

			//ToDo: configを一行ずつ読みだしてprintする

			// sec := viper.GetString("secret")
			// fmt.Println(sec)

			return nil
		},
	}
	return cmd
}
