package source

import (
	"bufio"
	"os"

	"github.com/morphy76/cgnlog/internal/event"
)

func ReadInputFile(fileName string, rowsChan chan event.Event, progressChan chan bool) error {

	progressChan <- true

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	progressChan <- true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ev, isJson := event.ToJson(line)
		ev.Line = line
		if isJson {
			rowsChan <- ev
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	rowsChan <- event.Event{Line: "<EOF>"}
	progressChan <- true

	return nil
}
