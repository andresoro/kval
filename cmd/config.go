package cmd

import "time"

// Config is basic config settings for a kval server
type Config struct {
	shardNum int
	duration time.Duration
	httpPort string
	rpcPort  string
}

// DefaultConfig returns a configuration with default usable values
func DefaultConfig() Config {
	return Config{
		shardNum: 4,
		duration: time.Minute,
		httpPort: ":8080",
		rpcPort:  ":7741",
	}
}
