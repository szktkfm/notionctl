package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"example.com/notion-go-cli/util"
	"mkuznets.com/go/tabwriter"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCmdGet(&GetOptions{}, os.Stdout))

	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>"
	// viper.SetDefault("license", "apache")
}

type GetOptions struct {
	DB   string
	Page string
	Wide bool
	Out  io.Writer
}

type utilFactory struct {
	client http.Client
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdGet(o *GetOptions, writer io.Writer) *cobra.Command {
	// option := &GetOptions{Out: writer}
	o.Out = writer
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

	cmd.Flags().StringVar(&o.Page, "page", o.Page, "page id") //このdefault値を引数で渡したoptionの値にする
	cmd.Flags().StringVar(&o.DB, "db", o.DB, "db id")
	cmd.Flags().BoolVar(&o.Wide, "wide", o.Wide, "wide print")
	return cmd
}

func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.DB = viper.GetString("db") // ToDo: db引数の値を優先する。でなければviperで読み取る。
	// viperをつかってconfigから読み取る(configから読み取らずviper packageのmethodをそのまま使ってしまうと、singletonみたいになって、実行順序よって結果が変わるかも)
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
	queryResult, err := client.QueryDatabase(context.Background(), o.DB, nil)

	if err != nil {
		return err
	}
	// print
	// TODO: tabwriterを別のpackageにする
	const padding = 4
	// w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	w := tabwriter.NewWriter(o.Out, 4, 0, padding, ' ', tabwriter.TabIndent)
	// util.PrintDatabaseQueryResponce(queryResult, w)
	util.NewDatabaseQueryResponcePrinter(queryResult, w).Print()
	w.Flush()
	// fmt.Println(queryResult)
	return nil

}

func getSecret() string {
	return viper.GetString("secret")
}
