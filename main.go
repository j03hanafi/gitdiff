package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileDiff struct {
	FilePath            string
	CompareFileType     string
	CompareFileSize     string
	CompareLastModified time.Time
	CurrentFileType     string
	CurrentFileSize     string
	CurrentLastModified time.Time
}

func main() {
	// Define command line flag
	compareCommit := flag.String("from", "", "commit hash to compare with")
	currentCommit := flag.String("to", "", "latest commit hash")
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
	cmd.Stderr = os.Stdout
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

	// checkout back to current commit
	if err := checkoutCommit(currentCommit); err != nil {
		return fileDetail, err
	}

	// get file details at current commit
	currentFileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Skipping file info at current commit: %v\n", err)
		fileDetail.CurrentFileSize = "-"
		fileDetail.CurrentLastModified = time.Time{}
	}

	// set file details
	if currentFileInfo != nil {
		fileDetail.CurrentFileSize = fmt.Sprintf("%.2f", float64(currentFileInfo.Size())/1024)
		fileDetail.CurrentLastModified = currentFileInfo.ModTime()
		fileDetail.CurrentFileType = filepath.Ext(filePath)
	}

	// checkout to compare commit
	if err := checkoutCommit(compareCommit); err != nil {
		return fileDetail, err
	}

	// get file details at compare commit
	compareFileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Skipping file info at compare commit: %v\n", err)
		fileDetail.CompareFileSize = "-"
		fileDetail.CompareLastModified = fileDetail.CurrentLastModified
		fileDetail.CompareFileType = fileDetail.CurrentFileType
	}

	// set file details
	if compareFileInfo != nil {
		fileDetail.CompareFileSize = fmt.Sprintf("%.2f", float64(compareFileInfo.Size())/1024)
		fileDetail.CompareLastModified = compareFileInfo.ModTime()
		fileDetail.CompareFileType = filepath.Ext(filePath)
	}

	return fileDetail, nil
}

func checkoutCommit(commit string) error {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Stderr = os.Stdout
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
	err = writer.Write([]string{"", "from", compareCommit, "", "", "to", currentCommit})
	if err != nil {
		return err
	}

	err = writer.Write([]string{"No", "File Name", "File Type", "Date Modified", "File Size (KB)", "File Name", "File Type", "Date Modified", "File Size (KB)", "Remarks"})
	if err != nil {
		return err
	}

	// Write data
	for i, diff := range fileDiffs {
		err = writer.Write([]string{
			strconv.Itoa(i + 1),
			diff.FilePath,
			diff.CompareFileType,
			diff.CompareLastModified.Format("02 Jan 2006"),
			diff.CompareFileSize,
			diff.FilePath,
			diff.CurrentFileType,
			diff.CurrentLastModified.Format("02 Jan 2006"),
			diff.CurrentFileSize,
			"Backend",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
