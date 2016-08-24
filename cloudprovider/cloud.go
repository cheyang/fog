package cloudprovider

import (
	"github.com/Sirupsen/logrus"
	aliyun_k8s "github.com/cheyang/cloudprovider/aliyun/k8s"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type CloudInterface interface {
	setConfigFromFlags(opts drivers.DriverOptions) error

	SetHosts(hosts []types.Host)

	Configure() error
}

func InitProivder(provider, clusterType string) CloudInterface {

	if providerFunc[provider][clusterType] == nil {
		logrus.Infof("Not able to find provider %s for %s", provider, clusterType)
	}

	return providerFunc[provider][clusterType]
}

var providerFunc = map[string](map[string]func() CloudInterface){
	"aliyun": map[string]func() CloudInterface{
		"k8s": aliyun_k8s.New(),
	},
}
