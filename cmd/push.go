package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"example.com/notion-go-cli/util"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	o.DB = viper.GetString("db") // ToDo: 面倒だからenv variableで読みだそうかな
	client := notion.NewClient(getSecret())

	// 一回dbを情報をgetしてきて、そこからparameterをbuildする。
	o.targetDB, _ = client.FindDatabaseByID(context.Background(), o.DB)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (o *PushOptions) Run(cmd *cobra.Command, args []string) error {
	fmt.Println(args)

	client := notion.NewClient(getSecret())
	params := o.newCreatePageParams() //ToDo:  param構造体を作る
	// param.build みたいな感じでparameterをbuildして

	// fmt.Printf("params: %#v\n", params)
	page, err := client.CreatePage(context.Background(), params)
	// fmt.Printf("bufer: %s", buf)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return err
	}
	// print
	// TODO: tabwriterを別のpackageにする
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
	// printPage(page, w)
	// fmt.Printf("%#v\n", page)
	util.NewCreatePageResponcePrinter(page, w).Print()
	w.Flush()
	// fmt.Println(queryResult)
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
				Title: getRitchText(o.Title),
			}
		case notion.DBPropTypeRichText:
			dbPageProp[k] = notion.DatabasePageProperty{
				RichText: getRitchText(o.Description),
			}
		}
	}

	fp, _ := os.Open(o.FilePath)
	// if err != nil {
	// 	return err
	// }
	scanner := bufio.NewScanner(fp)

	// children :=

	// return notion.CreatePageParams{
	// 	ParentType:             notion.ParentTypeDatabase,
	// 	ParentID:               o.Db,
	// 	DatabasePageProperties: &dbPageProp,
	// 	Children: []notion.Block{
	// 		notion.Heading1Block{
	// 			RichText: []notion.RichText{
	// 				{
	// 					Text: &notion.Text{
	// 						Content: "testtest",
	// 					},
	// 				},
	// 				{
	// 					Text: &notion.Text{
	// 						Content: "2行目",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	return notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               o.DB,
		DatabasePageProperties: &dbPageProp,
		Children:               convert(scanner),
	}
}

func convert(scanner *bufio.Scanner) []notion.Block {
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

func getRitchText(content string) []notion.RichText {
	return []notion.RichText{
		{
			Text: &notion.Text{
				Content: content,
			},
		},
	}
}
