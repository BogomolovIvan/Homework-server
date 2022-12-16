package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Note struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	NoteText string `json:"note_text"`
}

// var NoteStorage = []Note{}
var reader = bufio.NewReader(os.Stdin)
var httpClient = &http.Client{}

func main() {
	CreateNewNote()
	for {
		fmt.Println("что делать дальше? \n c - создать новую заметку \t l - вывести все записи \t q - завершить исполнение")
		Command := ReadLine("Введите команду: ")
		switch Command {
		case "c":
			CreateNewNote()
			break
		case "l":
			ListAllNotes()
			break
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Команда не найдена. Попробуйте ещё раз")
		}
	}
}

func ReadLine(helpText string) string {
	fmt.Print(helpText)
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSuffix(input, "\n")
}

func CreateNewNote() {
	NewNote := Note{}
	NewNote.Name = ReadLine("Введите имя: ")
	NewNote.Surname = ReadLine("Введите фамилию: ")
	NewNote.NoteText = ReadLine("Введите заметку: ")

	NewNote.SaveNote()
}

func (newNote *Note) SaveNote() {
	json, err := json.Marshal(newNote)
	if err != nil {
		log.Fatal(err)
	}

	bb := bytes.Buffer{}
	bb.Write(json)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/save_note", &bb)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("Введённые данные: \n  имя: %s\n  фамилия: %s\n  заметка: %s\n", newNote.Name, newNote.Surname, newNote.NoteText)
	} else {
		fmt.Printf("Error occured:\nStatus:%s", resp.Status)
	}
}

func ListAllNotes() {
	resp, err := http.Get("http://127.0.0.1:8000/list_all")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var noteList []Note
	if json.Unmarshal(body, &noteList) != nil {
		log.Fatal(err)
		return
	}
	for index, note := range noteList {
		fmt.Printf("\nЗаписка №%d\n  имя: %s\n  фамилия: %s\n  заметка: %s\n", index+1, note.Name, note.Surname, note.NoteText)
	}
}
