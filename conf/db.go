package conf

type DbConfig struct {
	Alias      string
	DriverName string
	Database   string
	Host       string
	Port       string
	User       string
	Password   string
	MaxOpen    int
	MaxIdle    int
}
