package core

import (
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

	for _, config := range configList {
		spearator(config)
	}

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

	return eachCluster
}

func filter2MidWare(config model.Config) []model.MidConfigWare {
	clusterMap := map2Cluster(config)
	contextMap := map2Context(config)
	userMap := map2User(config)

	log.Println("cluster:", clusterMap, "\ncontextMap:", contextMap, "\nuserMap", userMap)

	var midWareMap map[string]model.MidConfigWare
	for k, v := range contextMap {

		temp := model.MidConfigWare{
			Context: v,
		}
		cluster, ok := clusterMap[v.Context.Cluster]
		if ok {
			temp.Cluster = cluster
		}

		user, ok := userMap[v.Context.User]
		if ok {
			temp.User = user
		}
		if cluster != nil && cluster.Cluster.Server != "" {
			midWareMap[cluster.Cluster.Server] = temp
		}

	}
}

func map2Cluster(clusters model.Config) map[string]model.Clusters {
	clusterMap := make(map[string]model.Clusters)
	for _, cluster := range clusters.Clusters {
		if cluster.Cluster.Server == "" {
			log.Println("the server url is empty, skip current clusters.")
			continue
		}
		clusterMap[cluster.Cluster.Server] = cluster
	}
	return clusterMap
}

func map2Context(clusters model.Config) map[string]model.Contexts {

	clusterMap := make(map[string]model.Contexts)
	for _, context := range clusters.Contexts {
		if context.Name == "" {
			log.Println("the context name is empty, skip current context.")
			continue
		}
		clusterMap[context.Name] = context
	}
	return clusterMap
}

func map2User(clusters model.Config) map[string]model.Users {

	clusterMap := make(map[string]model.Users)
	for _, user := range clusters.Users {
		if user.Name == "" {
			log.Println("the user name is empty, skip current context.")
			continue
		}
		clusterMap[user.Name] = user
	}
	return clusterMap
}
