package pkg

// CduleConfig cdule configuration
type CduleConfig struct {
	Cdule    cdule    `yaml:"cdule"`
	Database database `yaml:"database"`
}

type cdule struct {
	Type        string `yaml:"type"`
	Consistency string `yaml:"consistency"`
}

type database struct {
	Type               string `yaml:"type"`
	Url                string `yaml:"url"`
	MaxIdleConnections int    `yaml:"maxIdleConnections"`
	MaxOpenConnections int    `yaml:"maxOpenConnections"`
}
