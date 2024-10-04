package destination

import (
	"fmt"
	"os"

	"github.com/morphy76/cgnlog/internal/event"
)

type outFormat struct {
	Header    string
	Separator string
	Footer    string
	Offset    int
	Parse     func(event.Event) string
	Extension string
}

var formats = map[string]outFormat{
	"html": {
		Header:    "<!DOCTYPE html><html><head><title>Log</title><style>body{font-size:0.75em;}table{width:100%;}td{border:1px solid darkgray;padding:0 3px;width:1%;}div{max-height:100px;overflow:scroll;text-wrap:wrap;}</style></head><body><table><thead><tr><th>Timestamp</th><th>Logger</th><th>Level</th><th>Tenant</th><th>Subscription</th><th>Trace</th><th>Message</th></tr></thead><tbody>",
		Separator: "\n",
		Footer:    "</tbody></table></body></html>",
		Offset:    0,
		Parse: func(ev event.Event) string {
			return fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td style='width:99%%;max-width:300px;'><div>%s</div></td></tr>", ev.Timestamp, ev.Logger, ev.Level, ev.MDC.Tenant, ev.MDC.Subscription, ev.MDC.Trace, ev.Message)
		},
		Extension: "html",
	},
	"json": {
		Header:    "[",
		Separator: ",\n",
		Footer:    "]",
		Offset:    1,
		Parse:     func(ev event.Event) string { return ev.Line },
		Extension: "json",
	},
}

func WriteTemporaryFile(lineChan chan event.Event, format string, doneChan chan string, progressChan chan bool) error {

	progressChan <- true

	useFormat := formats[format]

	tempFile, err := os.CreateTemp("", "cgnlog_*."+useFormat.Extension)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	progressChan <- true

	size, err := tempFile.WriteString(useFormat.Header)

	progressChan <- true

	for ev := range lineChan {
		line := ev.Line
		if line == "<EOF>" {
			tempFile.WriteAt([]byte(useFormat.Footer), int64(size-useFormat.Offset))
			doneChan <- tempFile.Name()
			return nil
		} else {
			bytes, _ := tempFile.WriteString(useFormat.Parse(ev) + useFormat.Separator)
			size += bytes
		}
	}

	progressChan <- true

	return nil
}
