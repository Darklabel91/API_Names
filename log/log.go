package log

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"os"
	"strings"
	"time"
)

func exportLog(filename string) error {
	//open the file
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Loop through the lines
	var logs []models.Log
	for scanner.Scan() {
		line := scanner.Text()
		log, _ := logBreaker(line)
		logs = append(logs, log)
	}

	//check for any errors during scanning
	if err = scanner.Err(); err != nil {
		return err
	}

	//uploads to the database
	database.Db.Create(&logs)

	// Return the line count and no error
	return nil
}

func logBreaker(message string) (models.Log, error) {
	stp1 := strings.Replace(message, "[GIN] ", "", -1)

	stp2 := strings.Split(strings.TrimSpace(stp1), "|")
	if len(stp2) != 5 {
		fmt.Println("oi")
		return models.Log{}, errors.New("unexpected string length")
	}

	stp3 := strings.Split(stp2[4], " ")
	if len(stp3) < 2 {
		return models.Log{}, errors.New("unexpected string length on second division")
	}

	return models.Log{
		Time:    strings.TrimSpace(stp2[0]),
		Status:  strings.TrimSpace(stp2[1]),
		Latency: strings.TrimSpace(stp2[2]),
		IP:      strings.TrimSpace(stp2[3]),
		Method:  strings.TrimSpace(stp3[1]),
		Path:    strings.TrimSpace(stp3[len(stp3)-1]),
	}, nil
}

func StartExportLog(filename string) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := exportLog(filename)
			if err != nil {
				fmt.Println("Error exporting log:", err)
			}
		}
	}
}
