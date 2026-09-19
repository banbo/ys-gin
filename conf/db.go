package conf

type DbConfig struct {
	Alias      string
	DriverName string
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	Charset    string
	Loc        string // DSN loc（DATETIME 读写时区）：Local（默认，兼容旧行为）/ UTC / Asia/Shanghai；配置写原值，拼 DSN 时自动转义
	MaxOpen    int
	MaxIdle    int
}
