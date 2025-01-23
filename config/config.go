package config

type Config struct {
	CacheDirectory  string
	ImportDirectory string
	Log             Log
	Theme           Theme
}

type Log struct {
	Debug     bool
	Directory string
}

type Theme struct {
	Name string
}

func New() *Config {
	return &Config{}
}
