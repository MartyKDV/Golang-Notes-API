package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	jsonPath := "../notes/notes.json"
	notes := getNotes(jsonPath)

	router := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("../static"))

	// Request Handlers
	router.Handle("/", fileServer)
	router.HandleFunc("/new", notes.handleAddNote)
	router.HandleFunc("/notes", notes.handleViewNote)
	router.HandleFunc("/details/{id}", notes.handleEditNote)
	router.HandleFunc("/editNoteData/{id}", notes.handleEditNoteData)
	router.HandleFunc("/deleteNote/{id}", notes.handleDeleteNote)

	log.Println("Server Has Successfully Started at Port :8080...")
	err := http.ListenAndServe(":8080", router)
	checkError(err)

}

type Note struct {
	Data           string `json:"data"`
	Author         string `json:"author"`
	TimeCreated    string `json:"time_created"`
	TimeLastEdited string `json:"time_last_edited"`
}

type notesCollection struct {
	notes    map[int]*Note
	fileName string
}

type noteInfo struct {
	Note int
	Data *Note
}

// Populate the in-memory map
func getNotes(path string) *notesCollection {

	var notesFromFiles *notesCollection = &notesCollection{}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := os.Create(path)
		checkError(err)
	}

	file, err := ioutil.ReadFile(path)
	checkError(err)

	notesFromFiles.notes = make(map[int]*Note)
	notesFromFiles.fileName = path

	json.Unmarshal(file, &notesFromFiles.notes)

	for note := range notesFromFiles.notes {

		log.Println(notesFromFiles.notes[note])
	}

	return notesFromFiles
}

func (h *notesCollection) handleAddNote(w http.ResponseWriter, r *http.Request) {

	// Return the form if the request is get and add the element if the request is post
	switch r.Method {
	case "POST":
		{
			author := r.FormValue("author")
			data := r.FormValue("data")

			currentTime := time.Now()
			lastIndex := len(h.notes)
			var added bool

			for i := 0; i < lastIndex; i++ {

				if _, exists := h.notes[i]; !exists {

					h.notes[i] = &Note{
						Author:         author,
						Data:           data,
						TimeCreated:    currentTime.Format("2006-01-02 3:4:5 PM"),
						TimeLastEdited: currentTime.Format("2006-01-02 3:4:5 PM"),
					}
					added = true
					break
				}
			}

			if !added {
				h.notes[lastIndex] = &Note{
					Author:         author,
					Data:           data,
					TimeCreated:    currentTime.Format("2006-01-02 3:4:5 PM"),
					TimeLastEdited: currentTime.Format("2006-01-02 3:4:5 PM"),
				}
			}

			log.Println("Note added: ", *h.notes[lastIndex])

			notesBytes, err := json.MarshalIndent(h.notes, "", " ")
			checkError(err)

			err = ioutil.WriteFile(h.fileName, notesBytes, 0644)
			checkError(err)

			http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
		}

	case "GET":
		{
			templ, err := template.ParseFiles("../static/addNote.html")
			checkError(err)

			err = templ.Execute(w, nil)
			checkError(err)
		}
	}
}

// Return the html to view all notes
func (h *notesCollection) handleViewNote(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../static/viewNotes.html")
	checkError(err)

	err = templ.Execute(w, h.notes)
	checkError(err)
}

// Return the html with the edit note form based on the id
func (h *notesCollection) handleEditNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	note, err := strconv.Atoi(vars["id"])
	checkError(err)

	var n noteInfo = noteInfo{Note: note, Data: h.notes[note]}

	templ, err := template.ParseFiles("../static/editNote.html")
	checkError(err)

	err = templ.Execute(w, n)
	checkError(err)

}

// Temporary function - will implement the algorithm
// to handleEditNote - switch on the type of request (PUT/GET)
func (h *notesCollection) handleEditNoteData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	note, err := strconv.Atoi(vars["id"])
	checkError(err)

	author := r.FormValue("author")
	data := r.FormValue("data")
	currentTime := time.Now()

	// Edit the in-memory data structure entry
	h.notes[note].Author = author
	h.notes[note].Data = data
	h.notes[note].TimeLastEdited = currentTime.Format("2006-01-02 3:4:5 PM")

	dat, err := json.MarshalIndent(h.notes, "", " ")
	checkError(err)
	err = ioutil.WriteFile(h.fileName, dat, 0644)
	checkError(err)

	http.Redirect(w, r, "http://localhost:8080/notes", http.StatusFound)
}

// Delete the note from in-memory and the JSON file
func (h *notesCollection) handleDeleteNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	note, err := strconv.Atoi(vars["id"])
	checkError(err)

	delete(h.notes, note)

	dat, err := json.MarshalIndent(h.notes, "", " ")
	checkError(err)

	err = ioutil.WriteFile(h.fileName, dat, 0644)
	checkError(err)

	http.Redirect(w, r, "http://localhost:8080/notes", http.StatusFound)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
