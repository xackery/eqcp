package config

const (
	IsDiscoveredOnly = 1
)

var (
	config map[int]string
)

func init() {
	config = make(map[int]string)
	config[IsDiscoveredOnly] = "0"
}
