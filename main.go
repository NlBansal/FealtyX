package main

import (
	"fealtyx-student-api/handlers" // Adjust if your module name or folder structure is different
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Routes to handlers for CRUD operations on students
	r.HandleFunc("/students", handlers.CreateStudent).Methods("POST")
	r.HandleFunc("/students", handlers.GetAllStudents).Methods("GET")
	r.HandleFunc("/students/{id}", handlers.GetStudentByID).Methods("GET")
	r.HandleFunc("/students/{id}", handlers.UpdateStudentByID).Methods("PUT")
	r.HandleFunc("/students/{id}", handlers.DeleteStudentByID).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", handlers.GetStudentSummary).Methods("GET")

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
