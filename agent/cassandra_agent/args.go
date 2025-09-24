package main

import "flag"

type Args struct {
	configPath string
}

func NewArgs() Args {
	data := flag.String("configPath","../conf/config.toml","config file path")
	flag.Parse()
	
	return Args{
		configPath: *data,
	}
}