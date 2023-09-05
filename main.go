package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"strings"
	"test/internal"
)

func readShard(filePathShard string, sep string) []string {
	file, err := os.ReadFile(filePathShard)
	if err != nil {
		log.Fatal("Ошибка чтения файла:", err)
	}
	strSplit := strings.Split(string(file), sep)
	return strSplit
}

func createManifests(f *internal.Flag) {
	file := readShard(f.ShardValuesFile, " ")
	for _, v := range file {
		name := fmt.Sprint("firestate", string(v))
		command := fmt.Sprint("select(.metadata.name == \"", name, "\" , .metadata.name == \"firestate-config\")")
		cmd := exec.Command("yq", command, f.ManifestsFile)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(string(err.(*exec.ExitError).Stderr))
			fmt.Println(err)
			return
		}
		name = fmt.Sprint(name, ".yaml")
		var outputFile string
		if f.OutFile != "" {
			outputFile = fmt.Sprint(f.OutFile, "/", name)
		} else {
			outputFile = name
		}
		err = os.WriteFile(outputFile, out, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func kustomizeUse(deployment *internal.Deployment, config *internal.Config, f *internal.Flag) {

	if f.OverrideFile == "" {
		return
	}
	configDeployment := readConfig(f.OverrideFile, config)
	for key, firestate := range configDeployment.Firestates {
		var outputFile string
		name := fmt.Sprint(key, ".yaml")
		if f.OutFile != "" {
			outputFile = fmt.Sprint(f.OutFile, "/", name)
		} else {
			outputFile = name
		}

		data := readShard(outputFile, "\n---\n")
		err := yaml.Unmarshal([]byte(data[2]), &deployment)
		if err != nil {
			log.Fatalf("cannot unmarshal data: %v", err)
		}
		originResources := &deployment.Spec.Template.Spec.Containers[0].Resources
		newResources := firestate.Resources
		if newResources.Limits.CPU != "" {
			originResources.Limits.CPU = newResources.Limits.CPU
		}
		if newResources.Limits.Memory != "" {
			originResources.Limits.Memory = newResources.Limits.Memory
		}
		if newResources.Requests.CPU != "" {
			originResources.Requests.CPU = newResources.Requests.CPU
		}
		if newResources.Requests.Memory != "" {
			originResources.Requests.Memory = newResources.Requests.Memory
		}
		dataResult, err := yaml.Marshal(deployment)
		if err != nil {
			log.Fatal(err)
		}
		resultManifest := fmt.Sprint(data[0], "\n---\n", data[1], "\n---\n", string(dataResult))

		err = os.WriteFile(outputFile, []byte(resultManifest), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readConfig(filename string, deployment *internal.Config) *internal.Config {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &deployment)
	if err != nil {
		log.Fatalf("Ошибка разбора YAML: %v", err)
	}
	return deployment
}

func main() {
	deployment := new(internal.Deployment)
	config := new(internal.Config)
	flag := new(internal.Flag)
	flag.Init()
	createManifests(flag)
	kustomizeUse(deployment, config, flag)
}
