package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"ttping/tt"
)

var cfg = &tt.Config{}
var rootCmd = &cobra.Command{
	Use:   "ttping [address]",
	Short: "A network pinging tool",
	Args: func(cmd *cobra.Command, args []string) error {
		if printVersion {
			return nil
		}

		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		err := validateArgs(args[0], cfg.Protocol)
		if err != nil {
			return err
		}
		cfg.Address = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion {
			fmt.Println("ttping version:", Version)
			return
		}
		tt.NewPinger(cfg).Ping()
	},
}

var Version = "0.1.0"
var printVersion bool

func init() {
	rootCmd.Flags().StringVarP((*string)(&cfg.Protocol), "protocol", "t", string(tt.DefaultProtocol), "指定使用的协议，有效值为 tcp、icmp 或 http，默认为 icmp")
	rootCmd.Flags().IntVarP(&cfg.Count, "count", "c", tt.DefaultCount, "指定连接尝试的次数，默认为 5")
	rootCmd.Flags().IntVarP(&cfg.Size, "size", "s", tt.DefaultSize, "指定要发送的数据包大小，默认为 64")
	rootCmd.Flags().IntVarP(&cfg.SampleInterval, "interval", "i", tt.DefaultSampleInterval, "指定采样间隔（秒），默认为 1")
	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "打印版本信息")
}

func validateArgs(address string, protocol tt.ProtocolType) error {
	switch protocol {
	case tt.TCP:
		if !tt.ContainsPort(address) {
			return errors.New("TCP 协议必须包含端口")
		}
	case tt.HTTP:
		u, err := url.Parse(address)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return errors.New("HTTP 协议必须以 http:// 或 https:// 开头")
		}
	case tt.ICMP:
		if tt.ContainsPort(address) {
			return errors.New("ICMP 协议不支持指定端口")
		}
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
