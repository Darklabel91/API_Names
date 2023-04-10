package models

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strings"
)

type Log struct {
	gorm.Model
	Time    string
	Status  string
	Latency string
	IP      string
	Method  string
	Path    string
}

func (*Log) Upload(db *gorm.DB, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//read the file content
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	//don't upload if content is 0
	if len(content) == 0 {
		return nil
	}

	// replace null bytes with empty string
	content = bytes.ReplaceAll(content, []byte{0}, []byte{})

	// write the modified content back to the file
	err = ioutil.WriteFile(file.Name(), content, 0666)
	if err != nil {
		return err
	}

	// Create a scanner
	scanner := bufio.NewScanner(file)

	// Read every line
	var logs []Log
	for scanner.Scan() {
		// Process the line
		line := scanner.Text()
		logLine, err := breakLog(line)
		if err != nil {
			return err
		}

		logs = append(logs, logLine)
	}

	//save the document to the database
	if len(logs) != 0 {
		err = db.Create(&logs).Error
		if err != nil {
			return err
		}

		//clear the file
		err = file.Truncate(0)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("upload scanner return was 0")
}

func breakLog(logLine string) (Log, error) {
	split1 := strings.Split(logLine, "|")
	if len(split1) != 5 {
		return Log{}, errors.New("unexpected length on first splitting")
	}

	split2 := strings.Split(split1[4], " ")
	if len(split2) < 7 {
		fmt.Println(len(split2))
		return Log{}, errors.New("unexpected length on second splitting")
	}

	return Log{
		Time:    strings.Replace(strings.TrimSpace(split1[0]), "[GIN]", "", -1),
		Status:  strings.TrimSpace(split1[1]),
		Latency: strings.TrimSpace(split1[2]),
		IP:      strings.TrimSpace(split1[3]),
		Method:  strings.TrimSpace(split2[1]),
		Path:    strings.TrimSpace(split2[len(split2)-1]),
	}, nil

}
