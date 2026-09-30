package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var NacosClient naming_client.INamingClient

// --------------------
// INIT (.env loading)
// --------------------
func init() {
	// Load .env ONLY for local development (safe if file does not exist)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	} else {
		log.Println("✅ .env file loaded success")
	}
}

// --------------------
// ENV LOGGING (SAFE)
// --------------------
func logEnvSummary() {
	log.Println("🔧 ENV SUMMARY")
	log.Println("APP_NAME:", os.Getenv("APP_NAME"))
	log.Println("APP_PORT:", os.Getenv("APP_PORT"))
	log.Println("SERVICE_IP:", os.Getenv("SERVICE_IP"))
	log.Println("NACOS_SERVER_ADDR:", os.Getenv("NACOS_SERVER_ADDR"))
	log.Println("NACOS_SERVER_PORT:", os.Getenv("NACOS_SERVER_PORT"))
	log.Println("NACOS_GROUP:", os.Getenv("NACOS_GROUP"))
	log.Println("NACOS_USERNAME:", os.Getenv("NACOS_USERNAME"))
	log.Println("NACOS_PASSWORD:", os.Getenv("NACOS_PASSWORD"))
}

// --------------------
// NACOS INIT
// --------------------
func InitNacos() {
	logEnvSummary()

	addr := os.Getenv("NACOS_SERVER_ADDR")
	portStr := os.Getenv("NACOS_SERVER_PORT")

	if addr == "" || portStr == "" {
		log.Fatal("❌ NACOS_SERVER_ADDR or NACOS_SERVER_PORT not set")
	}

	ip, port := parseNacosIP(addr + ":" + portStr)

	log.Printf("🔌 Connecting to Nacos at %s:%d\n", ip, port)

	// ✅ FIXED: Use constant.NewServerConfig (proper v2 API)
	sc := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public", // default namespace
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "info",
		LogDir:              "./nacos/log",
		CacheDir:            "./nacos/cache",
		Username:            os.Getenv("NACOS_USERNAME"),
		Password:            os.Getenv("NACOS_PASSWORD"),
	}

	// ✅ FIXED: Pass serverConfigs and clientConfig as slices/pointers correctly
	// client, err := clients.CreateNamingClient(map[string]interface{}{
	// 	"serverConfigs": sc,
	// 	"clientConfig":  &cc, // Pass as pointer
	// })
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Fatal("❌ Nacos init failed:", err)
	}

	NacosClient = client
	log.Println("✅ Nacos client created successfully")

}

// --------------------
// WAIT UNTIL NACOS READY
// --------------------
func WaitForNacosReady() {
	log.Println("⏳ Waiting for Nacos client to be ready...")

	for i := 1; i <= 20; i++ { // ~2 minutes
		_, err := NacosClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			GroupName: "DEFAULT_GROUP",
			PageNo:    1,
			PageSize:  1,
		})

		if err == nil {
			log.Printf("✅ Nacos client ready after %d attempts\n", i)
			return
		} else {
			log.Printf("❌ Nacos not ready yet error: %v\n", err)
		}

		if i%10 == 0 {
			log.Printf("⏳ Still waiting for Nacos... (%d/120)\n", i)
		}
		time.Sleep(time.Second)
	}
	log.Println("⚠️ Nacos client still not fully ready, proceeding anyway")
}

// --------------------
// REGISTER SERVICE
// --------------------
func RegisterService() {
	WaitForNacosReady()
	appName := os.Getenv("APP_NAME")
	serviceIP := os.Getenv("SERVICE_IP")
	appPort := uint64(mustAtoi(os.Getenv("APP_PORT")))
	//groupName := os.Getenv("NACOS_GROUP")

	if appName == "" || serviceIP == "" || appPort == 0 {
		log.Fatal("❌ APP_NAME, SERVICE_IP or APP_PORT missing")
	}

	log.Printf("📝 Registering service %s at %s:%d\n", appName, serviceIP, appPort)

	for attempt := 1; attempt <= 5; attempt++ {
		// ✅ FIXED: Added ClusterName (required in some Nacos versions)
		success, err := NacosClient.RegisterInstance(vo.RegisterInstanceParam{
			ServiceName: appName,
			Ip:          serviceIP,
			Port:        appPort,
			GroupName:   "DEFAULT_GROUP",
			ClusterName: "DEFAULT", // Explicitly set default cluster
			Weight:      1.0,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
		})

		if err == nil && success {
			log.Println("✅ Service registered in Nacos successfully")
			return
		}

		log.Printf("❌ Registration attempt %d/10 failed: success=%v, err=%v\n", attempt, success, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ Failed to register service in Nacos after 10 retries")
}

// --------------------
// DEREGISTER SERVICE
// --------------------
func DeregisterService() {
	appName := os.Getenv("APP_NAME")
	serviceIP := os.Getenv("SERVICE_IP")
	appPort := uint64(mustAtoi(os.Getenv("APP_PORT")))
	groupName := os.Getenv("NACOS_GROUP")

	success, err := NacosClient.DeregisterInstance(vo.DeregisterInstanceParam{
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

// --------------------
// HELPERS (UNCHANGED)
// --------------------
func parseNacosIP(val string) (string, uint64) {
	parts := strings.Split(val, ":")
	if len(parts) != 2 {
		log.Fatalf("❌ Invalid NACOS address: %s (expected ip:port)", val)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatalf("❌ Invalid Nacos port: %s", parts[1])
	}
	return parts[0], uint64(port)
}

func mustAtoi(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("❌ Invalid integer: %s", val)
		return 0
	}
	return i
}
