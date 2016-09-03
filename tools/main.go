package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

//kubectl get po -o json --namespace kube-system|grep \"image\"|uniq
var images = []string{
	"gcr.io/google_containers/elasticsearch:1.8",
	"gcr.io/google_containers/fluentd-elasticsearch:1.19",
	"gcr.io/google_containers/heapster:v1.0.2",
	"gcr.io/google_containers/addon-resizer:1.0",
	"gcr.io/google_containers/heapster:v1.0.2",
	"gcr.io/google_containers/addon-resizer:1.0",
	"gcr.io/google_containers/heapster:v1.0.2",
	"gcr.io/google_containers/addon-resizer:1.0",
	"gcr.io/google_containers/kibana:1.3",
	"gcr.io/google_containers/etcd-amd64:2.2.1",
	"gcr.io/google_containers/kube2sky:1.14",
	"gcr.io/google_containers/skydns:2015-10-13-8c72f8c",
	"gcr.io/google_containers/exechealthz:1.0",
	"gcr.io/google_containers/etcd-amd64:2.2.1",
	"gcr.io/google_containers/exechealthz:1.0",
	"gcr.io/google_containers/kube2sky:1.14",
	"gcr.io/google_containers/skydns:2015-10-13-8c72f8c",
	"gcr.io/google_containers/kubedash:v0.2.1",
	"gcr.io/google_containers/heapster_influxdb:v0.5",
	"gcr.io/google_containers/heapster_grafana:v2.6.0-2",
	"gcr.io/google_containers/heapster_influxdb:v0.5",
	"gcr.io/google_containers/pause:2.0",
}

func main() {

	fmt.Printf("images: %d\n", len(images))
	regular := fmt.Sprintf("gcr.io/google_containers/(\\w+):(.+)")
	re := regexp.MustCompile(regular)
	names := []string{}
	nameMap := map[string]string{}

	for _, image := range images {
		match := re.FindStringSubmatch(image)

		if len(match) > 0 {
			//log.Infof("image=%v", image)
			//log.Infof("match=%v", match)
			name := match[1]
			fmt.Println(name)
			if _, found := nameMap[name]; !found {
				names = append(names, name)
				nameMap[name] = image
			}
		}

	}

	fmt.Printf("names: %d\n", len(names))

	root := "/docker-ansible/google-containers"

	for k, v := range nameMap {
		dir := filepath.Join(root, k)
		if err := os.MkdirAll(dir, 0744); err != nil {
			fmt.Println(err)
			return
		}

		file := filepath.Join(dir, "Dockerfile")
		if err := ioutil.WriteFile(file, []byte(fmt.Sprintf("FROM %s\n", v)), 0644); err != nil {
			fmt.Println(err)
			return
		}

	}
}
