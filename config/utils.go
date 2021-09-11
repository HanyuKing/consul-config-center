package config

var (
	consulAddr = "127.0.0.1:8500"
)

func GetConsulAddr() string {
	return consulAddr
}
