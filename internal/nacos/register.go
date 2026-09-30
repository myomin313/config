package nacos

import (
	"log"
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func RegisterService(
	serviceName string,
	ip string,
	port uint64,
) error {

	_, err := NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		GroupName:   os.Getenv("NACOS_GROUP"),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true, // auto remove when service stops

	})
	return err
}

func DeregisterService() {
	appName := os.Getenv("APP_NAME")
	serviceIP := os.Getenv("SERVICE_IP")
	appPort := uint64(mustAtoi(os.Getenv("APP_PORT")))
	groupName := os.Getenv("NACOS_GROUP")

	success, err := NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: appName,
		Ip:          serviceIP,
		Port:        appPort,
		GroupName:   groupName,
	})

	if err != nil {
		log.Printf("⚠️ Failed to deregister service: %v\n", err)
	} else if success {
		log.Println("🛑 Service deregistered from Nacos successfully")
	} else {
		log.Println("⚠️ Service deregister response was false")
	}
}

func mustAtoi(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("❌ Invalid integer: %s", val)
		return 0
	}
	return i
}
