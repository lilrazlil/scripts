package internal

import (
	"flag"
	"log"
)

type Flag struct {
	ShardValuesFile string
	ManifestsFile   string
	OverrideFile    string
	OutFile         string
}

func (Flag *Flag) Init() {
	flag.StringVar(&Flag.ShardValuesFile, "shard-values", "", "")
	flag.StringVar(&Flag.ManifestsFile, "manifest", "", "")
	flag.StringVar(&Flag.OverrideFile, "override", "", "")
	flag.StringVar(&Flag.OutFile, "dir", "", "")
	flag.Parse()
	if Flag.ManifestsFile == "" {
		log.Fatal("Не указан файл с манифестами")
	}
	if Flag.ShardValuesFile == "" {
		log.Fatal("Не указан файл шард")
	}
}
