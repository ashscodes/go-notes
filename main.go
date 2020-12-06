package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Index defines the structure of the index page.
type Index struct {
	Title    string
	Links    []string
	Year     int
	Username string
}

// Page defines the structure of our data for a page.
type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "index.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// buildIndex creates a list of files in the current directory for the Index page.
func buildIndex(ext string) []string {
	var files []string
	filepath.Walk(".", func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ext {
				files = append(files, strings.TrimSuffix(f.Name(), ext))
			}
		}
		return nil
	})
	return files
}

// makeHandler takes a handler function and returns a type of http.HandlerFunc
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// renderTemplate renders a HTML page from the template files included.
func renderTemplate(w http.ResponseWriter, tmpl string, i interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// loadPage returns a page to the application from a text file.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	// If there is an issue with reading the file.
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// savePage is a function that saves a new page to text file.
func (p *Page) savePage() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// editHandler display a form where a user can edit the page.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// faviconHandler serves the favicon icon.
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "img/favicon.ico")
}

// listHandler displays a list of browsable pages on the root page.
func listHandler(w http.ResponseWriter, r *http.Request) {
	fileList := buildIndex(".txt")
	currentYear := time.Now().Year()
	i := &Index{Title: "Index", Links: fileList, Year: currentYear, Username: "User"}
	renderTemplate(w, "index", i)
}

// saveHandler handles the submission of forms that have are edited.
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.savePage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// viewHandler allows users to view a notes page.
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	fmt.Println("Server running on http://localhost:4646/")
	fmt.Println("Press 'CTRL+C' to stop the server.")
	log.Fatal(http.ListenAndServe(":4646", nil))
}
