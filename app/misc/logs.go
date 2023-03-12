package misc

import (
	"io"
	"path"
	"sync"

	"github.com/memnix/memnix-rest/config"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logsChan = make(chan string, config.LogChannelSize) // buffered channel

// LogWriter is a custom writer for logs to be used
// It implements the io.Writer interface
type LogWriter struct{}

// Write is the implementation of the io.Writer interface
// It writes the log to the logsChan channel to be processed by the LogWorker
func (e LogWriter) Write(p []byte) (int, error) {
	logsChan <- string(p) // write to channel
	return 0, nil
}

// LogWorker is the worker that will process the logsChan channel
func LogWorker(logs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done() // signal that we are done
	// Infinite loop
	for {
		logContent := <-logs                                 // read from channel
		_, err := newRollingFile().Write([]byte(logContent)) // write to file
		if err != nil {
			// log error
			log.Error().Err(err).Msg("Error writing to file")
			continue
		}
	}
}

// CreateLogger creates a new LogWorker in a goroutine and waits for it to finish
func CreateLogger() {
	var wg sync.WaitGroup       // wait group
	wg.Add(1)                   // add one to wait group
	go LogWorker(logsChan, &wg) // start LogWorker in a goroutine
	wg.Wait()                   // wait for LogWorker to finish
}

// newRollingFile creates a new rolling file
// It uses the lumberjack package to create a rolling file
// It returns an io.Writer interface
func newRollingFile() io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join("./logs", "api.log"), // file path
		MaxBackups: config.MaxBackupLogFiles,       // files
		MaxSize:    config.MaxSizeLogFiles,         // megabytes
	}
}
