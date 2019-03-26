package conf

type ApiConfig struct {
	HttpPort    string
	RpcPort     string
	RunMode     string
	LogPath     string
	LogLevel    string
	WorkerID    int64
	ParamSecret string
	Dbs         string
}
