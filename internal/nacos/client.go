package nacos

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var NamingClient naming_client.INamingClient

func Init() error {
	_ = godotenv.Load()
	// godotenv.oad()

	host := os.Getenv("NACOS_HOST")
	portStr := os.Getenv("NACOS_PORT")
	user := os.Getenv("NACOS_USER")     // Nacos username
	pass := os.Getenv("NACOS_PASSWORD") // Nacos password
	namespaceId := os.Getenv("NACOS_NAMESPACE_ID")

	if host == "" || user == "" || pass == "" || portStr == ""{
		return fmt.Errorf("Nacos host, username or password not set")
	}

	// 👇 Split multiple hosts
	hostList := strings.Split(host, ",")

	var serverConfigs []constant.ServerConfig

	for _, h := range hostList {
		h = strings.TrimSpace(h)

		// Expect format: ip:port
		parts := strings.Split(h, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid NACOS_HOST format: %s", h)
		}

		ip := parts[0]
		p, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return err
		}

		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: ip,
			Port:   p,
		})
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "info",
		Username:            user, // <-- username
		Password:            pass, // <-- password
	}
   var err error
	NamingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: serverConfigs,
			ClientConfig:  &clientConfig,
		},
	)
	return err
}
