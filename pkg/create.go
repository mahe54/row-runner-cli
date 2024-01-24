package pkg

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// Creator is an interface that needs to be implemented by the user of the library
//
// import "github.com/vbauerster/mpb/v8" <- progress bar library
//
//	type Creator interface {
//		CreateFromInput(input interface{}, bar *mpb.Bar, progress chan<- int, cancel <-chan struct{}, statusChan chan<- Status, status *Status)
//		ConvertToInput(data []string) interface{}
//	}
//
// Inputs is any csv file that can be converted to a slice of strings
// SemaphoreSize is the number of goroutines/csv lines that can run in parallel
// LogFileName is the name of the log file
func Start(creator Creator, inputs []interface{}, semaphoreSize int, logFileName string) {

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	p := mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithWidth(25))
	progressChannels := make([]chan int, len(inputs))
	cancelChannels := make([]chan struct{}, len(inputs))
	statusChannels := make([]chan Status, len(inputs))
	semaphore := make(chan struct{}, semaphoreSize)

	totalInputs := len(inputs)
	bars := make([]*mpb.Bar, totalInputs)
	statuses := make([]*Status, len(inputs))

	for i, input := range inputs {
		currentInputNumber := i + 1
		progress := make(chan int)
		progressChannels[i] = progress
		statusChan := make(chan Status)
		statusChannels[i] = statusChan
		cancel := make(chan struct{})
		cancelChannels[i] = cancel
		status := &Status{}
		status.Set(Starting)
		statuses[i] = status

		barText := fmt.Sprintf("Processing %d/%d", currentInputNumber, totalInputs)
		bar := p.New(100,
			mpb.BarStyle().Lbound("[\u001b[32;1m█").Filler("█").Tip("█\u001b[36;1m").Padding("\u001b[0m░").Rbound("\u001b[0m]╟"),
			mpb.PrependDecorators(
				// decor.OnCondition(statusUpdate(&statuses[i], barText, 5), string(statuses[i].Current) == string(Working)),
				// decor.OnCondition(statusUpdate(&statuses[i], barText, 5), string(statuses[i].Current) == string(Failed)),

				// decor.Name(string(statuses[i].Current)+" "+barText), //, decor.WC{W: len(barText) + 1, C: decor.DidentRight}),
				// decor.OnComplete(
				// 	decor.OnCondition(statusUpdate(&statuses[i], barText, 5), string(statuses[i].Current) == string(Failed)), "jh",
				// ),

				decor.Conditional(string(statuses[i].Current) == string(Starting),
					statusUpdate(statuses[i], barText, 11, decor.WC{W: len(string(statuses[i].Current)) + len(barText) + 2}), statusUpdate(statuses[i], barText, 9, decor.WC{W: len(string(statuses[i].Current)) + len(barText) + 3})),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WC{W: 5}),
			),
		)

		bars[i] = bar

		semaphore <- struct{}{}
		go func(input interface{}, bar *mpb.Bar, progress chan<- int, cancel <-chan struct{}, statusChan chan<- Status, status *Status) {
			defer func() { <-semaphore }()
			creator.CreateFromInput(input, bar, progress, cancel, statusChan, status)
			close(progress)
		}(input, bar, progress, cancel, statusChan, status)

		go func(bar *mpb.Bar, statusChan <-chan Status, status *Status) {
			for newStatus := range statusChan {

				if newStatus.Current == Completed {
					bar.SetCurrent(100)
					bar.Completed()
				}
				if newStatus.Current == Failed {
					bar.Abort(false)
				}
				*status = newStatus
			}
		}(bar, statusChan, status)

		go func(bar *mpb.Bar, progress <-chan int, status *Status) {
			for p := range progress {
				bar.SetCurrent(int64(p))
				if p == 100 {
					*status = Status{Current: Completed}
				}
			}
		}(bar, progress, status)
	}

	p.Wait()
}

func statusUpdate(status *Status, text string, ws uint, wcc ...decor.WC) decor.Decorator {

	f := func(s decor.Statistics) string {
		return fmt.Sprintf("%s %s", status.Current, text)
		// if status.Current == Working {
		// 	return string(Working) + " " + string(text)
		// }
		// if status.Current == Completed {
		// 	return string(Completed) + " " + string(text)
		// }
		// return string(Failed) + " " + string(text)
	}
	return decor.Any(f, wcc...)
}
