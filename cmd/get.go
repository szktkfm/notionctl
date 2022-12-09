package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"example.com/notion-go-cli/util"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCmdGet())

	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>"
	// viper.SetDefault("license", "apache")
}

type GetOptions struct {
	Db   string
	Page string
	Wide bool
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdGet() *cobra.Command {
	o := &GetOptions{}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "config on .notion-go",
		RunE: func(cmd *cobra.Command, args []string) error {
			//debug
			// fmt.Println(getSecret())
			// flagの穴埋め
			o.Complete(cmd, args)
			// run
			o.Run(cmd, args)
			return nil
		},
	}

	cmd.Flags().StringVar(&o.Page, "page", "", "page id")
	cmd.Flags().StringVar(&o.Db, "db", "", "db id")
	cmd.Flags().BoolVar(&o.Wide, "wide", false, "wide print")
	return cmd
}

func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Db = viper.GetString("db") // ToDo: db引数の値を優先する。でなければviperで読み取る。
	// viperをつかってconfigから読み取る(configから読み取らずviper packageのmethodをそのまま使ってしまうと、singletonみたいになって、実行順序よって結果が変わるかも)
	o.Page, _ = cmd.Flags().GetString("page")

	return nil
}

func (o *GetOptions) Run(cmd *cobra.Command, args []string) error {
	if len(o.Page) > 0 {
		// ToDo: get markdown and build name
		name := "sample"
		fmt.Printf("get markdown to %s.md\n", name)
		return nil
	}
	client := notion.NewClient(getSecret())
	queryResult, err := client.QueryDatabase(context.Background(), o.Db, nil)

	if err != nil {
		return err
	}
	// print
	// TODO: tabwriterを別のpackageにする
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
	// util.PrintDatabaseQueryResponce(queryResult, w)
	util.NewDatabaseQueryResponcePrinter(queryResult, w).Print()
	w.Flush()
	// fmt.Println(queryResult)
	return nil

}

func readConfig() {
	fmt.Println(viper.GetString("secret"))
}

func getSecret() string {
	return viper.GetString("secret")
}
