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
	MaxOpen    int
	MaxIdle    int
	Slaves     []struct {
		Host     string
		Port     string
		User     string
		Password string
	}
}
