package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	http.HandleFunc("/viewNotes", notes.handleViewNote)
	http.HandleFunc("/editNote", notes.handleEditNote)
	http.HandleFunc("/deleteNote", notes.handleDeleteNote)

	// Debug to see all notes
	printNotes(notes)

	// Start server
	log.Println("Server started listening at port :8080...")
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

// Request Handle Functions
func (h *notesCollection) handleAddNote(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	data := r.FormValue("data")

	h.notes[name] = data

	// Create files - placeholder for DB operations
	fileName := "notes/" + name + ".txt"

	file, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	io.WriteString(file, data)
	log.Println("File Created...")

	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)

}

func (h *notesCollection) handleViewNote(w http.ResponseWriter, r *http.Request) {

	for key, value := range h.notes {
		fmt.Fprintf(w, key+" - "+value+"\n")
	}

}

func (h *notesCollection) handleEditNote(w http.ResponseWriter, r *http.Request) {

	nameOfNote := r.FormValue("nameOfNote")
	name := r.FormValue("name")
	data := r.FormValue("note")

	path := "notes/" + nameOfNote + ".txt"

	if data != "" {

		file, err := os.OpenFile(path, os.O_WRONLY, 0644)

		if err != nil {
			log.Println(err)
			return
		}

		os.Truncate(path, 0)
		io.WriteString(file, data)
		h.notes[nameOfNote] = data
		file.Sync()

		file.Close()
	}

	if name != "" {
		err := os.Rename(path, "notes/"+name+".txt")
		if err != nil {
			panic(err)
		}
		h.notes[name] = h.notes[nameOfNote]
		delete(h.notes, nameOfNote)
	}

	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)

}

func (h *notesCollection) handleDeleteNote(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	path := "notes/" + name + ".txt"
	os.Remove(path)

	delete(h.notes, name)

	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
}
