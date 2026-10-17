package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Done       TaskStatus = "done"
	fileName   string     = "task-tracker.json"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task

	err = json.Unmarshal(data, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func getNextID(tasks []Task) int {
	maxID := 0

	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1
}

func addTask(description string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Gagal membaca file: %v\n", err)
		return
	}

	now := time.Now()
	newTask := Task{
		ID:          getNextID(tasks),
		Description: description,
		Status:      Todo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Gagal menyimpan task: %v\n", err)
		return
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
}

func listTasks(filter string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Gagal membaca file: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("Belum ada task yang tersimpan.")
		return
	}

	found := false
	for _, task := range tasks {
		if filter == "" || string(task.Status) == filter {
			found = true
			fmt.Printf("[%d] %s (Status: %s) - Dibuat: %s\n",
				task.ID,
				task.Description,
				task.Status,
				task.CreatedAt.Format("2006-01-02 15:04"),
			)
		}
	}

	if !found {
		fmt.Printf("Tidak ada task dengan status '%s'.\n", filter)
	}
}

func updateTask(id int, newDescription string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Gagal membaca file: %v\n", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Task dengan ID %d tidak ditemukan.\n", id)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Gagal menyimpan perubahan: %v\n", err)
		return
	}

	fmt.Println("Task updated successfully")
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Gagal membaca file: %v\n", err)
		return
	}

	foundIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		fmt.Printf("Error: Task dengan ID %d tidak ditemukan.\n", id)
		return
	}

	/*
		// Buat slice baru untuk menampung task yang tidak dihapus
		var newTasks []Task

		for i, task := range tasks {
			if i != foundIndex {
				newTasks = append(newTasks, task)
			}
		}
	*/
	tasks = append(tasks[:foundIndex], tasks[foundIndex+1:]...)

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Gagal menyimpan perubahan: %v\n", err)
		return
	}

	fmt.Println("Task deleted successfully")
}

func markStatus(id int, status TaskStatus) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Gagal membaca file: %v\n", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Task dengan ID %d tidak ditemukan.\n", id)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Gagal mengubah status: %v\n", err)
		return
	}

	fmt.Printf("Task marked as %s\n", status)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gunakan perintah: task-cli <command> [arguments]")
		fmt.Println("Contoh command: add, update, delete, mark-in-progress, mark-done, list")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Masukkan deskripsi task. Contoh: task-cli add \"Beli susu\"")
			return
		}
		addTask(os.Args[2])

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Error: Format salah. Contoh: task-cli update <id> \"Deskripsi baru\"")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ID harus berupa angka.")
			return
		}
		updateTask(id, os.Args[3])

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Format salah. Contoh: task-cli delete <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ID harus berupa angka.")
			return
		}
		deleteTask(id)

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Error: Format salah. Contoh: task-cli mark-in-progress <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ID harus berupa angka.")
			return
		}
		markStatus(id, InProgress)

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Error: Format salah. Contoh: task-cli mark-done <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ID harus berupa angka.")
			return
		}
		markStatus(id, Done)

	case "list":
		filter := ""
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		listTasks(filter)

	default:
		fmt.Printf("Perintah '%s' tidak dikenali.\n", command)
	}
}
