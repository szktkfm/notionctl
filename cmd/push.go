package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/szktkfm/notionctl/util"
	"github.com/yuin/goldmark"
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

func newCmdPush(o *PushOptions, writer io.Writer) *cobra.Command {
	o.Out = writer
	cmd := &cobra.Command{
		Use:   "push",
		Short: "push text",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}
			err = o.Run(cmd, args)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&o.Title, "title", o.Title, "title string")
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "description string")
	cmd.Flags().StringVarP(&o.FilePath, "file", "f", o.FilePath, "file path")
	return cmd
}

func (o *PushOptions) Complete(cmd *cobra.Command, args []string) error {
	db, err := getDBID()
	if err != nil {
		return err
	}
	o.DB = db

	secret, err := getSecret()
	if err != nil {
		return err
	}
	client := notion.NewClient(secret)

	if len(o.Title) == 0 {
		return fmt.Errorf("required flag(s) title not set")
	}

	// If --file option is set to -, read from stdin
	if o.FilePath != "-" {
		o.In, _ = os.Open(o.FilePath)
		// TODO close o.In
	} else {
		o.In = cmd.InOrStdin()
	}

	o.targetDB, _ = client.FindDatabaseByID(context.Background(), o.DB)
	return nil
}

func (o *PushOptions) Run(cmd *cobra.Command, args []string) error {

	secret, err := getSecret()
	if err != nil {
		return err
	}

	client := notion.NewClient(secret)
	params := o.newCreatePageParams()
	_, err = client.CreatePage(context.Background(), params)

	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "%s is created\n", o.Title)

	return nil
}

func (o *PushOptions) newCreatePageParams() notion.CreatePageParams {
	dbPageProp := make(notion.DatabasePageProperties)

	for k, dp := range o.targetDB.Properties {
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
		Children:               o.buildBlocksFromInput(),
	}
}

func (o *PushOptions) buildBlocksFromInput() []notion.Block {

	// TODO: markdownからパースする

	source, _ := io.ReadAll(o.In)

	var buf bytes.Buffer
	buf.WriteString(`{"results": [`)
	md := goldmark.New(
		goldmark.WithExtensions(util.NotionExtension),
	)
	md.Convert(source, &buf)

	buf.Truncate(buf.Len() - 1)
	buf.WriteString(`]}`)

	// fmt.Printf("%v", buf.String())
	var param notion.BlockChildrenResponse
	json.Unmarshal(buf.Bytes(), &param)

	fmt.Printf("%v", buf.String())

	return param.Results
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
