package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Note struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	NoteText string `json:"note_text"`
}

var NoteStorage = []Note{}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Hi %s", r.URL.Query().Get("name"))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func saveNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Sorry, only POST method is supported.", 400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Bad request.", 400)
		return
	}

	newNote := Note{}
	if json.Unmarshal(body, &newNote) != nil {
		http.Error(w, "Bad request.", 400)
		return
	}
	NoteStorage = append(NoteStorage, newNote)
	fmt.Printf("Введённые данные: \n  имя: %s\n  фамилия: %s\n  заметка: %s\n", newNote.Name, newNote.Surname, newNote.NoteText)
}

func listAllNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Sorry, only GET method is supported.", 400)
		return
	}

	jsonResp, err := json.Marshal(NoteStorage)
	if err != nil {
		http.Error(w, "error happened", 500)
	}

	w.Write(jsonResp)
	return
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/save_note", saveNote)
	http.HandleFunc("/list_all", listAllNotes)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
