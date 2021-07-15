package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	// Path to folder with notes
	folderPath := "notes"
	notes := getNotes(folderPath)

	// File server that contains the html files
	fileServer := http.FileServer(http.Dir("../static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/addNote", notes.handleAddNote)
	http.HandleFunc("/viewNote", notes.handleViewNote)
	http.HandleFunc("/editNote", notes.handleEditNote)
	http.HandleFunc("/deleteNote", notes.handleDeleteNote)

	// Debug to see all notes
	printNotes(notes)

	// Start server
	log.Println("Server Has Successfully Started at Port :8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}

}

// Temp func to print all notes on console - name of note, data of note
func printNotes(n *notesCollection) {
	for key, value := range n.notes {
		fmt.Println(key, value)
	}
}

type notesCollection struct {
	notes map[string]string
}

// Retrieve all the notes from the folder and return as a notesCollection pointer
func getNotes(folderPath string) *notesCollection {

	var notesFromFiles *notesCollection = &notesCollection{}
	notesFromFiles.notes = make(map[string]string)

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		name := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]

		fileName := folderPath + "/" + file.Name()
		data, err := ioutil.ReadFile(fileName)

		if err != nil {
			panic(err)
		}

		notesFromFiles.notes[name] = string(data)
	}

	return notesFromFiles
}

// Request Handle Functrions
func (h *notesCollection) handleAddNote(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	data := r.FormValue("data")

	h.notes[name] = data
}

func (h *notesCollection) handleViewNote(w http.ResponseWriter, r *http.Request) {

}

func (h *notesCollection) handleEditNote(w http.ResponseWriter, r *http.Request) {

}

func (h *notesCollection) handleDeleteNote(w http.ResponseWriter, r *http.Request) {

}
