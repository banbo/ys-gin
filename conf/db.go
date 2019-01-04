package conf

type DbConfig struct {
	DriverName string
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	MaxOpen    int
	MaxIdle    int
}
