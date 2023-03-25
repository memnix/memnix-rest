package misc

import (
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/memnix/memnix-rest/infrastructures"
	"sync"
	"time"
)

var logsChan = make(chan write.Point, 100) // buffered channel

// LogWriter is a custom writer for logs to be used
// It implements the io.Writer interface
type LogWriter struct{}

// Write is the implementation of the io.Writer interface
// It writes the log to the logsChan channel to be processed by the LogWorker
func (LogWriter) Write(p write.Point) (int, error) {
	logsChan <- p // write to channel
	return 1, nil
}

// LogWorker is the worker that will process the logsChan channel
func LogWorker(logs <-chan write.Point, wg *sync.WaitGroup) {
	defer wg.Done() // signal that we are done

	writeAPI := (*infrastructures.GetInfluxDBClient()).WriteAPI("memnix", "logs") // get writeAPI

	go func() {
		// Flush every 5s
		for range time.Tick(5 * time.Second) {
			writeAPI.Flush()
		}
	}()

	// Infinite loop
	for {
		logContent := <-logs
		writeAPI.WritePoint(&logContent)
	}
}

// CreateLogger creates a new LogWorker in a goroutine and waits for it to finish
func CreateLogger() {
	var wg sync.WaitGroup       // wait group
	wg.Add(1)                   // add one to wait group
	go LogWorker(logsChan, &wg) // start LogWorker in a goroutine
	wg.Wait()                   // wait for LogWorker to finish
}
