# Marks

## Bookmarks on the command line

### Still under development! But it works well enough so to have a crack at it now:

1) [Install go](https://golang.org/doc/install)

2) Install marks:
```
go get github.com/tomguerney/marks
echo "contentPath: $HOME" > $HOME/.marks.yaml
touch $HOME/bookmarks.yaml
```

3) Try it out:
```
marks add "Abc News" --url www.abc.net.au --tag news --tag "current affairs"
marks open news
```
### Usage
```
Usage:
  marks [command]

Available Commands:
  add         Add a bookmark
  copy        Copy a bookmark to the clipboard
  delete      Delete a bookmark
  help        Help about any command
  open        Open a url in a browser
  update      Update a bookmark

Flags:
      --config string   config file (default is $HOME/.marks.yaml)
      --debug           output debug logs
  -h, --help            help for marks

Use "marks [command] --help" for more information about a command.
```