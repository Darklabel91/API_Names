package log

import (
	"bytes"
	"fmt"
	database2 "github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"os"
	"sync"
)

//user, err := c.Cookie("user")
//if err != nil {
//	c.AbortWithStatus(http.StatusUnauthorized)
//	return
//}

const logBufferSize = 1000

type databaseWriter struct {
	buffer   bytes.Buffer
	db       *gorm.DB
	mutex    sync.Mutex
	logCount int
}

func (w *databaseWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	n, err = w.buffer.Write(p)
	if err != nil {
		return
	}

	// Count the number of logs written to the buffer
	w.logCount += bytes.Count(p, []byte{'\n'})

	// If the buffer is full, flush it to the database
	if w.logCount >= logBufferSize {
		err = w.flush()
		if err != nil {
			return
		}
	}

	return
}

func (w *databaseWriter) flush() error {
	// Split the buffer into individual log messages
	logs := bytes.Split(w.buffer.Bytes(), []byte{'\n'})
	for _, log := range logs {
		// Skip empty log messages
		if len(log) == 0 {
			continue
		}

		// Insert the log into the database
		err := w.db.Create(&models.Log{Message: string(log)}).Error
		if err != nil {
			return err
		}
	}

	// Reset the buffer and log count
	w.buffer.Reset()
	w.logCount = 0

	return nil
}

func LocalLog() gin.HandlerFunc {
	// Create the file writer and database writer
	file, err := os.Create("log/gin.txt")
	if err != nil {
		fmt.Println("Log not created")
		return nil
	}
	database := &databaseWriter{db: database2.InitDb()}

	// Use a multi-writer to write to both the file and database writers
	writer := io.MultiWriter(file, database)

	// Create the gin logger with the custom writer
	logger := gin.LoggerWithWriter(writer)

	return func(c *gin.Context) {
		// Call the logger middleware
		logger(c)

		// Flush the logs to the database if necessary
		if database.logCount >= logBufferSize {
			err = database.flush()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
