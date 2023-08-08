package main

import (
	"create-shard-deploy-config/internal"
	"io"
	"log"
	"os"
	"strings"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func parsing(c *internal.Config, fs *linkedhashmap.Map) {
	file, err := os.Open(c.FilePath)
	if err != nil {
		log.Fatal("Ошибка чтения файла:")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Не получилось закрыть файл")
		}
	}(file)
	p, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Ошибка чтения файла:")
	}
	splitStr := strings.Split(string(p), "---\n")
	for _, line := range strings.Split(splitStr[1], "\n") {
		if len(line) == 0 {
			continue
		}
		name := strings.Split(line, ":")
		partsFs := strings.Split(name[1], ",")
		for _, lines := range partsFs {
			if c.Target == "firestate" {
				fs.Put(lines, name[0])
			}
			if c.Target == "backup" {
				fs.Put(name[0], lines)
			}
		}
	}
	for _, line := range strings.Split(splitStr[0], "\n") {
		if len(line) == 0 {
			continue
		}
		name := strings.Split(line, ":")
		partsGroup := strings.Split(name[1], ",")

		row := name[0] + " " + strings.Join(partsGroup, " ")
		group := strings.Split(row, " ")
		if c.Target == "firestate" {
			internal.PrintFsBlock(group, fs, c.Env)
		}
		if c.Target == "backup" {
			internal.PrintBackupBlock(group, fs, c.Env)
		}
	}
}

func main() {
	conf := new(internal.Config)
	conf.Init()
	fs := linkedhashmap.New()
	parsing(conf, fs)
}
