package tt

import (
	"fmt"
	"time"
)

type ProcessData struct {
	Elapsed time.Duration
	Address string
	Size    int
}

type PingPrinter interface {
	ProcessPrint(elapsed time.Duration, address string, size int)
	ResultPrint()
}

type StartPrintData struct {
	Protocol   ProtocolType
	Address    string
	RawAddress string
}

type ProcessPrintData struct {
	Protocol ProtocolType
	Address  string
	Size     int
	Elapsed  time.Duration
}

type ErrorPrintData struct {
	Protocol ProtocolType
	Address  string
	Error    error
}

type ResultPrintData struct {
	FailedCount    int
	SucceededCount int
	MinDelay       time.Duration
	AvgDelay       time.Duration
	MaxDelay       time.Duration
	Percentage     float64
}

type SimplePrinter struct {
}

func (s *SimplePrinter) StartPrint(data StartPrintData) {
	fmt.Printf("PING %s %s (%s)\n", data.Protocol, data.Address, data.RawAddress)
}

func (s *SimplePrinter) ProcessPrint(data ProcessPrintData) {
	fmt.Printf("%v bytes to %v  %v\n", data.Size, data.Address, data.Elapsed)
}

func (s *SimplePrinter) ResultPrint(data ResultPrintData) {
	fmt.Printf("%v failed, %v succeeded, %.2f%% loss\n", data.FailedCount,
		data.SucceededCount, data.Percentage)
	fmt.Printf("min/avg/max =  %.2v/%.2v/%.2v ms\n", data.AvgDelay.Milliseconds(),
		data.MinDelay.Milliseconds(), data.MaxDelay.Milliseconds())
}

func (s *SimplePrinter) ErrorPrint(data ErrorPrintData) {
	fmt.Printf("%s to %s failed: %v\n", data.Protocol, data.Address, data.Error)
}
