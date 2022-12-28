package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"example.com/notion-go-cli/util"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mkuznets.com/go/tabwriter"
)

func init() {
	rootCmd.AddCommand(newCmdPush(&PushOptions{}, os.Stdout))
}

type PushOptions struct {
	DB          string
	Title       string
	Description string
	targetDB    notion.Database
	FilePath    string
	Out         io.Writer
	In          io.Reader
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdPush(o *PushOptions, writer io.Writer) *cobra.Command {
	// o := &PushOptions{}
	cmd := &cobra.Command{
		Use:   "push",
		Short: "push text",
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

	// cmd.Flags().StringVar(&o.Db, "db", "", "db id")
	cmd.Flags().StringVar(&o.Title, "title", o.Title, "title string")
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "description string")
	cmd.Flags().StringVarP(&o.FilePath, "file", "f", o.FilePath, "file path")
	return cmd
}

func (o *PushOptions) Complete(cmd *cobra.Command, args []string) error {
	o.DB = viper.GetString("db")
	client := notion.NewClient(getSecret())

	// --file - のとき stdinから読み込む
	if o.FilePath != "-" {
		o.In, _ = os.Open(o.FilePath)
	} else {
		o.In = cmd.InOrStdin()
	}

	o.targetDB, _ = client.FindDatabaseByID(context.Background(), o.DB)

	return nil
}

func (o *PushOptions) Run(cmd *cobra.Command, args []string) error {
	fmt.Println(args)

	client := notion.NewClient(getSecret())

	params := o.newCreatePageParams()

	page, err := client.CreatePage(context.Background(), params)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	util.NewCreatePageResponcePrinter(page, w).Print()
	w.Flush()
	return nil

}

func (o *PushOptions) newCreatePageParams() notion.CreatePageParams {
	dbPageProp := make(notion.DatabasePageProperties)

	for k, dp := range o.targetDB.Properties {
		// fmt.Printf("properties: %#v\n", k)
		// fmt.Printf("Properties value: %#v\n", dp)
		switch dp.Type {
		case notion.DBPropTypeTitle:
			dbPageProp[k] = notion.DatabasePageProperty{
				Title: stringToRichTexts(o.Title),
			}
		case notion.DBPropTypeRichText:
			dbPageProp[k] = notion.DatabasePageProperty{
				RichText: stringToRichTexts(o.Description),
			}
		}
	}

	return notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               o.DB,
		DatabasePageProperties: &dbPageProp,
		Children:               o.buildBlocksFromFile(),
	}
}

func (o *PushOptions) buildBlocksFromFile() []notion.Block {
	scanner := bufio.NewScanner(o.In)

	var blocks []notion.Block
	for scanner.Scan() {
		blocks = append(blocks,
			notion.ParagraphBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: scanner.Text(),
						},
					},
				},
			},
		)
	}
	return blocks
}

func stringToRichTexts(content string) []notion.RichText {
	return []notion.RichText{
		{
			Text: &notion.Text{
				Content: content,
			},
		},
	}
}
