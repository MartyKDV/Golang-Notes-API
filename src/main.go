package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	notes := newNotesCollection()

	fileServer := http.FileServer(http.Dir("../static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/addNote", notes.handleAddNote)
	http.HandleFunc("/viewNote", notes.handleViewNote)
	http.HandleFunc("/editNote", notes.handleEditNote)
	http.HandleFunc("/deleteNote", notes.handleDeleteNote)

	printNotes(notes)
	log.Println("Server Has Successfully Started at Port :8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}

}

func printNotes(n *notesCollection) {
	for key, value := range n.notes {
		fmt.Println(key, value)
	}
}

type notesCollection struct {
	notes map[string]string
}

func newNotesCollection() *notesCollection {

	return &notesCollection{
		notes: map[string]string{
			"test": "testing message",
		},
	}
}

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
