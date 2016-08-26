package persist

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/mcnerror"
)

type Filestore struct {
	Path string
}

func NewFilestore(path string) *Filestore {
	return &Filestore{
		Path: path,
	}
}

func (s Filestore) GetMachinesDir() string {
	return filepath.Join(s.Path, "machines")
}

func (s Filestore) CreateStorePath(name string) error {
	hostPath := filepath.Join(s.GetMachinesDir(), name)

	// Ensure that the directory we want to save to exists.
	if err := os.MkdirAll(hostPath, 0700); err != nil {
		return err
	}

	return nil
}

func (s Filestore) Save(host *types.Host) error {
	data, err := json.MarshalIndent(host, "", "    ")
	if err != nil {
		return err
	}

	hostPath := filepath.Join(s.GetMachinesDir(), host.Name)

	// Ensure that the directory we want to save to exists.
	if err := os.MkdirAll(hostPath, 0700); err != nil {
		return err
	}

	logrus.Debugf("config.json: %s", filepath.Join(hostPath, "config.json"))

	return s.saveToFile(data, filepath.Join(hostPath, "config.json"))
}

func (s Filestore) saveToFile(data []byte, file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return ioutil.WriteFile(file, data, 0600)
	}

	tmpfi, err := ioutil.TempFile(filepath.Dir(file), "config.json.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfi.Name())

	if err = ioutil.WriteFile(tmpfi.Name(), data, 0600); err != nil {
		return err
	}

	if err = tmpfi.Close(); err != nil {
		return err
	}

	if err = os.Remove(file); err != nil {
		return err
	}

	if err = os.Rename(tmpfi.Name(), file); err != nil {
		return err
	}
	return nil
}

func (s Filestore) List() ([]string, error) {
	dir, err := ioutil.ReadDir(s.GetMachinesDir())
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	hostNames := []string{}

	for _, file := range dir {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			hostNames = append(hostNames, file.Name())
		}
	}

	return hostNames, nil
}

func (s Filestore) Exists(name string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.GetMachinesDir(), name))

	if os.IsNotExist(err) {
		return false, nil
	} else if err == nil {
		return true, nil
	}

	return false, err
}

func (s Filestore) loadConfig(h *types.Host) error {
	data, err := ioutil.ReadFile(filepath.Join(s.GetMachinesDir(), h.Name, "config.json"))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, h)

}

func (s Filestore) Load(name string) (*types.Host, error) {
	hostPath := filepath.Join(s.GetMachinesDir(), name)

	if _, err := os.Stat(hostPath); os.IsNotExist(err) {
		return nil, mcnerror.ErrHostDoesNotExist{
			Name: name,
		}
	}

	host := &types.Host{
		Name: name,
	}

	if err := s.loadConfig(host); err != nil {
		return nil, err
	}

	return host, nil
}
