package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func PrintBackupBlock(group []string, fs *linkedhashmap.Map, env string) {
	for _, line := range fs.Keys() {
		var variable strings.Builder
		for i := 1; i < len(group); i++ {
			variable.WriteString(fmt.Sprintf("    - ansible-playbook ansible/deploy.yaml --connection=local -t %s -e K8S_TOKEN=$K8S_ORDER_%s -e TARGET=backup/%s/backup%s.yaml\n", line, strings.ToUpper(line.(string)), line.(string), group[i]))
		}
		fmt.Printf("%s_BACKUP_%s_%s.deploy:\n  stage: deploy\n  image:\n    name: docker.hub:1.0.0\n    entrypoint: [ \"\" ]\n  tags:\n    - linux-docker-executor\n  when: manual\n  script:\n%s\n  needs:\n    - pipeline: $PARENT_PIPELINE_ID\n      job: %s_BACKUP.manifests\n", strings.ToUpper(string(env)), group[0], line, variable.String(), strings.ToUpper(string(env)))
	}
}

func PrintFsBlock(group []string, fs *linkedhashmap.Map, env string) {
	var variable strings.Builder
	for i := 1; i < len(group); i++ {
		line, t := fs.Get(group[i])
		if t == true {
			variable.WriteString(fmt.Sprintf("    - ansible-playbook ansible/deploy.yaml --connection=local -t %s -e K8S_TOKEN=$K8S_ORDER_%s -e TARGET=firestate/%s/firestate%s.yaml\n", line, strings.ToUpper(line.(string)), line, group[i]))
		} else {
			log.Fatal("Для элемента \"", group[i], "\" не указан кластер")
		}
	}
	fmt.Printf("%s_FIRESTATE_%s.deploy:\n  stage: deploy\n  image:\n    name: docker.hub:1.0.0\n    entrypoint: [ \"\" ]\n  tags:\n    - linux-docker-executor\n  when: manual\n  script:\n%s\n  needs:\n    - pipeline: $PARENT_PIPELINE_ID\n      job: %s_FIRESTATE.manifests\n", strings.ToUpper(string(env)), group[0], variable.String(), strings.ToUpper(string(env)))
}
