package cloudprovider

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type CloudInterface interface {
	SetConfigFromFlags(opts drivers.DriverOptions) error

	SetHosts(hosts []types.Host)

	Configure() error
}

func initProivder(provider, clusterType string) CloudInterface {

	providerFunc := providerFuncMap[provider][clusterType]

	if providerFunc == nil {
		logrus.Infof("Not able to find provider %s for %s, ignore it...", provider, clusterType)
		return nil
	}

	return providerFunc()
}

func RegisterProvider(cloudDriverName string, clusterType string, method func() CloudInterface) error {

	providerFuncMap[cloudDriverName] = map[string]func() CloudInterface{
		clusterType: method,
	}

	return nil
}

var providerFuncMap = map[string](map[string]func() CloudInterface){
	"aliyun": map[string]func() CloudInterface{
		"k8s": aliyun_k8s.New,
	},
}
