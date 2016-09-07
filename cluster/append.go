package cluster

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cluster/ansible"
	"github.com/cheyang/fog/cluster/deploy"
	"github.com/cheyang/fog/host"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
)

// key is the vmspec name, value is the host name list
var (
	runningHostMap map[string][]string

	splitHostname = "(.+)-(\\d+)"
)

func ExpandCluster(s persist.Store, spec types.Spec, requiredRoles []string) error {
	spec.Update = true
	hosts, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	err = buildRunningMap(hosts)
	if err != nil {
		return err
	}
	for i, vmSpec := range spec.VMSpecs {
		spec.VMSpecs[i].Start, err = nextNumber(vmSpec.Name)
		if err != nil {
			return err
		}
	}
	vmSpecs, err := host.BuildHostConfigs(spec)
	if err != nil {
		return err
	}

	newHosts, err := provisionVMs(spec)
	if err != nil {
		return err
	}
	cp := provider_registry.GetProvider(spec.CloudDriverName, spec.ClusterType)
	if cp != nil {
		cp.SetHosts(newHosts)
		cp.Configure() // configure IaaS
	}
	var deployer deploy.Deployer
	deployer, err = ansible.NewDeployer(spec.Name)
	if err != nil {
		return err
	}
	deployer.SetHosts(newHosts)
	if len(spec.Run) > 0 {
		deployer.SetCommander(spec.Run)
	} else {
		deployer.SetCommander(spec.DockerRun)
	}

	return deployer.Run()
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
		logrus.Infof("The max of %s is %d", name, max)
		return max + 1, nil
	}
	return 0, nil
}

// parse host name to two parts, spec name and id, take master-1 as example, spec name is master, id is 1
func parseHostname(hostname string) (specName string, id int, err error) {
	re := regexp.MustCompile(splitHostname)
	match := re.FindStringSubmatch(hostname)
	specName = match[1]
	id, err = strconv.Atoi(match[2])
	return specName, id, err
}

func buildRunningMap(hosts []*types.Host) (err error) {
	runningHostMap = make(map[string][]string)

	for _, host := range hosts {
		// build running host map
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
	return nil
}

type ByHostname []string

func (s ByHostname) Len() int {
	return len(s)
}
func (s ByHostname) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByHostname) Less(i, j int) bool {
	_, si, err := parseHostname(s[i])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	_, sj, err := parseHostname(s[j])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	return si < sj
}
