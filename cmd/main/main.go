package main

import (
	"fmt"
	"os"
	"strconv"
	"taskTrackerEasy/internal/task"
	"taskTrackerEasy/internal/trackerMethods"
)

func main() {
	// fmt.Println(trackerMethods.AddTask("wq", "3"))

	// fmt.Println(trackerMethods.UpdateTask(2, "dadfsfsdffsfd", "fsdfsdfsfds", "dfsdsd"))

	// fmt.Println(trackerMethods.GetLine(buf))

	args := os.Args
	if len(args) < 1 {
		fmt.Println("Not enough parameters")
		return
	}

	switch args[1] {
	case "add":
		if len(args) == 3 {
			trackerMethods.AddTask(args[2], task.ToDo)
		} else if len(args) == 4 {
			trackerMethods.AddTask(args[2], args[3])
		} else {
			fmt.Println("Not enough parameters")
		}
	case "update":
		if len(args) == 4 {
			id, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Enter a number")
				return
			}

			trackerMethods.UpdateTask(id, args[3], "")

		} else if len(args) == 5 {
			status, err := trackerMethods.CorrectStatus(args[4])
			if err != nil {
				fmt.Println(err)
				return
			}

			id, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Enter a number")
				return
			}

			trackerMethods.UpdateTask(id, args[3], status)

		} else {
			fmt.Println("Not enough parameters")
		}
	case "list":
		if len(args) > 2 {
			fmt.Println("There are too many parameters")
			return
		}

		tasks, err := trackerMethods.GetAllLines()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, task := range tasks {
			fmt.Println(task)
		}
	case "delete":
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Enter a number")
			return
		}
		fmt.Println(trackerMethods.DeleneTask(id))
	}
}
