package main

import (
	"encoding/json"
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

// AppConfig defines the structure for the json configuration file.
type AppConfig struct {
	Username string `json:"username"`
}

// Index defines the structure of the index page.
type Index struct {
	Title    string
	Links    []string
	Year     int
	Username string
	Errors   map[string]string
}

// Page defines the structure of our data for a page.
type Page struct {
	Title string
	Body  []byte
}

var appConfig *AppConfig
var appConfigFile = "app-config.json"
var currentDir string
var currentYear = time.Now().Year()
var fileList []string
var templates = template.Must(template.ParseGlob("tmpl/*.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([A-Za-z-]*)$")

// buildIndex creates a list of files in the current directory for the Index page.
func buildIndex(ext string) []string {
	var files []string
	filepath.Walk("docs/", func(path string, f os.FileInfo, _ error) error {
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
	filename := "docs/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	// If there is an issue with reading the file.
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// loadAppConfig returns the app configuation file.
func loadAppConfig() error {
	config, err := ioutil.ReadFile(appConfigFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(config), &appConfig)
	if err != nil {
		return err
	}
	return nil
}

// updateAppConfig updates the app configuration file.
func updateAppConfig(newConfig *AppConfig) {
	if appConfig == nil {
		err := loadAppConfig()
		if err != nil {
			log.Fatal("Could not load the app configuration.")
		}
	}
	appConfig.Username = newConfig.Username
	jsonString, err := json.Marshal(appConfig)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(appConfigFile, jsonString, 0600)
}

// savePage is a function that saves a new page to text file.
func (p *Page) savePage() error {
	filename := "docs/" + p.Title + ".txt"
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
	switch r.Method {
	case "GET":
		fileList = buildIndex(".txt")
		i := &Index{Title: "Index", Links: fileList, Year: currentYear, Username: appConfig.Username}
		renderTemplate(w, "index", i)
	case "POST":
		r.ParseForm()
		username := r.Form["username"][0]
		if len(username) > 0 {
			var newConfig *AppConfig
			newConfig = new(AppConfig)
			newConfig.Username = username
			updateAppConfig(newConfig)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			i := &Index{Title: "Index", Links: fileList, Year: currentYear, Username: appConfig.Username}
			i.Errors = make(map[string]string)
			i.Errors["Username"] = "Please enter a username."
			renderTemplate(w, "index", i)
		}
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
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

func init() {
	// If the docs folder doesn't exist, we need to create it.
	if _, err := os.Stat("docs"); os.IsNotExist(err) {
		os.Mkdir("docs", 0777)
	}

	// If there isn't an app-config.json, we should create a new one.
	err := loadAppConfig()
	if err != nil {
		var newConfig *AppConfig
		newConfig = new(AppConfig)
		newConfig.Username = "New User"
		appConfig = newConfig
		updateAppConfig(newConfig)
	}
}

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	fmt.Println("Server running on http://localhost:4646/")
	fmt.Println("Press 'CTRL+C' to stop the server.")
	log.Fatal(http.ListenAndServe(":4646", nil))
}
