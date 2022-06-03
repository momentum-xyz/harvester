package mysql

// Config : structure to hold MYSQL configuration
type Config struct {
	Database string `yaml:"database" envconfig:"DB_DATABASE"`
	Host     string `yaml:"host" envconfig:"DB_HOST"`
	Port     uint   `yaml:"port" envconfig:"DB_PORT"`
	Username string `yaml:"username" envconfig:"DB_USERNAME"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	Migrate  bool   `yaml:"migrate" envconfig:"DB_MIGRATE"`
}

func (x *Config) Init() {
	x.Database = "harvester"
	x.Host = "localhost"
	x.Password = ""
	x.Username = "root"
	x.Port = 3306
}
