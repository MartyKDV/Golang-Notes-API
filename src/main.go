package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	jsonPath := "../notes/notes.json"
	notes := getNotes(jsonPath)

	fileServer := http.FileServer(http.Dir("../static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/addNote", notes.handleAddNote)
	http.HandleFunc("/viewNotes", notes.handleViewNote)
	http.HandleFunc("/editNote", notes.handleEditNote)
	http.HandleFunc("/deleteNote", notes.handleDeleteNote)

	log.Println("Server Has Successfully Started at Port :8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}

}

type Note struct {
	Data           string `json:"data"`
	Author         string `json:"author"`
	TimeCreated    string `json:"time_created"`
	TimeLastEdited string `json:"time_last_edited"`
}

type notesCollection struct {
	notes    []Note
	fileName string
}

func getNotes(path string) *notesCollection {

	var notesFromFiles *notesCollection = &notesCollection{}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			panic(err)
		}
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	notesFromFiles.notes = make([]Note, 0)
	notesFromFiles.fileName = path

	json.Unmarshal(file, &notesFromFiles.notes)

	for note := range notesFromFiles.notes {

		log.Println(notesFromFiles.notes[note])
	}

	return notesFromFiles
}

func (h *notesCollection) handleAddNote(w http.ResponseWriter, r *http.Request) {

	author := r.FormValue("author")
	data := r.FormValue("data")
	currentTime := time.Now()

	lastIndex := len(h.notes)
	h.notes = append(h.notes, Note{Data: data, Author: author, TimeCreated: currentTime.Format("2006-01-02 3:4:5 PM"), TimeLastEdited: currentTime.Format("2006-01-02 3:4:5 PM")})
	log.Println("Note added: ", h.notes[lastIndex])

	notesBytes, err := json.MarshalIndent(h.notes, "", " ")
	if err != nil {
		panic(err)
	}

	errWriteJSON := ioutil.WriteFile(h.fileName, notesBytes, 0644)
	if errWriteJSON != nil {
		panic(errWriteJSON)
	}

	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)

}

func (h *notesCollection) handleViewNote(w http.ResponseWriter, r *http.Request) {

	for i := range h.notes {

		fmt.Fprintf(w, "--------------------------------------------------------------------------------\n")
		fmt.Fprintf(w, "Author: "+h.notes[i].Author+"\n")
		fmt.Fprintf(w, h.notes[i].Data+"\n")
		fmt.Fprintf(w, "Time Created: "+h.notes[i].TimeCreated+"\nLast Edited:  "+h.notes[i].TimeLastEdited+"\n")
	}
	fmt.Fprintf(w, "--------------------------------------------------------------------------------\n")

}

func (h *notesCollection) handleEditNote(w http.ResponseWriter, r *http.Request) {

}

func (h *notesCollection) handleDeleteNote(w http.ResponseWriter, r *http.Request) {

}
