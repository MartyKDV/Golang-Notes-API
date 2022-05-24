Golang Notes API
==============
https://github.com/MartyKDV/Golang-Notes-API  

Overview
------------
This is a project to create an API that serves to handle notes a person might take. The project has been inspired by my desire to learn the Go language and start creating useful applications with it as a part of my SAP summer practice program. The current functions of the project are:

* Create notes
* View Notes
* Edit Notes
* Delete Notes

The project is written in Go and the HTML pages are templates which Go has configured so data can be passed from the server. This has been realized through the use of the html/template package. The requests are routed through the use of the `gorilla/mux` package, which creates a router that is used to route and send the request to be handles by the handler functions. The notes are stored as a map in-memory and as objects in a JSON file located in the `notes` folder. The project is run by building the `main.go` file and afterwards running it to start the server. Currently everything works on the localserver at port 8080.  
The current endpoints are as follows:
* /new - to add a note
* /notes - to view all notes and possibly edit/delete any

Example
=======

Docker run image
----------------

**When running the image - port forward a port to 8080 for the container (with -p \<host port\>:8080)**  

\----------------------------------------------------------------------------------------------------------------------

**URL: localhost:8080 - Homepage. Navigate to add a note or view notes**

![Homepage](https://i.ibb.co/6P78KdX/homepage.png)
 
 **URL: localhost:8080/notes - View all notes. Buttons to edit or delete a specific note**
 
 ![View Notes](https://i.ibb.co/JCVhBYP/notes.png)

**URL: localhost:8080/new - Create a new note**

![Create Note](https://i.ibb.co/yRw6sZD/new-note.png)

Credits
======
For more information into the `gorilla/mux` package (the router): 
[https://www.gorillatoolkit.org/pkg/mux](https://www.gorillatoolkit.org/pkg/mux) 
