package handlers

import (
	"bytes"
	"encoding/json"
	"fealtyx-student-api/models"
	"fealtyx-student-api/storage"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateStudent handler
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storage.Mutex.Lock()
	student.ID = storage.IDCounter
	storage.Students[student.ID] = student
	storage.IDCounter++
	storage.Mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// GetAllStudents handler
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	students := make([]models.Student, 0, len(storage.Students))
	for _, student := range storage.Students {
		students = append(students, student)
	}

	json.NewEncoder(w).Encode(students)
}

// GetStudentByID handler
func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	storage.Mutex.Lock()
	student, exists := storage.Students[id]
	storage.Mutex.Unlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(student)
}

// UpdateStudentByID handler
func UpdateStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storage.Mutex.Lock()
	if _, exists := storage.Students[id]; !exists {
		storage.Mutex.Unlock()
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	updatedStudent.ID = id
	storage.Students[id] = updatedStudent
	storage.Mutex.Unlock()

	json.NewEncoder(w).Encode(updatedStudent)
}

// DeleteStudentByID handler
func DeleteStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	storage.Mutex.Lock()
	if _, exists := storage.Students[id]; !exists {
		storage.Mutex.Unlock()
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	delete(storage.Students, id)
	storage.Mutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// GetStudentSummary handler using Ollama API
func GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	storage.Mutex.Lock()
	student, exists := storage.Students[id]
	storage.Mutex.Unlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	prompt := fmt.Sprintf("Generate a summary for a student named %s, who is %d years old with the email %s.", student.Name, student.Age, student.Email)

	response, err := http.Post("http://localhost:11434/v1/generate", "application/json", bytes.NewBuffer([]byte(`{"prompt": "`+prompt+`"}`)))
	if err != nil || response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to connect to Ollama", http.StatusInternalServerError)
		return
	}

	var summary map[string]string
	if err := json.NewDecoder(response.Body).Decode(&summary); err != nil {
		http.Error(w, "Failed to parse Ollama response", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"summary": summary["text"]})
}
