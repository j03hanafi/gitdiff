package main

import (
	"bytes"
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
	remark := flag.String("remark", "", "latest commit hash")
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
	fileDiffs := make([]FileDiff, len(changedFiles))

	for i, file := range changedFiles {
		fileDiffs[i].FilePath = file
	}

	err = getFileDetails(fileDiffs, *currentCommit, *compareCommit)
	if err != nil {
		log.Fatalf("Error getting file details")
	}

	// Write to csv
	err = writeToCSV(fileDiffs, *currentCommit, *compareCommit, *remark)
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

func getFileDetails(fileDiffs []FileDiff, currentCommit, compareCommit string) error {
	// Get commit dates
	currentCommitDate, err := getCommitDate(currentCommit)
	if err != nil {
		return err
	}
	compareCommitDate, err := getCommitDate(compareCommit)
	if err != nil {
		return err
	}

	// Checkout to the current commit once
	if err := checkoutCommit(currentCommit); err != nil {
		return err
	}

	// Loop over fileDiffs to get file info at the current commit
	for i := range fileDiffs {
		filePath := fileDiffs[i].FilePath

		currentFileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Skipping file info at current commit for %s: %v\n", filePath, err)
			fileDiffs[i].CurrentFileSize = "-"
			fileDiffs[i].CurrentFileType = ""
		} else {
			fileDiffs[i].CurrentFileSize = fmt.Sprintf("%.2f", float64(currentFileInfo.Size())/1024)
			fileDiffs[i].CurrentFileType = filepath.Ext(filePath)
		}
		fileDiffs[i].CurrentLastModified = currentCommitDate
	}

	// Checkout to the compare commit once
	if err := checkoutCommit(compareCommit); err != nil {
		return err
	}

	// Loop over fileDiffs to get file info at the compare commit
	for i := range fileDiffs {
		filePath := fileDiffs[i].FilePath

		compareFileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Skipping file info at compare commit for %s: %v\n", filePath, err)
			fileDiffs[i].CompareFileSize = "-"
			fileDiffs[i].CompareFileType = fileDiffs[i].CurrentFileType
		} else {
			fileDiffs[i].CompareFileSize = fmt.Sprintf("%.2f", float64(compareFileInfo.Size())/1024)
			fileDiffs[i].CompareFileType = filepath.Ext(filePath)
		}
		fileDiffs[i].CompareLastModified = compareCommitDate
	}

	return nil
}

func checkoutCommit(commit string) error {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func getCommitDate(commit string) (time.Time, error) {
	cmd := exec.Command("git", "show", "-s", "--format=%ci", commit)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return time.Time{}, err
	}
	dateStr := strings.TrimSpace(out.String())
	// Parse the date string into time.Time
	commitDate, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return commitDate, nil
}

func writeToCSV(fileDiffs []FileDiff, currentCommit, compareCommit, remark string) error {
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
	err = writer.Write([]string{"", "from", compareCommit, "", "", "", "to", currentCommit})
	if err != nil {
		return err
	}

	err = writer.Write([]string{"No", "File Name", "File Type", "Date Modified", "File Size (KB)", "No", "File Name", "File Type", "Date Modified", "File Size (KB)", "Remark"})
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
			strconv.Itoa(i + 1),
			diff.FilePath,
			diff.CurrentFileType,
			diff.CurrentLastModified.Format("02 Jan 2006"),
			diff.CurrentFileSize,
			remark,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
