package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	folderPath := "../notes"
	notes := getNotes(folderPath)

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

func getNotes(path string) *notesCollection {

	var notesFromFiles *notesCollection

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}

		name := path[:len(path)-len(filepath.Ext(path))]

		data, inputErr := ioutil.ReadFile(path)

		if inputErr != nil {
			fmt.Println(inputErr)
			return nil
		}

		notesFromFiles.notes[name] = string(data)

		return nil
	})

	if err != nil {
		panic(err)
	}

	/*return &notesCollection{
		notes: map[string]string{
			"test": "testing message",
		},
	}*/

	return notesFromFiles
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
