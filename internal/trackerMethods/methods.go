package trackerMethods

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"taskTrackerEasy/internal/task"
)

func CorrectStatus(status string) (string, error) {
	if strings.ToLower(status) == "todo" {
		return task.ToDo, nil
	} else if strings.ToLower(status) == "in-progress" {
		return task.InProgress, nil
	} else if strings.ToLower(status) == "done" {
		return task.Done, nil
	} else {
		return "", fmt.Errorf("There is no such status\ntodo\tin-progress\tdone")
	}
}

func AddTask(description, status string) (*task.Task, error) {
	file, err := os.OpenFile("Tracker.jsonl", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("File opening error: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var id int

	if fileInfo.Size() == 0 {
		id = 0
	} else {
		id, err = getLastId()
		fmt.Println(err)
		if err != nil && err.Error() != "jsonl unmarshal error" {
			return nil, err
		}
		id += 1
	}

	var task *task.Task = task.NewTask(id, description, status)

	jsonlTask, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("jsonl serialization error: %w", err)
	}

	_, err = file.Write([]byte("\n"))
	_, err = file.Write(jsonlTask)

	if err != nil {
		return nil, fmt.Errorf("Error writing to a file: %w", err)
	}

	return task, nil
}

func getLastId() (int, error) {
	file, err := os.Open("Tracker.jsonl")
	if err != nil {
		return 0, fmt.Errorf("File opening error: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("File info error: %w", err)
	}

	sizeFile := fileInfo.Size()
	var buferSize int64 = 100
	var lastLine []byte
	offset := sizeFile

	for {
		readSize := buferSize
		if buferSize > offset {
			readSize = offset
		}

		offset -= readSize

		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			return 0, fmt.Errorf("Shift error: %w", err)
		}

		buffer := make([]byte, readSize)

		_, err = file.Read(buffer)
		if err != nil {
			return 0, fmt.Errorf("Error reading file: %w", err)
		}

		for i := len(buffer) - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				lastLine = append(buffer[i+1:], lastLine...)

				var t task.Task
				json.Unmarshal(lastLine, &t)
				return t.Id, nil
			}
		}
		lastLine = append(buffer, lastLine...)

		if offset == 0 {
			break
		}

	}
	fmt.Println(string(lastLine))

	var t task.Task
	if err := json.Unmarshal(lastLine, &t); err != nil {
		return 0, fmt.Errorf("jsonl unmarshal error: %w", err)
	}
	return t.Id, nil
}

func GetAllLines() ([]task.Task, error) {
	file, err := os.Open("Tracker.jsonl")
	if err != nil {
		return nil, fmt.Errorf("Error when opening file: %w", err)
	}
	defer file.Close()

	currentLine := 0
	lineScanner := bufio.NewScanner(file)
	var t1 task.Task
	var t []task.Task

	for lineScanner.Scan() {

		if err := lineScanner.Err(); err != nil {
			return nil, fmt.Errorf("Error while reading file: %s", err)
		}

		if err = json.Unmarshal([]byte(lineScanner.Text()), &t1); err != nil {
			return nil, fmt.Errorf("jsonl unmarshal error: %w", err)

		}
		t = append(t, t1)

		currentLine++
	}

	return t, nil
}

func UpdateTask(id int, description, status string) error {
	file, err := os.OpenFile("Tracker.jsonl", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("File opening error: %w", err)
	}
	defer file.Close()

	lines, err := GetAllLines()
	if err != nil {
		return err
	}

	if description != "" {
		lines[id].Description = description
	}
	if status != "" {
		lines[id].Status = status
	}

	for _, line := range lines {
		jsonlLines, err := json.Marshal(line)
		if err != nil {
			return fmt.Errorf("jsonl serialization error: %w", err)
		}
		file.Write(jsonlLines)
		file.Write([]byte("\n"))
	}
	return nil

}

func DeleneTask(id int) error {
	if id < 0 {
		return fmt.Errorf("Uncorrect input")
	}

	tasks, err := GetAllLines()
	if err != nil {
		return err
	}

	var sliceLines []task.Task

	for _, task := range tasks {
		if task.Id != id {
			sliceLines = append(sliceLines, task)
		}
	}

	file, err := os.Create("Tracker.jsonl")
	if err != nil {
		return fmt.Errorf("File opening error: %w", err)
	}
	defer file.Close()

	for i, task := range sliceLines {
		jsonData, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("JSON marshal error: %w", err)
		}

		if _, err := file.Write(jsonData); err != nil {
			return fmt.Errorf("Write error: %w", err)
		}

		if i != len(sliceLines)-1 {
			if _, err := file.Write([]byte{'\n'}); err != nil {
				return fmt.Errorf("Write newline error: %w", err)
			}
		}
	}

	return nil
}
