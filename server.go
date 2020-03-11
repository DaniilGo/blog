package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	db        *sql.DB
	Title     string
	Templates map[templateName]*template.Template
}

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	server := Server{
		Title:     "Posts from Habr",
		Templates: createTemplates(),
		db:        db,
	}

	server.insertDefault()

	router := mux.NewRouter()
	router.HandleFunc("/", server.handlePostsList)
	router.HandleFunc("/post/{id:[0-9]+}", server.handleSinglePost)
	router.HandleFunc("/edit/{id:[0-9]+}", server.handleEditPost)
	router.HandleFunc("/results", server.handleResults)

	port := "8080"
	log.Printf("start server on port: %v", port)

	go func() {
		_ = http.ListenAndServe(":"+port, router)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	server.truncate()
}

func (server *Server) handlePostsList(wr http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate(server.Templates, List)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	posts, err := getPosts(server.db)
	if err != nil {
		log.Println(err)
		return
	}

	if err := tmpl.ExecuteTemplate(wr, "page", posts); err != nil {
		err = errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (server *Server) handleSinglePost(wr http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate(server.Templates, Single)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	vars := mux.Vars(req)

	id := vars["id"]
	if len(id) == 0 {
		log.Println(errors.New("empty id"))
		return
	}

	post, err := getPost(server.db, id)
	if err != nil {
		err := errors.Wrap(err, "empty post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if err := tmpl.ExecuteTemplate(wr, "page", post); err != nil {
		err := errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (server *Server) handleEditPost(wr http.ResponseWriter, req *http.Request) {
	tmpl := getTemplate(server.Templates, Edit)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	vars := mux.Vars(req)

	id := vars["id"]
	if len(id) == 0 {
		log.Println("edit: empty id")
		return
	}

	post, err := getPost(server.db, id)
	if err != nil {
		err := errors.Wrap(err, "empty post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if err := tmpl.ExecuteTemplate(wr, "page", post); err != nil {
		err = errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (server *Server) handleResults(wr http.ResponseWriter, req *http.Request) {
	idVal := req.FormValue("id")
	if len(idVal) == 0 {
		log.Print("results: empty id")
		return
	}

	id, err := strconv.Atoi(idVal)
	if err != nil {
		err := errors.Wrapf(err, "id from form value: %v", idVal)
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	post := Post{
		Id:      id,
		Title:   req.FormValue("title"),
		Date:    req.FormValue("date"),
		Link:    req.FormValue("link"),
		Comment: req.FormValue("comment"),
	}

	if err := editPost(server.db, post, idVal); err != nil {
		log.Print(err)
		return
	}

	http.Redirect(wr, req, "/", http.StatusFound)
}
