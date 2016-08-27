package cluster

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/persist"
)

func Remove(s persist.Store) error {

	hostList, _, err := persist.LoadAllHosts(s)
	if err != nil {
		return err
	}

	for _, host := range hostList {
		err := host.Driver.Remove()
		s.Remove(host.Name)
		if err != nil {
			logrus.Infof("host err: %v", err)
		}
	}

	return nil
}
