//
// get add and remove item events by parsing the mpd log
//

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type LogEvents struct {
	MpdReady chan bool
	EsReady  chan bool
	added    chan string
	deleted  chan string
}

var (
	addedString   = "update: added "
	deletedString = "update: removing "
)

// process to read log to create add and remove events
func NewLogEventParser(logFile string) (*LogEvents, error) {
	fmt.Printf("Create MPD log pipe: %s\n", logFile)
	syscall.Mkfifo(logFile, 0600)

	f, err := os.OpenFile(logFile, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)

	e := &LogEvents{
		added:   make(chan string),
		deleted: make(chan string),
	}

	go e.readLog(reader)

	return e, nil
}

// parse logs and send items to add and remove channels
func (e *LogEvents) readLog(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		fmt.Printf("%s", line)

		if strings.Contains(line, addedString) {
			str := strings.Split(line, addedString)
			e.added <- strings.TrimSuffix(str[len(str)-1], "\n")

		} else if strings.Contains(line, deletedString) {
			str := strings.Split(line, deletedString)
			e.deleted <- strings.TrimSuffix(str[len(str)-1], "\n")
		}
	}
}
