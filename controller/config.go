package controller

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sjian_mstr/cluster-management/enums"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	"github.com/sjian_mstr/cluster-management/util"
	"github.com/sjian_mstr/cluster-management/util/homedir"
	"k8s.io/client-go/rest"
)

var (
	kube = map[string]*KubeConfig{
		"insecure":      nil,
		"prod":          nil,
		"dev":           nil,
	}
)

// init 函数自己执行
func init() {
	fmt.Println("... init Config yaml")
	// var array []KubeConfig
	var fileMap = util.SingletonConfig
	for _, fileInfo := range fileMap {
		fmt.Printf("Name: %s\n", fileInfo.Name)
		fmt.Printf("Path: %s\n", fileInfo.Path)
		kube[fileInfo.Name] = new(KubeConfig)
		kube[fileInfo.Name].Name = fileInfo.Name
		kube[fileInfo.Name].Configpath = fileInfo.Path

	}
	fmt.Println(kube)

}

type KubeConfig struct {
	Name       string
	Configpath string
}

// defaultPort initialized the port representation, give it a default value
var defaultPort = []util.PortReprsent{
	{PortType: enums.PortApplication, PortNum: 8080},
	{PortType: enums.PortSchedule, PortNum: 9090},
	{PortType: enums.PortDebug, PortNum: 5050},
}

var generate = func() []util.PortReprsent {
	var (
		rd                 = randomUint()
		application uint64 = 8000
		debug       uint64 = 5000
		sche        uint64 = 9000
	)

	return []util.PortReprsent{
		{PortType: enums.PortApplication, PortNum: atomicAdd(&application, rd)},
		{PortType: enums.PortDebug, PortNum: atomicAdd(&debug, rd)},
		{PortType: enums.PortSchedule, PortNum: atomicAdd(&sche, rd)},
	}
}

func parseUnixTimeStamp(ts int) string {
	unixTimeUTC := time.Unix(int64(ts), 0)
	return unixTimeUTC.Format(time.RFC3339)
}

func strToInt(str string) int {
	r, _ := strconv.Atoi(str)
	return r
}

func randomUint() uint64 {
	seed := rand.NewSource(time.Now().UnixNano())
	return uint64(rand.New(seed).Int63n(999))
}

func parseBool(str string) bool {
	return util.ParseBool(str)
}

func atomicAdd(addr *uint64, delta uint64) int {
	return int(atomic.AddUint64(addr, delta))
}

////////////////////////////// self define ////////////////////////////

func GetKubeConfig() (*rest.Config, error) {
	// Load Kubernetes configuration from the default location or a specified kubeconfig file.
	var kubeconfigPath string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = home + "/.kube/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}
