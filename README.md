# notionctl
A toy notion cli

## Installation
```
go install github.com/szktkfm/notionctl...
```

## Usage

Create a Notion Integration. Refer https://developers.notion.com/docs.


Set the API key and the ID of the Database to be integrated into an environment variable.
```
export NOTION_API_KEY=<your notion api key>
export NOTION_DATABASE=<integrated db id>
```

You can display the values from the Database in a table-like format as follows
```bash
$ notionctl get
NAME               AGE     FOOD GROUP    TAGS     DESCRIPTION
Mandarin orange    5d4h    Fruit         -        small citrus tree fruit
Tuscan Kale        32d     Vegetable     salad    A dark green leafy veget
```


You can create pages in the database from markdown files or raw text.

```bash
$ notionctl push --title readme --file readme.md

$ notionctl push --title todo --description "Buy toothpaste."

$ cat <<EOF | notionctl --title todo --file -
Buy toothpaste.
EOF

```