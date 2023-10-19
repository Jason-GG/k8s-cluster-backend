package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	builtin "net/url"
	"strings"

	"github.com/sjian_mstr/cluster-management/enums"
	"golang.org/x/crypto/ssh"
)

var (
	__sockPath__ = "/var/run/docker.sock"
	__socket__   = "--unix-socket"
)

const (
	MethodGet    = "-XGET"
	MethodPost   = "-XPOST"
	MethodPut    = "-XPUT"
	MethodHead   = "-XHEAD"
	MethodDelete = "-XDELETE"
)

type ipcBuilder struct {
	url      string
	header   http.Header
	form     builtin.Values
	method   string
	socket   bool
	sockAddr string
	sockFlag string

	body []byte
	json interface{}
}

type CurlResult struct {
}

func CurlGet(uri string) *ipcBuilder {
	return newIpcBuilder(uri, MethodGet)
}

func CurlPost(uri string) *ipcBuilder {
	return newIpcBuilder(uri, MethodPost)
}

func CurlDelete(uri string) *ipcBuilder {
	return newIpcBuilder(uri, MethodDelete)
}

func CurlPut(uri string) *ipcBuilder {
	return newIpcBuilder(uri, MethodPut)
}

func CurlHead(uri string) *ipcBuilder {
	return newIpcBuilder(uri, MethodHead)
}

func (ipc *ipcBuilder) Form(form builtin.Values) *ipcBuilder {
	ipc.form = form
	return ipc
}

func (ipc *ipcBuilder) Json(json interface{}) *ipcBuilder {
	// ipc.header.Set(ContentType, ApplicationJSON)
	ipc.json = json
	return ipc
}

func (ipc *ipcBuilder) createForm(form builtin.Values) error {
	if form == nil {
		return errors.New("invalid form")
	}

	// ipc.header.Set(ContentType, ApplicationFormUrlencoded)
	ipc.form = form

	return nil
}

func (ipc *ipcBuilder) createJson() error {
	b, err := json.Marshal(ipc.json)
	if err != nil {
		return err
	}

	ipc.body = b
	return nil
}

func (ipc *ipcBuilder) buildHeader() string {
	var final string

	if ipc.header == nil {
		return final
	}

	for k, v := range ipc.header {
		var concat string
		if len(v) > 0 {
			concat = strings.Join(v, ",")
		}

		final = fmt.Sprintf(" -H \"%s: %s\" ", k, concat)
	}

	return final
}

func (ipc *ipcBuilder) buildBase() string {
	if ipc.socket {
		ipc.sockAddr = __sockPath__
		ipc.sockFlag = __socket__
	}

	return fmt.Sprintf("curl %s %s %s ", ipc.sockFlag, ipc.sockAddr, ipc.method)
}

func (ipc *ipcBuilder) Build() string {
	final := ipc.buildHeader()

	base := ipc.buildBase()

	var sendBody string

	if ipc.json != nil {
		if err := ipc.createJson(); err != nil {
			panic(err)
		}
	}

	reqUrl := ipc.url

	if ipc.form != nil {
		f := ipc.form.Encode()
		reqUrl = fmt.Sprintf("\"%s?%s\"", ipc.url, f)
		// reqUrl = ipc.url + "?" + f
	}

	if ipc.body == nil {
		sendBody = ""
	} else {
		sendBody = fmt.Sprintf(" -d \"%v\"", string(ipc.body))
	}

	return base + final + reqUrl + sendBody
}

func newIpcBuilder(url string, method string) *ipcBuilder {
	return &ipcBuilder{
		url:    url,
		header: http.Header{},
		form:   builtin.Values{},
		method: method,
		socket: true,
	}
}

type HostMeta struct {
	Hostname         string            `json:"hostname" bson:"hostname"`
	InstanceId       string            `json:"instanceId" bson:"instanceId"`
	VpcId            string            `json:"vpcId" bson:"vpcId"`
	VswitchId        string            `json:"vswitchId" bson:"vswitchId"`
	Region           string            `json:"region" bson:"region"`
	PrivateIpAddress string            `json:"privateIpAddress" bson:"privateIpAddress"`
	ElasticIpAddress string            `json:"elasticIpAddress" bson:"elasticIpAddress"`
	BeingUsed        bool              `json:"beingUsed" bson:"beingUsed"`
	Apps             []App             `json:"apps" bson:"apps"`
	Env              int               `json:"env"`
	ContainerIDs     []string          `json:"containerIds" bson:"containerIds"`
	Stats            Stats             `json:"stats" bson:"stats"`
	InstanceType     string            `json:"instanceType" bson:"instanceType"`
	Tag              map[string]string `json:"tag" bson:"tag"`
	//SprintPlanId     int      `json:"sprintPlanId"`
	//AssociateProject []int    `json:"associateProject"`
}

type App struct {
	Env          int            `json:"env" bson:"env"`
	SprintPlanId int            `json:"sprintPlanId" bson:"sprintPlanId"`
	AppName      string         `json:"appName" bson:"appName"`
	UsedPort     []PortReprsent `json:"usedPort" bson:"usedPort"`
}

type PortReprsent struct {
	PortType enums.PortType `json:"portType" bson:"portType"`
	PortNum  int            `json:"port" bson:"port"`
}

func PruneContainer(client *ssh.Client) ([]byte, error) {
	return shellCmd(client, CurlPost(enums.ApiContainerPrune).Json(nil).Build())
}

func ListContainer(client *ssh.Client) ([]byte, error) {
	return shellCmd(client, CurlGet(enums.ApiContainerList).Json(nil).Build())
}

func InspectContainer(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlGet(fmt.Sprintf(enums.ApiContainerInspect, id)).Json(nil).Build())
}

func ListContainerProc(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlGet(fmt.Sprintf(enums.ApiContainerListProc, id)).Json(nil).Build())
}

func LogsContainer(client *ssh.Client, id string, form builtin.Values) ([]byte, error) {
	return shellCmd(client, CurlGet(fmt.Sprintf(enums.ApiContainerLogs, id)).Form(form).Build())
}

func StartContainer(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlPost(fmt.Sprintf(enums.ApiContainerStart, id)).Json(nil).Build())
}

func StopContainer(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlPost(fmt.Sprintf(enums.ApiContainerStop, id)).Json(nil).Build())
}

func RestartContainer(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlPost(fmt.Sprintf(enums.ApiContainerRestart, id)).Json(nil).Build())
}

func KillContainer(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlPost(fmt.Sprintf(enums.ApiContainerKill, id)).Json(nil).Build())
}

//	func UpdateContainer(client *ssh.Client, id string) ([]byte, error) {
//		return shellCmd(client, fmt.Sprintf("%s %s %s %s", __curlSocket__, MethodPost, id, enums.ApiContainerUpdate))
//	}
func RemoveContainer(client *ssh.Client, id string) ([]byte, error) {
	return shellCmd(client, CurlDelete(fmt.Sprintf(enums.ApiContainerRemove, id)).Json(nil).Build())
}

func ListContainerImage(client *ssh.Client) ([]byte, error) {
	return shellCmd(client, CurlGet(enums.ApiContainerImageList).Json(nil).Build())
}

func PruneContainerImage(client *ssh.Client) ([]byte, error) {
	return shellCmd(client, CurlPost(enums.ApiContainerImagePrune).Json(nil).Build())
}
