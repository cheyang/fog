package cluster

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
)

type Deployer interface {
	SetHosts(hosts []types.Host)
	Run() error
}

type ansibleDeployer struct {
	name    string
	hosts   []types.Host
	roleMap map[string][]types.Host
}

func (this ansibleDeployer) Run() error {

	inventoryFile, err := this.createInventoryFile()
	if err != nil {
		return err
	}
	logrus.Infof("inventory file: %s\n", inventoryFile)

	return nil
}

func (this *ansibleDeployer) SetHosts(hosts []types.Host) {

	this.hosts = hosts
	this.roleMap = make(map[string][]types.Host)

	for _, host := range hosts {

		for _, role := range host.Roles {

			if _, found := this.roleMap[role]; !found {
				this.roleMap[role] = make([]types.Host, 0)
			}

			this.roleMap[role] = append(this.roleMap[role], host)
		}

	}

	for _, value := range this.roleMap {
		sort.Sort(byHostName(value))
	}
}

func (this *ansibleDeployer) createInventoryFile() (path string, err error) {
	storePath, err := util.GetStorePath(this.name)
	if err != nil {
		return path, err
	}
	storage := persist.NewFilestore(storePath)
	err = storage.CreateDeploymentDir()
	if err != nil {
		return path, err
	}

	deploymentDir := storage.GetDeploymentDir()
	f, err := os.Create(filepath.Join(deploymentDir, "inventory"))
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for k, hosts := range this.roleMap {
		_, err = w.WriteString(fmt.Sprintf("[%s]\n", k))
		if err != nil {
			return path, err
		}

		for _, h := range hosts {
			_, err = w.WriteString(fmt.Sprintf("%s ansible_host=%s ansible_user=%s ansible_ssh_private_key_file=%s\n",
				h.Name,
				h.SSHHostname,
				h.SSHUserName,
				h.SSHKeyPath))
			if err != nil {
				return path, err
			}
		}

		_, err = w.WriteString("\n")
		if err != nil {
			return path, err
		}

		if this.name != "" {
			_, err = w.WriteString(fmt.Sprintf("[%s:children]\n", this.name))
			if err != nil {
				return path, err
			}

			for k, _ := range this.roleMap {
				_, err := w.WriteString(k)
				if err != nil {
					return path, err
				}
			}
		}
	}

	return path, err
}

type byHostName []types.Host

func (s byHostName) Len() int {
	return len(s)
}
func (s byHostName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byHostName) Less(i, j int) bool {
	sai := strings.Split(s[i].Name, "-")
	si, err := strconv.Atoi(sai[len(sai)-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	saj := strings.Split(s[j].Name, "-")
	sj, err := strconv.Atoi(saj[len(saj)-1])
	if err != nil {
		logrus.Infof("err: %v", err)
	}
	return si < sj
}
