package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func Resolve(serviceName string) (string, error) {
	instances, err := NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
		HealthyOnly: true,
	})
	if err != nil || len(instances) == 0 {
		return "", err
	}

	instance := instances[0]
	return instance.Ip + ":" + fmt.Sprint(instance.Port), nil
}
