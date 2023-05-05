package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/russross/blackfriday/v2"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	runServer()
}

func runServer() {
	log.Println("Starting server...")

	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Get("/blog", blogsHandler)
	r.Get("/blog/{slug}", blogHandler)
	r.Get("/cv", cvHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "cv", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func blogsHandler(w http.ResponseWriter, r *http.Request) {
	filesInDir, err := os.ReadDir("articles")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	posts := BlogPosts{
		Posts: []string{},
	}

	for _, file := range filesInDir {
		if file.IsDir() {
			continue
		}

		nameWithoutExt := strings.TrimSuffix(file.Name(), ".md")
		posts.Posts = append(posts.Posts, nameWithoutExt)
	}

	err = templates.ExecuteTemplate(w, "blogs", posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "slug")

	postContent, err := os.ReadFile("articles/" + param + ".md")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	htmlContent := blackfriday.Run(postContent)

	post := BlogPost{
		Content: template.HTML(string(htmlContent)),
	}

	err = templates.ExecuteTemplate(w, "blog", post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
}
