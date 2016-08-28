package ansible

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	docker_client "github.com/docker/engine-api/client"
	docker "github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

const ansibleEtc = "/etc/ansible"

var dockerClient *docker_client.Client

var (
	ansibleHostFile   = filepath.Join(ansibleEtc, "hosts")
	ansibleSSHkeysDir = filepath.Join(ansibleEtc, "machines")
)

type Deployer interface {
	SetHosts(hosts []types.Host)
	SetCommander(run interface{}) error
	Run() error
}

type ansibleManager struct {
	name                  string
	hosts                 []types.Host
	containerCreateConfig *docker.ContainerCreateConfig
	run                   []string
	roleMap               map[string][]types.Host
	store                 persist.Store
}

func NewDeployer(name string) (Deployer, error) {
	storePath, err := util.GetStorePath(name)
	if err != nil {
		return nil, err
	}

	return &ansibleManager{
		name:  name,
		store: persist.NewFilestore(storePath),
	}, nil
}

func (this ansibleManager) Run() error {

	inventoryFile, err := this.createInventoryFile()
	if err != nil {
		return err
	}
	logrus.Infof("inventory file: %s\n", inventoryFile)

	if this.containerCreateConfig != nil {
		err := this.dockerRun()
		if err != nil {
			return err
		}
	} else {

	}

	return nil
}

func (this *ansibleManager) SetCommander(cmd interface{}) error {
	switch cmd.(type) {
	case []string:
		this.run = cmd.([]string)
	case *docker.ContainerCreateConfig:
		this.containerCreateConfig = cmd.(*docker.ContainerCreateConfig)
	default:
		return fmt.Errorf("Unrecongized type %v", cmd)
	}
	return nil
}

func (this *ansibleManager) SetHosts(hosts []types.Host) {

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

func (this *ansibleManager) dockerRun() error {
	ctx := context.Background()
	dockerClient, err := docker_client.NewEnvClient()
	if err != nil {
		return err
	}

	config := this.containerCreateConfig.Config
	hostConfig := this.containerCreateConfig.HostConfig
	newtworkConfig := this.containerCreateConfig.NetworkingConfig
	hostConfig.Binds = append(hostConfig.Binds, this.genBindsForAnsible()...)
	config.Env = append(config.Env, this.genEnvsForAnsible()...)
	resp, err := dockerClient.ContainerCreate(ctx, config, hostConfig, newtworkConfig, "")
	if err != nil {
		return err
	}
	for _, w := range resp.Warnings {
		logrus.Warnf("Docker create: %v", w)
	}

	id := resp.ID
	options := docker.ContainerStartOptions{}
	err = dockerClient.ContainerStart(ctx, id, options)
	if err != nil {
		return err
	}
	return nil
}

// create the inventory file which is used by ansible
func (this *ansibleManager) createInventoryFile() (path string, err error) {
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
	path = filepath.Join(deploymentDir, "inventory")
	f, err := os.Create(path)
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
				this.mappingKeyPath(h.SSHKeyPath)))
			if err != nil {
				return path, err
			}
		}

		_, err = w.WriteString("\n")
		if err != nil {
			return path, err
		}
	}

	if this.name != "" {
		_, err = w.WriteString(fmt.Sprintf("[%s:children]\n", this.name))
		if err != nil {
			return path, err
		}

		for k, _ := range this.roleMap {
			_, err := w.WriteString(fmt.Sprintf("%s\n", k))
			if err != nil {
				return path, err
			}
		}
	}

	return path, err
}

func (this *ansibleManager) genBindsForAnsible() (binds []string) {

	binds = append(binds,
		fmt.Sprintf("%s:%s:ro", filepath.Join(this.store.GetDeploymentDir(), "inventory"), ansibleHostFile),
		fmt.Sprintf("%s:%s:ro", filepath.Join(this.store.GetMachinesDir(), ansibleSSHkeysDir)),
	)

	return binds
}

func (this *ansibleManager) genEnvsForAnsible() []string {
	return []string{
		"ANSIBLE_HOST_KEY_CHECKING=False",
	}
}

func (this *ansibleManager) mappingKeyPath(keyPath string) string {
	if this.containerCreateConfig != nil {
		return strings.Replace(keyPath, this.store.GetMachinesDir(), ansibleSSHkeysDir, 1)
	}
	return keyPath
}
