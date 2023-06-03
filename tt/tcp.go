package tt

import (
	"net"
	"time"
)

// TCPinger 实现了 Pinger 接口，用于执行 TCP ping
type TCPinger struct {
	address        string
	count          int // 连接次数
	size           int // 发送的包大小
	sampleInterval int // 采样间隔
}

func NewTCPinger(config *Config) *TCPinger {
	return &TCPinger{
		address:        config.Address,
		count:          config.Count,
		size:           config.Size,
		sampleInterval: config.SampleInterval,
	}
}

func (t *TCPinger) Ping() {
	var (
		totalDelay      time.Duration
		elapsed         time.Duration
		minDelay        time.Duration
		maxDelay        time.Duration
		avgDelay        time.Duration
		successfulCount = 0
		failedCount     = 0

		printer = &SimplePrinter{}
	)
	hostPort, err := resolveAddress(t.address)
	if err != nil {
		printer.ErrorPrint(ErrorPrintData{
			Protocol: TCP,
			Address:  t.address,
			Error:    err,
		})
		return
	}
	printer.StartPrint(StartPrintData{
		Protocol:   TCP,
		Address:    t.address,
		RawAddress: hostPort,
	})
	for i := 1; i <= t.count; i++ {
		start := time.Now()
		conn, err := net.Dial("tcp", hostPort)
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: TCP,
				Address:  hostPort,
				Error:    err,
			})
			failedCount++
			continue
		}

		_, err = conn.Write(make([]byte, t.size))
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: TCP,
				Address:  hostPort,
				Error:    err,
			})
			failedCount++
			continue
		}
		defer conn.Close()
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: TCP,
				Address:  hostPort,
				Error:    err,
			})
			failedCount++
			continue
		}

		elapsed = time.Since(start)
		printer.ProcessPrint(ProcessPrintData{
			Protocol: TCP,
			Address:  t.address,
			Size:     t.size,
			Elapsed:  elapsed,
		})

		successfulCount++
		totalDelay += elapsed

		if minDelay == 0 || elapsed < minDelay {
			minDelay = elapsed
		}
		if elapsed > maxDelay {
			maxDelay = elapsed
		}
		time.Sleep(time.Duration(t.sampleInterval) * time.Second)
	}

	if successfulCount > 0 {
		avgDelay = totalDelay / time.Duration(successfulCount)
	}
	percentage := float64(failedCount/(failedCount+successfulCount)) * 100
	printer.ResultPrint(ResultPrintData{
		FailedCount:    failedCount,
		SucceededCount: successfulCount,
		MinDelay:       minDelay,
		AvgDelay:       avgDelay,
		MaxDelay:       maxDelay,
		Percentage:     percentage,
	})
}
