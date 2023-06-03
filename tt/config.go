package tt

const (
	DefaultProtocol       = ICMP
	DefaultCount          = 5
	DefaultSize           = 64
	DefaultSampleInterval = 1
)

type ProtocolType string

const (
	TCP  ProtocolType = "tcp"
	ICMP ProtocolType = "icmp"
	HTTP ProtocolType = "http"
)

type Config struct {
	Protocol       ProtocolType
	Address        string
	Count          int
	Size           int
	SampleInterval int
}
