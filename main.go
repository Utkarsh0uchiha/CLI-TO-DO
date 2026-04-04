package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          int
	Name        string
	Status      string
	CreatedAt   string
	CompletedAt string
}

func getNextID() int {
	file, err := os.Open("todo.csv")

	if err != nil {
		return 1
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 1
	}
	n := len(records)

	if n < 2 {
		return 1
	}

	lastRow := records[n-1]

	lastID, err := strconv.Atoi(lastRow[0])
	if err != nil {
		panic(err)

	}

	return lastID + 1

}

func saveTask(task Task) {
	// string to write
	header := []string{"ID", "Name", "Status", "CreatedAt", "CompletedAt"}
	s := []string{strconv.Itoa(task.Id), task.Name, task.Status, task.CreatedAt, task.CompletedAt}
	// opening file
	file, err := os.OpenFile("todo.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// checking error in opening file
	if err != nil {
		fmt.Println("Error! file don't exist")
		panic(err)
	}

	// closing file
	defer file.Close()

	// creating a writer
	writer := csv.NewWriter(file)

	// getting file info
	info, err := file.Stat()
	// checking error in getting info
	if err != nil {
		fmt.Println("Error getting file info: ", err)
		return
	}

	// if file size is empty insert into header in file
	if info.Size() == 0 {
		err := writer.Write(header)
		if err != nil {
			panic(err)
		}
	}

	// writing data in file
	err = writer.Write(s)

	if err != nil {
		panic(err)
	}
	writer.Flush()
}

func getAllTasks() []Task {
	// opening file
	file, err := os.Open("todo.csv")
	// handling file opening error
	if err != nil {
		fmt.Println("Error! file don't exist")
		panic(err)
	}

	// defer file closing
	defer file.Close()
	// creating a reader
	reader := csv.NewReader(file)
	//reading the file
	records, err := reader.ReadAll()
	// handling file read error
	if err != nil {
		fmt.Println("Error! Can't read file")
		panic(err)
	}
	// creating a struct to convert from CSV
	n := len(records)
	List := []Task{}
	if n <= 1 {
		return List
	}

	for i := 1; i < n; i++ {
		row := records[i]
		if len(row) < 5 {
			continue
		}
		// converting string to Int
		ID, err := strconv.Atoi(row[0])
		if err != nil {
			panic(err)
		}

		task := Task{
			Id:          ID,
			Name:        row[1],
			Status:      row[2],
			CreatedAt:   row[3],
			CompletedAt: row[4],
		}
		// append the struct in the slice
		List = append(List, task)

	}

	return List
}

func completeTask(id int, tm time.Time) bool {
	file, err := os.Open("todo.csv")
	if err != nil {
		fmt.Println("Error! file don't exist")
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error! Can't read file")
		panic(err)
	}

	found := false
	for i := 1; i < len(records); i++ {
		row := records[i]
		ID, err := strconv.Atoi(row[0])
		if err != nil {
			panic(err)
		}
		if ID == id {
			row[2] = "completed"
			row[4] = tm.Format("2006-01-02 15:04")
			found = true
		}
	}
	if found {
		editFile, err := os.Create("todo.csv")
		if err != nil {
			panic(err)
		}
		defer editFile.Close()

		writer := csv.NewWriter(editFile)
		err = writer.WriteAll(records)
		if err != nil {
			panic(err)
		}

		writer.Flush()
	}
	return found
}

func deleteTask(id int) error {
	file, err := os.Open("todo.csv")
	if err != nil {

		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		return err
	}

	var updateRecords [][]string

	updateRecords = append(updateRecords, records[0])
	for i := 1; i < len(records); i++ {
		row := records[i]
		ID, err := strconv.Atoi(row[0])
		if err != nil {
			return err
		}

		if ID != id {
			updateRecords = append(updateRecords, row)
		}
	}
	editFile, err := os.Create("todo.csv")

	if err != nil {
		return err
	}

	defer editFile.Close()

	writer := csv.NewWriter(editFile)
	defer writer.Flush()

	err = writer.WriteAll(updateRecords)

	if err != nil {
		return err
	}

	return nil

}
func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage:\n add 'task name' \n complete 'task id' \n delete 'task id' \n list \n")
	} else {
		command := args[1]
		tm := time.Now()

		switch command {
		case "add":
			if len(args) < 3 {
				fmt.Println("Please provide task name")
				return
			}
			work := args[2]

			task := Task{getNextID(), work, "pending", tm.Format("2006-01-02 15:04"), ""}
			saveTask(task)
			fmt.Printf("Adding task: %s\n", work)
		case "list":
			List := getAllTasks()
			fmt.Println("Listing tasks")

			fmt.Printf("%-5s %-12s %-10s %-20s %-20s\n", "ID", "Name", "Status", "CreatedAt", "CompletedAt")
			for _, row := range List {

				fmt.Printf("%-5d %-12s %-10s %-20s %-20s\n", row.Id, row.Name, row.Status, row.CreatedAt, row.CompletedAt)
			}
		case "complete":

			if len(args) < 3 {
				fmt.Println("Please provide the ID")
				return
			}

			id := args[2]
			intId, err := strconv.Atoi(id)
			if err != nil {
				panic(err)
			}
			if completeTask(intId, tm) {
				List := getAllTasks()
				for _, row := range List {
					if row.Id == intId {
						fmt.Printf("%-5d %-12s %-10s %-20s %-20s\n", row.Id, row.Name, row.Status, row.CreatedAt, row.CompletedAt)
						break
					}
				}
			} else {
				fmt.Println("Task not found")
			}
		case "delete":
			if len(args) < 3 {
				fmt.Println("Please Provide an ID")
				return
			}
			id := args[2]
			intId, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			err = deleteTask(intId)

			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("Record Deleted Successfully!")
			}
		default:
			fmt.Println("Unknown Command")
		}
	}

}
