package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// HTTPFS represents a httpfs handler
type HTTPFSHandler struct {
	path   string
	err    error
	isFile bool
}

// UNIX calls

// CatFile reproduces `cat $file_path` unix operation
func (h *HTTPFSHandler) CatFile(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(h.path)
	if err != nil {
		h.err = err
		h.HandleNotFound(w, r)
		return
	}

	log.Printf("path=%s,method=%s,unix_op=cat %s,response_code=%d", h.path, r.Method, h.path, http.StatusOK)
	_, _ = w.Write(b)
}

// LsDir reproduces `ls $dir_path` unix operation
func (h *HTTPFSHandler) LsDir(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(h.path)
	if err != nil {
		h.err = err
		h.HandleNotFound(w, r)
		return
	}

	for _, file := range files {
		line := fmt.Sprintf("%s\n", file.Name())
		w.Write([]byte(line))
	}
}

// TouchFile reproduces `touch $file_path` unix operation
func (h *HTTPFSHandler) TouchFile(w http.ResponseWriter, r *http.Request) {
	_, err := os.OpenFile(h.path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		h.err = err
		h.HandleConflict(w, r)
		return
	}

	log.Printf("path=%s,method=%s,unix_op=touch %s,response_code=%d", h.path, r.Method, h.path, http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}

// Handlers

// HandleNotFound handles NotFound scenarios
func (h *HTTPFSHandler) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("path=%s,method=%s,response_code=%d,err=%s", h.path, r.Method, http.StatusNotFound, h.err)
	http.NotFound(w, r)
}

// HandleConflict handles Conflict scenarios
func (h *HTTPFSHandler) HandleConflict(w http.ResponseWriter, r *http.Request) {
	log.Printf("path=%s,method=%s,response_code=%d,err=%s", h.path, r.Method, http.StatusConflict, h.err)
	w.WriteHeader(http.StatusConflict)
	io.WriteString(w, "conflict\n")
}

// HandleBadRequest handles BadRequest scenarios
func (h *HTTPFSHandler) HandleBadRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("path=%s,method=%s,response_code=%d,err=%s", h.path, r.Method, http.StatusBadRequest, h.err)
	w.WriteHeader(http.StatusConflict)
	io.WriteString(w, "bad request\n")
}

// HandleGet handles GET requests
func (h *HTTPFSHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if h.isFile {
		h.CatFile(w, r)
		return
	}

	h.LsDir(w, r)
}

// HandlePost handles POST requests
func (h *HTTPFSHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	// if file exists we exit
	if h.err == nil {
		h.err = errors.New("file/dir exists")
		h.HandleConflict(w, r)
		return

	}

	// parse form looking for type
	r.ParseForm()
	t := r.PostFormValue("type")
	if t == "" {
		h.err = errors.New("type missing")
		h.HandleBadRequest(w, r)
		return
	}

	if t == "file" {
		h.TouchFile(w, r)
		return
	}

	if t == "dir" {
		// TODO
		log.Printf("path=%s,method=%s,response_code=%d,err=%s", h.path, r.Method, http.StatusNotImplemented, "not implemented")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	h.err = errors.New("wrong type")
	h.HandleBadRequest(w, r)
}

// HandlePut handles PUT requests
func (h *HTTPFSHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// HandleDelete handles DELETE requests
func (h *HTTPFSHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// ServeHTTP satisfies http.Handler interface
func (h *HTTPFSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// build file/dir path
	h.path = fmt.Sprintf("%s%s", UserHomeDir(), r.URL)
	h.err = nil

	// find if path is file or directory
	isFile, err := IsFile(h.path)
	if err != nil {
		h.err = err
	}
	h.isFile = isFile

	switch {
	case r.Method == "GET":
		// cat and ls
		h.HandleGet(w, r)
		return
	case r.Method == "POST":
		// touch and mkdir
		h.HandlePost(w, r)
		return
	case r.Method == "PUT":
		// echo
		h.HandlePut(w, r)
		return
	case r.Method == "DELETE":
		// rm
		h.HandleDelete(w, r)
		return
	}
}

func main() {
	http.Handle("/", &HTTPFSHandler{})
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
