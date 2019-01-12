package core

import (
	"fmt"
	"github.com/Doublemine/komposer/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func Compose(paths []string, isForce bool, isSpecial bool) {
	var configList []model.Config
	for _, path := range paths {
		configList = append(configList, parse2Config(path))
	}
	fmt.Println(configList)

}

func parse2Config(path string) model.Config {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	auth := model.Config{}
	err = yaml.Unmarshal(content, &auth)
	if err != nil {
		log.Fatal("the file: " + path + " can not resolve")
	}
	return auth
}

func spearator(config model.Config) map[string]model.Config {
	eachCluster := make(map[string]model.Config)

	if len(config.Clusters) != len(config.Contexts) {
		log.Fatalln("only support the same num config.")
	}

	//for index, cluster := range config.Clusters {
	//
	//
	//}

	return eachCluster
}
