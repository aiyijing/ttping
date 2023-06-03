package tt

import (
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ICMPPinger 实现了 Pinger 接口，用于执行 ICMP ping
type ICMPPinger struct {
	address        string
	count          int // 连接次数
	size           int // 发送的包大小
	sampleInterval int // 采样间隔
}

// NewICMPPinger 创建一个新的 ICMPPinger 实例
func NewICMPPinger(config *Config) *ICMPPinger {
	return &ICMPPinger{
		address:        config.Address,
		count:          config.Count,
		size:           config.Size,
		sampleInterval: config.SampleInterval,
	}
}

// Ping 执行 ICMP ping
func (i *ICMPPinger) Ping() {
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

	for j := 1; j <= i.count; j++ {
		start := time.Now()

		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}
		defer conn.Close()

		ipAddr, err := net.ResolveIPAddr("ip4", i.address)
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  j,
				Data: make([]byte, i.size),
			},
		}
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}

		_, err = conn.WriteTo(msgBytes, ipAddr)
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}

		respBuf := make([]byte, 1500)
		err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}

		n, _, err := conn.ReadFrom(respBuf)
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}

		_, err = icmp.ParseMessage(1, respBuf[:n])
		if err != nil {
			printer.ErrorPrint(ErrorPrintData{
				Protocol: ICMP,
				Address:  i.address,
				Error:    err,
			})
			failedCount++
			continue
		}

		elapsed = time.Since(start)
		printer.ProcessPrint(ProcessPrintData{
			Protocol: ICMP,
			Address:  i.address,
			Size:     i.size,
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
		time.Sleep(time.Duration(i.sampleInterval) * time.Second)
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
