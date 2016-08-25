package cluster

import (
	"fmt"

	"github.com/cheyang/fog/types"
)

type Deployer interface {
	SetHosts(hosts []types.Host)
	Run() error
}

type ansibleDeployer struct {
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
