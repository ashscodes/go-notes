# Go Notes

Go Notes is a [Go](https://golang.org/) app that allows the user to store personal notes among other features.

I built this to learn how to do a number of things in Go. Let me know if I am doing something wrong.

The front end is static HTML templates that use [Semantic UI](https://semantic-ui.com/)

This project is a work in progress.

## Go Packages Used

encoding/json, fmt, html/template, io/ioutil, log, net/http, os, path/filepath, regexp, strings, time

## Instructions For Use

```bash
git clone https://github.com/ashscodes/go-notes.git
cd go-notes
go build .

# Unix
go-notes

# Windows
.\go-notes.exe
```

Navigate to [http://localhost:4646](http://localhost:4646) and follow the instructions on the home page.

### Installing Go

[Download Go for your system.](https://golang.org/dl/)

## Immediate Worklist

- Expand HTML templates for Edit and View pages.

- Add an notes index page.

- Add the ability to delete notes.

- Refactor main.go.

## Future Plans

- Use Markdown instead of plain text.

- Ability to replicate files to cloud storage provider. 
