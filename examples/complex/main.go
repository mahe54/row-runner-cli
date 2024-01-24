// This is a sample on how you could implement the package in your own CLI
package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	lib "github.com/mahe54/generic-go-cli-file-input/pkg"
	"github.com/vbauerster/mpb/v8"
)

type Car struct {
	RegNumber string
	Brand     string
	Model     string
	Year      string
	Mileage   string
	Insured   bool
}

type MyCreatorImpl struct{}

func (t MyCreatorImpl) ProcessInput(input interface{}, bar *mpb.Bar, progress chan<- int, cancel <-chan struct{}, statusChan chan<- lib.Status, status *lib.Status) {

	car := input.(Car)
	retry := false

	sleepTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond

	//Simulation of work
	time.Sleep(sleepTime)
	statusChan <- lib.Status{Current: lib.Working}

	//Report progress
	progress <- 25

	//Simulate some more work
	time.Sleep(sleepTime)

	//Report progress
	progress <- 50

	//Simulate some even more work, that might fail 50% chance
	time.Sleep(sleepTime)
	if rand.Intn(100) < 50 {
		retry = true
		statusChan <- lib.Status{Current: lib.Error}
		log.Printf("Error processing car: %s", car.RegNumber)
		time.Sleep(sleepTime)
	}

	for retry {
		//Simulate some even more work, that might succed 50% chance to recover
		statusChan <- lib.Status{Current: lib.Retrying}
		time.Sleep(sleepTime)
		if rand.Intn(100) < 50 {
			statusChan <- lib.Status{Current: lib.Working}
			log.Printf("recovered: %s", car.RegNumber)
			time.Sleep(sleepTime)
			retry = false
		} else {
			time.Sleep(sleepTime)
			log.Printf("Error processing car: %s", car.RegNumber)
			statusChan <- lib.Status{Current: lib.Failed}
			return
		}
	}

	//Report progress
	time.Sleep(sleepTime)
	progress <- 75

	//Simulate some even more work
	time.Sleep(sleepTime)
	progress <- 85

	//Report progress
	progress <- 100
}

func (t MyCreatorImpl) ConvertToInput(data []string) interface{} {

	return Car{
		RegNumber: data[0],
		Brand:     data[1],
		Model:     data[2],
		Year:      data[3],
		Mileage:   data[4],
		Insured:   data[5] == "true",
	}
}

func main() {
	inputFile := flag.String("file", "./input.csv", "Input file to process")
	semaphoreSize := flag.Int("s", 2, "Semaphore size")
	logFileName := flag.String("log", "app.log", "Log file name")

	flag.Parse()

	creator := MyCreatorImpl{}
	inputs, err := lib.ReadDataFromFile(*inputFile, creator)
	if err != nil {
		log.Fatal(err)
	}

	lib.Start(creator, inputs, *semaphoreSize, *logFileName)
}
