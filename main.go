package main

import (
	"bufio"
	f "file-directory/utilities/files"
	t "file-directory/utilities/times"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	Description string
}

var userCommands []string

func addUserCommand(command string) {
	userCommands = append(userCommands, command)
}

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}

	fmt.Println("File directory system to do file manipulation")
	fmt.Println("-------------------------------")
	fmt.Println("Current Directory:", currentPath)
	fmt.Println("-------------------------------")
	displayOptions()

	for {
		fmt.Print("Enter option: ")
		var input int
		fmt.Scanln(&input)

		switch input {
		case 1:
			displayOptions()
		case 2:
			listFiles(currentPath)
		case 3:
			fmt.Print("Enter directory name: ")
			var dirName string
			fmt.Scanln(&dirName)
			currentPath = changeDirectory(currentPath, dirName)
		case 4:
			createDirectory()
		case 5:
			readFile()
		case 6:
			copyFile()
		case 7:
			editFile()
		case 8:
			createFile()
		case 9:
			deleteFile()
		case 10:
			exitCLI()
		default:
			fmt.Print("Unknown command: ", input)
		}
		fmt.Println("\n-------------------------------")
	}
}

func displayOptions() {
	commands := []Command{
		{"Display Options"},
		{"List files in the current directory"},
		{"Change directory"},
		{"Create a new directory"},
		{"Read a file"},
		{"Copy File"},
		{"Edit File"},
		{"Create file"},
		{"Delete file"},
		{"Exit the CLI"},
	}

	fmt.Println("Options:")
	for index, cmd := range commands {
		fmt.Printf("  %d - %s\n", index+1, cmd.Description)
	}
}

func listFiles(path string) {
	addUserCommand("List files in the current directory")
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println("Files in " + path + " :")
	fmt.Printf("%-15s %-25s %-20s %-12s %s\n", "Permissions", "Name", "Size(bytes)", "Type", "Last updated")
	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("Error getting file info for %s: %s\n", file.Name(), err)
			continue
		}
		fileType := "File"
		if fileInfo.IsDir() {
			fileType = "Directory"
		}
		fmt.Printf("%-15s %-25s %-20s %-12s %s\n",
			f.GetPermissions(fileInfo),
			file.Name(),
			f.FormatSize(fileInfo.Size(), fileType),
			fileType,
			t.FormatRelativeTime(fileInfo.ModTime()),
		)
	}
}

func changeDirectory(currentPath, dirName string) string {
	addUserCommand("Change directory")
	fmt.Print("Enter directory name: ")
	newPath := filepath.Join(currentPath, dirName)
	_, err := os.Stat(newPath)
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return currentPath
	}

	return newPath
}

func createDirectory() {
	addUserCommand("Create a new directory")
	fmt.Print("Enter the directory name: ")
	var dirName string
	fmt.Scanln(&dirName)

	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	fmt.Println("Directory created successfully!")
}

func readFile() {
	addUserCommand("Read a file")
	fmt.Print("Enter the filename to read:")
	var filename string
	fmt.Scanln(&filename)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func copyFile() {
	addUserCommand("Copy a file")
	fmt.Println("Enter the source filename:")
	var sourceFilename string
	fmt.Scanln(&sourceFilename)

	fmt.Println("Enter the destination filename:")
	var destFilename string
	fmt.Scanln(&destFilename)

	source, err := os.Open(sourceFilename)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			fmt.Println("File closing error: ", err)
		}
	}(source)

	destination, err := os.Create(destFilename)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}

	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			fmt.Println("File manipulation error: ", err)
		}
	}(destination)

	_, err = io.Copy(destination, source)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	fmt.Println("File copied successfully!")
}

func editFile() {
	addUserCommand("Edit file")
	fmt.Print("Enter the filename to edit: ")
	var filename string
	fmt.Scanln(&filename)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // Default to nano if EDITOR is not set
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening file in editor:", err)
		return
	}

	fmt.Println("File edited successfully!")
}

func createFile() {
	addUserCommand("Create a new file")
	fmt.Print("Enter the filename to create: ")
	var filename string
	fmt.Scanln(&filename)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	fmt.Println("Enter the content of the file. Press Ctrl+D (Unix-like systems) or Ctrl+Z (Windows) to finish:")
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	_, err = file.Write(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File created successfully!")
}

func deleteFile() {
	addUserCommand("Delete a file")
	fmt.Print("Enter filename to be deleted: ")
	var filename string
	fmt.Scanln(&filename)

	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Error deleting file")
		return
	}

	fmt.Println("File deleted successfully")
}

func exitCLI() {
	var input string
	for {
		input = ""
		fmt.Print("Do you want to exit application? [y/N]: ")
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input == "y" || input == "n" {
			if input == "y" {
				fmt.Println("Exiting CLI file directory system.")
				os.Exit(1)
			} else {
				fmt.Println("\n-------------------------------")
				break
			}
		} else {
			fmt.Println("Incorrect Option.")
			break
		}
	}
	displayOptions()
}
