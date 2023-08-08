package internal

import (
	"flag"
	"log"
)

type Config struct {
	Env      string
	FilePath string
	Target   string
}

func (c *Config) Init() {
	flag.StringVar(&c.Env, "env", "", "an string")
	flag.StringVar(&c.FilePath, "file", "", "an string")
	flag.StringVar(&c.Target, "target", "", "an string")
	flag.Parse()
	if c.FilePath == "" {
		log.Fatal("Не указано имя файла")
	}
	if c.Env == "" {
		log.Fatal("Не указано окружение")
	}
	if c.Target != "firestate" && c.Target != "backup" {
		log.Fatal("Не указана цель", c.Target)
	}
}
