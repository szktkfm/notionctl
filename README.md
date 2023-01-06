# notionctl
A toy notion cli

## Installation
```
go install github.com/szktkfm/notionctl...
```

## Usage

Create a Notion Integration. 

- https://developers.notion.com/docs


Set the API key and the ID of the Database to be integrated into an environment variable.
```
export NOTION_API_KEY=<your notion api key>
export NOTION_DATABASE=<integrated db id>
```

You can display the values from the Database in a table-like format as follows
```Bash
$ notionctl get
NAME               AGE     FOOD GROUP    TAGS     DESCRIPTION
Mandarin orange    5d4h    Fruit         -        small citrus tree fruit
Tuscan Kale        32d     Vegetable     salad    A dark green leafy veget
```

You can create a page in the Database. You can create a page paragraph by reading from a file or stdin.
```Bash
$ notionctl push --title todo --description "Buy toothpaste."

$ notionctl push --title todo --file todo.txt

$ cat <<EOF | notionctl --title todo --file
Buy toothpaste.
EOF
```
