package persist

import "github.com/cheyang/fog/host"

type Store interface {
	// Exists returns whether a machine exists or not
	Exists(name string) (bool, error)

	// List returns a list of all hosts in the store
	List() ([]string, error)

	// Load loads a host by name
	Load(name string) (*host.HostHandler, error)

	// Remove removes a machine from the store
	Remove(name string) error

	// Save persists a machine in the store
	Save(host *host.HostHandler) error
}

func LoadHosts(s Store, hostNames []string) ([]*host.HostHandler, map[string]error) {
	loadedHosts := []*host.HostHandler{}
	errors := map[string]error{}

	for _, hostName := range hostNames {
		h, err := s.Load(hostName)
		if err != nil {
			errors[hostName] = err
		} else {
			loadedHosts = append(loadedHosts, h)
		}
	}

	return loadedHosts, errors
}

func LoadAllHosts(s Store) ([]*host.HostHandler, map[string]error, error) {
	hostNames, err := s.List()
	if err != nil {
		return nil, nil, err
	}
	loadedHosts, hostInError := LoadHosts(s, hostNames)
	return loadedHosts, hostInError, nil
}
