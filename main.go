package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var (
	tasks   = []Task{}
	nextID  = 1
	taskMux sync.Mutex
)

func main() {
	// Routes
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskByID)

	// Serve static frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Handle list + add
func handleTasks(w http.ResponseWriter, r *http.Request) {
	taskMux.Lock()
	defer taskMux.Unlock()

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(tasks)

	case "POST":
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t.ID = nextID
		nextID++
		tasks = append(tasks, t)
		json.NewEncoder(w).Encode(t)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle delete
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	taskMux.Lock()
	defer taskMux.Unlock()

	// extract ID
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "DELETE":
		for i, t := range tasks {
			if t.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "task not found", http.StatusNotFound)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
