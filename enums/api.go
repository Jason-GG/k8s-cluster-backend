package enums

// aliyun host meta data api
const (
	ApiContainerList     = "http://127.0.0.1/containers/json"
	ApiContainerInspect  = `http://127.0.0.1/containers/%s/json`
	ApiContainerListProc = `http://127.0.0.1/containers/%s/top`
	ApiContainerLogs     = `http://127.0.0.1/containers/%s/logs`
	ApiContainerStart    = `http://127.0.0.1/containers/%s/start`
	ApiContainerStop     = `http://127.0.0.1/containers/%s/stop`
	ApiContainerRestart  = `http://127.0.0.1/containers/%s/restart`
	ApiContainerKill     = `http://127.0.0.1/containers/%s/kill`
	ApiContainerUpdate   = `http://127.0.0.1/containers/%s/update`
	ApiContainerRemove   = `http://127.0.0.1/containers/%s`
	ApiContainerPause    = `http://127.0.0.1/containers/%s/pause`
	ApiContainerUnpause  = `http://127.0.0.1/containers/%s/unpause`
	ApiContainerPrune    = "http://127.0.0.1/containers/prune"
)

const (
	ApiContainerImageList    = "http://127.0.0.1/images/json"
	ApiContainerImageInspect = `http://127.0.0.1/images/%s/json`
	ApiContainerImagePrune   = `http://127.0.0.1/images/prune`
)

const (
	StatusCreated    = "created"
	StatusRestarting = "restarting"
	StatusRunning    = "running"
	StatusRemoving   = "removing"
	StatusPaused     = "paused"
	StatusExited     = "exited"
	StatusDead       = "dead"
)

var GpuInstanceSpecs = new(GpuInstanceSpec)

type GpuInstanceSpec struct{}

// func (*GpuInstanceSpec) HasSpec(typ string) bool {
// 	if _, ok := _gpuInstanceSpec[typ]; ok {
// 		return ok
// 	}
// 	return false
// }

// var _gpuInstanceSpec = map[string]interface{}{
// 	"ecs.gn6i-c4g1.xlarge":    struct{}{},
// 	"ecs.gn6i-c8g1.2xlarge":   struct{}{},
// 	"ecs.gn6i-c16g1.4xlarge":  struct{}{},
// 	"ecs.gn6i-c24g1.6xlarge":  struct{}{},
// 	"ecs.gn6i-c24g1.12xlarge": struct{}{},
// 	"ecs.gn6i-c24g1.24xlarge": struct{}{},
// 	"ecs.gn6e-c12g1.3xlarge":  struct{}{},
// 	"ecs.gn6e-c12g1.12xlarge": struct{}{},
// 	"ecs.gn6e-c12g1.24xlarge": struct{}{},
// 	"ecs.gn6v-c8g1.2xlarge":   struct{}{},
// 	"ecs.gn6v-c8g1.8xlarge":   struct{}{},
// 	"ecs.gn6v-c8g1.16xlarge":  struct{}{},
// 	"ecs.gn6v-c10g1.20xlarge": struct{}{},
// 	"ecs.ebmgn6e.24xlarge":    struct{}{},
// 	"ecs.ebmgn6v.24xlarge":    struct{}{},
// 	"ecs.ebmgn6i.24xlarge":    struct{}{},
// 	"ecs.gn5-c4g1.xlarge":     struct{}{},
// 	"ecs.gn5-c8g1.2xlarge":    struct{}{},
// 	"ecs.gn5-c4g1.2xlarge":    struct{}{},
// 	"ecs.gn5-c8g1.4xlarge":    struct{}{},
// 	"ecs.gn5-c28g1.7xlarge":   struct{}{},
// 	"ecs.gn5-c8g1.8xlarge":    struct{}{},
// 	"ecs.gn5-c28g1.14xlarge":  struct{}{},
// 	"ecs.gn5-c8g1.14xlarge":   struct{}{},
// 	"ecs.gn5i-c2g1.large":     struct{}{},
// 	"ecs.gn5i-c4g1.xlarge":    struct{}{},
// 	"ecs.gn5i-c8g1.2xlarge":   struct{}{},
// 	"ecs.gn5i-c16g1.4xlarge":  struct{}{},
// 	"ecs.gn5i-c16g1.8xlarge":  struct{}{},
// 	"ecs.gn5i-c28g1.14xlarge": struct{}{},
// 	"ecs.gn4-c4g1.xlarge":     struct{}{},
// 	"ecs.gn4-c8g1.2xlarge":    struct{}{},
// 	"ecs.gn4.8xlarge":         struct{}{},
// 	"ecs.gn4-c4g1.2xlarge":    struct{}{},
// 	"ecs.gn4-c8g1.4xlarge":    struct{}{},
// 	"ecs.gn4.14xlarge":        struct{}{},
// }
