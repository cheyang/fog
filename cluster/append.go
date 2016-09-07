package cluster

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
)

// key is the vmspec name, value is the host name list
var (
	runningHostMap map[string][]string

	splitHostname = "(.+)-(\\d+)"
)

func ExpandCluster(s persist.Store, appendSpec types.Spec, requiredRoles []string) error {

	hostList, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	var spec *types.Spec
	spec, err = s.LoadSpec()
	if err != nil {
		return nil
	}

	return nil
}

// next number of the specified vmspec name
func nextNumber(name string) (int, error) {
	if orderedHostnames, found := runningHostMap[name]; found {
		maxIndex := len(orderedHostnames) - 1
		// s := strings.Split(orderedHostnames[maxIndex], "-")
		// max, err := strconv.Atoi(s[len(s)-1])
		hostname := orderedHostnames[maxIndex]
		_, max, err := parseHostname(hostname)
		if err != nil {
			return 0, err
		}
		return max + 1, nil
	}

	return 0, nil
}

func parseHostname(hostname string) (specName string, id int, err error) {
	re := regexp.MustCompile(splitHostname)
	match := re.FindStringSubmatch(s)
	specName = match[1]
	id, err = strconv.Atoi(match[2])
	if err != nil {
		return "", "", err
	}
	return specName, id, nil
}

func buildRunningHostMap(hosts []*types.Host, err error) {
	runningHostMap = make(map[string][]string)

	for _, host := range hosts {
		hostname := host.Name
		key, _, err := parseHostname(hostname)
		if err != nil {
			return err
		}

		if _, found := runningHostMap[key]; !found {
			runningHostMap[key] = []string{}
		}

		runningHostMap[key] = append(runningHostMap[key], hostname)

	}

	for _, v := range runningHostMap {
		sort.Sort(ByHostname(v))
	}
}

type ByHostname []string

func (s ByHostname) Len() int {
	return len(s)
}
func (s ByHostname) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByHostname) Less(i, j int) bool {
	si, err := strconv.Atoi(strings.Split(s[i], "-")[len(strings.Split(s[i], "-"))-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	sj, err := strconv.Atoi(strings.Split(s[j], "-")[len(strings.Split(s[j], "-"))-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	return si < sj
}
