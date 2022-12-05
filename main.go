package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

var (
	srcDir = flag.String("src", "./seiten", "Inhalte-Dir.")
	tmpDir = flag.String("tmp", "./templates", "Template-Dir.")
)

type Page struct {
	Title       string
	LastChanged time.Time
	Content     template.HTML
}

type Pages []Page

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*.templ.html")
	router.Static("/static", "./static")

	router.GET("/", indexHandler)
	router.GET("/page/:topic", blogHandler)
	log.Print("Listening on :9000 ....")
	err := router.Run(":9000")
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(c *gin.Context) {
	const template = "index" // definiert in "index.templ.html"
	ps, err := loadPages(*srcDir)
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, template, ps)
}

func blogHandler(c *gin.Context) {
	const template = "page" // definiert in "page.templ.html"
	f := c.Param("topic")
	fpath := filepath.Join(*srcDir, f)
	p, err := loadPage(fpath)
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, template, p)
}

func loadPage(fpath string) (Page, error) {
	var p Page
	fi, err := os.Stat(fpath)
	if err != nil {
		return p, fmt.Errorf("loadPage: %w", err)
	}
	p.Title = fi.Name()
	p.LastChanged = fi.ModTime()
	b, err := os.ReadFile(fpath)
	if err != nil {
		return p, fmt.Errorf("loadPage.ReadFile: %w", err)
	}
	p.Content = template.HTML(blackfriday.Run(b))
	return p, nil
}

func loadPages(src string) (Pages, error) {
	var ps Pages
	fs, err := os.ReadDir(src)
	if err != nil {
		return ps, fmt.Errorf("loadPages.ReadDir: %w", err)
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		fpath := filepath.Join(src, f.Name())
		p, err := loadPage(fpath)
		if err != nil {
			return ps, fmt.Errorf("loadPages.loadPage: %w", err)
		}
		ps = append(ps, p)
	}
	return ps, nil
}
