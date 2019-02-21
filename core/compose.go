package core

import (
	"github.com/Doublemine/komposer/model"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

var (
	//auto append suffix
	duplicateSuffix = "-k6r"
	match           = regexp.MustCompile("(" + duplicateSuffix + ")+")
)

const (
	DEFAULT_KUBECONFIG_NAME  = "cluster.kubeconfig"
	USER_DIR_KUBECONFIG_NAME = "config"
	KUBE_DIR                 = ".kube"
)

func Compose(paths []string, isForce bool, isSpecial bool, suffix string) {
	if len(suffix) > 0 {
		duplicateSuffix = suffix
	}
	var configList []model.Config
	for _, path := range paths {
		configList = append(configList, parse2Config(path))
	}
	conbin := separator(configList)
	config := merge2Config(conbin)
	writeConfig2File(config, isSpecial, isForce)

}

func merge2Config(data []model.MidConfigWare) model.Config {
	renameWare := make([]model.MidConfigWare, 0)
	for _, item := range data {
		if len(renameWare) <= 0 {
			renameWare = append(renameWare, item)
		} else {
			for _, _item := range renameWare {
				renameDuplicateName(_item, &item)
			}
			renameWare = append(renameWare, item)
		}
	}
	contexts := make([]model.Contexts, 0)
	clusters := make([]model.Clusters, 0)
	users := make([]model.Users, 0)

	for _, item := range renameWare {
		contexts = append(contexts, item.Context)
		clusters = append(clusters, item.Cluster)
		users = append(users, item.User)
	}

	return model.Config{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: contexts[0].Name,
		Preferences:    model.Preferences{},
		Users:          users,
		Clusters:       clusters,
		Contexts:       contexts,
	}
}

func renameDuplicateName(name model.MidConfigWare, todo *model.MidConfigWare) {
	if name.Cluster.Name == todo.Cluster.Name {
		todo.Cluster.Name += duplicateSuffix
		todo.Context.Context.Cluster += duplicateSuffix
	}

	if name.User.Name == todo.User.Name {
		todo.User.Name += duplicateSuffix
		todo.Context.Context.User += duplicateSuffix
	}

	if name.Context.Name == todo.Context.Name {
		todo.Context.Name += duplicateSuffix
	}

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

// parse to kube-config by file path.
func writeConfig2File(config model.Config, isSpecial bool, isForce bool) {

	var path string
	if isSpecial {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatalln("can not found homedir, error:", err)
		}
		path = filepath.Join(dir, KUBE_DIR, USER_DIR_KUBECONFIG_NAME)

		_, _err := os.Stat(filepath.Join(dir, KUBE_DIR))
		if _err != nil && os.IsNotExist(_err) {
			_ = os.MkdirAll(filepath.Join(dir, KUBE_DIR), os.ModePerm)
		}
	} else {
		path = DEFAULT_KUBECONFIG_NAME
	}

	if !isForce && FileExist(path) {
		log.Infof("the config already exists, you can use flag: --force to overwrite it.")
		os.Exit(0)
	}

	content, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		log.Fatalf("write config to %s failed!", path)
	}
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

	noSuffix := make([]model.MidConfigWare, 0)
	for _, item := range mayDupWare {
		item.Cluster.Name = match.ReplaceAllString(item.Cluster.Name, ``)
		item.User.Name = match.ReplaceAllString(item.User.Name, ``)
		item.Context.Name = match.ReplaceAllString(item.Context.Name, ``)
		item.Context.Context.User = match.ReplaceAllString(item.Context.Context.User, ``)
		item.Context.Context.Cluster = match.ReplaceAllString(item.Context.Context.Cluster, ``)
		noSuffix = append(noSuffix, item)
	}
	//cleanup slice
	mayDupWare = mayDupWare[:0]

	noduplicate := make([]model.MidConfigWare, 0)
	for _, item := range noSuffix {
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
