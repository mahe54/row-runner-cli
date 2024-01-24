package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	lib "github.com/mahe54/row-runner-cli/pkg"
	"github.com/vbauerster/mpb/v8"
)

type Input struct {
	Name        string
	Description string
	Value       string
}

type CreatorExampleImpl struct{}

func (t CreatorExampleImpl) ProcessInput(input interface{}, bar *mpb.Bar, progress chan<- int, cancel <-chan struct{}, statusChan chan<- lib.Status, status *lib.Status) {
	// Simulate long-running task, updates progress channel every 500-1500ms by 10%
	for i := 1; i <= 10; i++ {
		select {
		case progress <- i * 10:
			statusChan <- lib.Status{Current: lib.Working}
			sleepTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
			time.Sleep(sleepTime)
		case <-cancel:
			return
		}
	}
	progress <- 100
}

func (t CreatorExampleImpl) ConvertToInput(data []string) interface{} {
	return Input{Name: data[0], Description: data[1], Value: data[2]}
}

func main() {
	inputFile := flag.String("file", "./input.csv", "Input file to process")
	semaphoreSize := flag.Int("s", 2, "Semaphore size")
	logFileName := flag.String("log", "app.log", "Log file name")

	flag.Parse()

	creator := CreatorExampleImpl{}
	inputs, err := lib.ReadDataFromFile(*inputFile, creator)
	if err != nil {
		log.Fatal(err)
	}

	lib.Start(creator, inputs, *semaphoreSize, *logFileName)
}
