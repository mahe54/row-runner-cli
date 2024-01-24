package pkg

import "github.com/vbauerster/mpb/v8"

type Creator interface {
	CreateFromInput(input interface{}, bar *mpb.Bar, progress chan<- int, cancel <-chan struct{}, statusChan chan<- Status, status *Status)
	ConvertToInput(data []string) interface{}
}

type Status struct {
	Current Current
}

type Current string

const (
	Starting  Current = "starting"
	Working   Current = "working"
	Error     Current = "error"
	Retrying  Current = "retrying"
	Completed Current = "completed"
	Failed    Current = "failed"
)

func (s *Status) Set(status Current) {
	s.Current = status
}
