package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client        *mongo.Client
	mongoURL      string
	ctx           context.Context
	notesDatabase *mongo.Database
	collection    *mongo.Collection
	err           error
)

func main() {

	mongoURL = "mongodb://mongodb:27017"
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	checkError(err)

	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	err = client.Connect(ctx)
	checkError(err)
	defer client.Disconnect(ctx)

	notesDatabase = client.Database("go-notes")
	collection = notesDatabase.Collection("notes")

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	fmt.Println("connected to nosql database:", mongoURL)

	router := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("../static"))

	// Request Handlers
	router.Handle("/", fileServer)
	router.HandleFunc("/new", handleAddNote)
	router.HandleFunc("/notes", handleViewNote)
	router.HandleFunc("/details/{id}", handleEditNote)
	router.HandleFunc("/editNote/{id}", handleEditNoteData)
	router.HandleFunc("/deleteNote/{id}", handleDeleteNote)

	log.Println("Server Has Successfully Started at Port :8080...")
	err = http.ListenAndServe(":8080", router)
	checkError(err)

}

type Note struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Data           string             `json:"data,omitempty" bson:"data,omitempty"`
	Author         string             `json:"author,omitempty" bson:"author,omitempty"`
	TimeCreated    string             `json:"time_created,omitempty" bson:"time_created,omitempty"`
	TimeLastEdited string             `json:"time_last_edited,omitempty" bson:"time_last_edited,omitempty"`
}

type noteInfo struct {
	Note int
	Data *Note
}

func handleAddNote(w http.ResponseWriter, r *http.Request) {

	// Return the form if the request is get and add the element if the request is post
	switch r.Method {
	case "POST":
		{
			author := r.FormValue("author")
			data := r.FormValue("data")
			currentTime := time.Now().Format("2006-01-02 3:4:5 PM")
			var note Note = Note{
				Author:         author,
				Data:           data,
				TimeCreated:    currentTime,
				TimeLastEdited: currentTime,
			}
			ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
			_, err := collection.InsertOne(ctx, note)
			checkError(err)
			var notes []Note
			ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
			cursor, err := collection.Find(ctx, bson.M{})
			checkError(err)
			defer cursor.Close(ctx)

			for cursor.Next(ctx) {
				var n Note
				cursor.Decode(&n)
				notes = append(notes, n)
			}

			log.Println(notes)

			http.Redirect(w, r, "/", http.StatusFound)
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
func handleViewNote(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../static/viewNotes.html")
	checkError(err)

	var notes []Note
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	checkError(err)

	cursor, err := collection.Find(ctx, bson.M{})
	checkError(err)

	for cursor.Next(ctx) {
		var n Note
		cursor.Decode(&n)
		notes = append(notes, n)
		log.Println(n)
	}

	err = templ.Execute(w, notes)
	checkError(err)
}

// Return the html with the edit note form based on the id
func handleEditNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	checkError(err)

	var n Note
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor := collection.FindOne(ctx, bson.M{"_id": idPrimitive})
	cursor.Decode(&n)

	templ, err := template.ParseFiles("../static/editNote.html")
	checkError(err)

	err = templ.Execute(w, n)
	checkError(err)

}

// Temporary function - will implement the algorithm
// to handleEditNote - switch on the type of request (PUT/GET)
func handleEditNoteData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	checkError(err)

	author := r.FormValue("author")
	data := r.FormValue("data")
	currentTime := time.Now().Format("2006-01-02 3:4:5 PM")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idPrimitive},
		bson.M{
			"$set": bson.M{"author": author, "data": data, "time_last_edited": currentTime},
		},
	)
	checkError(err)

	http.Redirect(w, r, "/notes", http.StatusFound)
}

// Delete the note from in-memory and the JSON file
func handleDeleteNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	checkError(err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = collection.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	checkError(err)

	http.Redirect(w, r, "/notes", http.StatusFound)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
