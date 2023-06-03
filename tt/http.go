package tt

import (
	"net/http"
	"time"
)

// HTTPPinger 实现了 Pinger 接口，用于执行 HTTP ping
type HTTPPinger struct {
	address        string
	count          int // 连接次数
	size           int // 发送的包大小
	sampleInterval int // 采样间隔
}

// NewHTTPPinger 创建一个新的 HTTPPinger 实例
func NewHTTPPinger(config *Config) *HTTPPinger {
	return &HTTPPinger{
		address:        config.Address,
		count:          config.Count,
		size:           config.Size,
		sampleInterval: config.SampleInterval,
	}
}

// Ping 执行 HTTP ping
func (h *HTTPPinger) Ping() {
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
	//hostPort, err := resolveAddress(h.address)
	//if err != nil {
	//	printer.ErrorPrint(ErrorPrintData{
	//		Protocol: TCP,
	//		Address:  h.address,
	//		Error:    err,
	//	})
	//	return
	//}
	printer.StartPrint(StartPrintData{
		Protocol:   HTTP,
		Address:    h.address,
		RawAddress: h.address,
	})
	for i := 1; i <= h.count; i++ {
		start := time.Now()
		resp, err := http.Get(h.address)
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: HTTP,
				Address:  h.address,
				Error:    err,
			})
			failedCount++
			continue
		}
		defer resp.Body.Close()

		// 读取响应内容，但不做任何操作
		_, _ = resp.Body.Read(make([]byte, h.size))

		elapsed = time.Since(start)
		printer.ProcessPrint(ProcessPrintData{
			Protocol: HTTP,
			Address:  h.address,
			Size:     h.size,
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

		// 暂停指定的采样间隔
		time.Sleep(time.Duration(h.sampleInterval) * time.Second)
	}

	if successfulCount > 0 {
		avgDelay = totalDelay / time.Duration(successfulCount)
	}
	percentage := float64(failedCount) / float64(failedCount+successfulCount) * 100
	printer.ResultPrint(ResultPrintData{
		FailedCount:    failedCount,
		SucceededCount: successfulCount,
		MinDelay:       minDelay,
		AvgDelay:       avgDelay,
		MaxDelay:       maxDelay,
		Percentage:     percentage,
	})
}
