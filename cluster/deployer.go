package cluster

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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

	for k, hosts := range this.roleMap {
		fmt.Printf("[%s]\n", k)

		for _, h := range hosts {
			fmt.Printf("%s ansible_host=%s ansible_user=%s ansible_ssh_private_key_file=%s",
				h.Name,
				h.SSHHostname,
				h.SSHUserName,
				h.SSHKeyPath)
		}

		fmt.Println("")
	}

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
}

func (this *ansibleDeployer) createInventoryFile() (path string, err error) {
	storePath, err := util.GetStorePath(this.name)
	if err != nil {
		return
	}
	storage := persist.NewFilestore(storePath)
	err = storage.CreateDeploymentDir()
	if err != nil {
		return
	}

	deploymentDir := storage.GetDeploymentDir()
	f, err := os.Create(filepath.Join(deploymentDir, "inventory"))
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for k, hosts := range this.roleMap {
		_, err = w.WriteString(fmt.Sprintf("[%s]\n", k))
		if err != nil {
			return
		}

		for _, h := range hosts {
			_, err = w.WriteString(fmt.Sprintf("%s ansible_host=%s ansible_user=%s ansible_ssh_private_key_file=%s\n",
				h.Name,
				h.SSHHostname,
				h.SSHUserName,
				h.SSHKeyPath))
			if err != nil {
				return
			}
		}

		_, err = w.WriteString("\n")
		if err != nil {
			return
		}
	}

	return
}
