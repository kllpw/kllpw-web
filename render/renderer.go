package render

import (
	"html/template"
	"io"
	"log"
)

type Page struct {
	Filename string
	Title string
	Header string
	ContentTitle string
	Content map[string]interface{}
}

// GetPageTemplate returns the populated template
func GetPageTemplate(page Page) *template.Template {
	t := template.New(page.Filename) // Create a template.
	t, _ = t.ParseFiles("./render/templates/" +  page.Filename) // Parse template file.
	return t
}
// WritePageToTemplate writes the page contents out via supplied template to io.Writer
func WritePageToTemplate(w io.Writer, page Page, t *template.Template) {
	err := t.Execute(w, page)
	if err != nil {
		log.Print(err.Error())
	}
}
