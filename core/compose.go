package core

import (
	"fmt"
	"github.com/Doublemine/komposer/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

func Compose(paths []string, isForce bool, isSpecial bool) {
	var configList []model.Config
	for _, path := range paths {
		configList = append(configList, parse2Config(path))
	}
	conbin := separator(configList)
	fmt.Println(len(conbin))
}

// parse to kube-config by file path.
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

func separator(configList []model.Config) []model.MidConfigWare {
	allMidConfWare := make([]model.MidConfigWare, 0)
	for _, config := range configList {
		if len(config.Clusters) != len(config.Contexts) {
			log.Fatalln("only support the same num config.")
		}
		allMidConfWare = append(allMidConfWare, filter2MidWare(config)...)
	}
	return filterDuplicates(allMidConfWare)
}

func filterDuplicates(mayDupWare []model.MidConfigWare) []model.MidConfigWare {
	if len(mayDupWare) <= 1 {
		return mayDupWare
	}
	noduplicate := make([]model.MidConfigWare, 0)
	for _, item := range mayDupWare {
		if len(noduplicate) == 0 {
			noduplicate = append(noduplicate, item)
		} else {
			for _index, _item := range noduplicate {
				if reflect.DeepEqual(item, _item) {
					break
				}
				if _index == len(noduplicate)-1 {
					noduplicate = append(noduplicate, item)
				}
			}
		}
	}

	return noduplicate
}

func filter2MidWare(config model.Config) []model.MidConfigWare {
	clusterMap := map2Cluster(config)
	contextMap := map2Context(config)
	userMap := map2User(config)

	midWareSlice := make([]model.MidConfigWare, 0)
	//may lost context name.
	for _, v := range contextMap {

		temp := model.MidConfigWare{
			Context: v,
		}

		if cluster, ok := clusterMap[v.Context.Cluster]; ok {
			temp.Cluster = cluster
		} else {
			log.Debugf("the cluster not exist!")
		}

		if user, ok := userMap[v.Context.User]; ok {
			temp.User = user
		} else {
			log.Debugf("the user not exist!")
		}

		if temp.Valid() {
			midWareSlice = append(midWareSlice, temp)
		}
	}
	return midWareSlice
}

// map to dict, cluster server url as key
func map2Cluster(clusters model.Config) map[string]model.Clusters {
	clusterMap := make(map[string]model.Clusters)
	for _, cluster := range clusters.Clusters {
		if cluster.Cluster.Server == "" {
			log.Debugf("the server url is empty, skip current clusters.")
			continue
		}
		clusterMap[cluster.Name] = cluster
	}
	return clusterMap
}

// map to dict, config name as key
func map2Context(clusters model.Config) map[string]model.Contexts {

	clusterMap := make(map[string]model.Contexts)
	for _, context := range clusters.Contexts {
		if context.Name == "" {
			log.Debugf("the context name is empty, skip current context.")
			continue
		}
		clusterMap[context.Name] = context
	}
	return clusterMap
}

// map to dict, user name as key
func map2User(clusters model.Config) map[string]model.Users {

	clusterMap := make(map[string]model.Users)
	for _, user := range clusters.Users {
		if user.Name == "" {
			log.Debugf("the user name is empty, skip current context.")
			continue
		}
		clusterMap[user.Name] = user
	}
	return clusterMap
}
