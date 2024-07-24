package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type FileDiff struct {
	FilePath            string
	CompareFileSize     string
	CompareLastModified time.Time
	CurrentFileSize     string
	CurrentLastModified time.Time
}

func main() {
	// Define command line flag
	compareCommit := flag.String("from", "", "commit hash to compare with")
	currentCommit := flag.String("to", "", "current commit hash")
	flag.Parse()

	// Validate flags
	if *currentCommit == "" || *compareCommit == "" {
		log.Fatalln("Both current and compare commit hashes are required")
	}

	// Get the list of changed files
	changedFiles, err := getChangedFiles(*currentCommit, *compareCommit)
	if err != nil {
		log.Fatalf("Error getting changed files: %v\n", err)
	}

	// Get file details
	var fileDiffs []FileDiff
	for _, file := range changedFiles {
		fileDiff, err := getFileDetails(file, *currentCommit, *compareCommit)
		if err != nil {
			log.Printf("Error getting file details: %v\n", err)
			continue
		}
		fileDiffs = append(fileDiffs, fileDiff)
	}

	// Write to csv
	err = writeToCSV(fileDiffs, *currentCommit, *compareCommit)
	if err != nil {
		log.Fatalf("Error writing to csv: %v\n", err)
	}

	log.Println("Successfully written to csv")
}

func getChangedFiles(current, compare string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", compare, current)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files, nil
}

func getFileDetails(filePath, currentCommit, compareCommit string) (FileDiff, error) {
	fileDetail := FileDiff{
		FilePath: filePath,
	}

	// checkout to compare commit
	if err := checkoutCommit(compareCommit); err != nil {
		return fileDetail, err
	}

	// get file details at compare commit
	compareFileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Skipping file info at compare commit: %v\n", err)
		fileDetail.CompareFileSize = "N/A"
		fileDetail.CompareLastModified = time.Time{}
	}

	// checkout back to current commit
	if err := checkoutCommit(currentCommit); err != nil {
		return fileDetail, err
	}

	// get file details at current commit
	currentFileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Skipping file info at current commit: %v\n", err)
		fileDetail.CurrentFileSize = "N/A"
		fileDetail.CurrentLastModified = time.Time{}
	}

	// set file details
	if compareFileInfo != nil {
		fileDetail.CompareFileSize = fmt.Sprintf("%.2f", float64(compareFileInfo.Size())/1024)
		fileDetail.CompareLastModified = compareFileInfo.ModTime()
	}
	if currentFileInfo != nil {
		fileDetail.CurrentFileSize = fmt.Sprintf("%.2f", float64(currentFileInfo.Size())/1024)
		fileDetail.CurrentLastModified = currentFileInfo.ModTime()
	}

	return fileDetail, nil
}

func checkoutCommit(commit string) error {
	cmd := exec.Command("git", "checkout", commit)
	return cmd.Run()
}

func writeToCSV(fileDiffs []FileDiff, currentCommit, compareCommit string) error {
	if len(currentCommit) > 5 {
		currentCommit = currentCommit[:5]
	}
	if len(compareCommit) > 5 {
		compareCommit = compareCommit[:5]
	}

	file, err := os.Create(fmt.Sprintf("diff_%s_%s.csv", compareCommit, currentCommit))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Printf("Error closing file: %v\n", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"from", compareCommit, "", "to", currentCommit, ""})
	if err != nil {
		return err
	}

	err = writer.Write([]string{"File Name", "File Size (KB)", "Modified", "File Name", "File Size (KB)", "Modified"})
	if err != nil {
		return err
	}

	// Write data
	for _, diff := range fileDiffs {
		err = writer.Write([]string{
			diff.FilePath,
			diff.CompareFileSize,
			diff.CompareLastModified.Format("02/01/2006"),
			diff.FilePath,
			diff.CurrentFileSize,
			diff.CurrentLastModified.Format("02/01/2006"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
