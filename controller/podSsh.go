package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/sjian_mstr/cluster-management/kubernetes"
	"github.com/sjian_mstr/cluster-management/kubernetes/scheme"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	"github.com/sjian_mstr/cluster-management/tools/remotecommand"
	v1 "k8s.io/api/core/v1"
)

// ssh流式处理器
type streamHandler struct {
	wsConn      *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

// web terminal package
type xtermMessage struct {
	MsgType string `json:"type"`  // class:resize client type
	Input   string `json:"input"` // when sgtype=input use it
	Rows    uint16 `json:"rows"`  // when msgtype=resize use it
	Cols    uint16 `json:"cols"`  // when msgtype=resize use it
}

// Next executor retrive web command to check if resize
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

func (handler *streamHandler) Read(p []byte) (size int, err error) {
	var (
		msg      *WsMessage
		xtermMsg xtermMessage
	)

	if msg, err = handler.wsConn.WsRead(); err != nil {
		return
	}

	// new
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return
	}
	if xtermMsg.MsgType == "resize" {
		handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" {
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

func (handler *streamHandler) Write(p []byte) (size int, err error) {
	var copyData []byte
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.wsConn.WsWrite(websocket.TextMessage, copyData)
	return
}

func ExecuteCommand(c *websocket.Conn) {
	var (
		wsConn   *WsConnection
		executor remotecommand.Executor
		handler  *streamHandler
		err      error
	)
	fmt.Println("=====================>>>>>>>>>>>>>>>>>>>ExecuteCommand<<<<<<<<<<<<<<<<<<<<<==========================")
	env := c.Query("env")
	podName := c.Query("podName")
	namespace := c.Query("namespace")
	containerName := c.Query("containerName")
	fmt.Println(env)
	fmt.Println(podName)
	fmt.Println(namespace)
	fmt.Println(containerName)

	config, err := clientcmd.BuildConfigFromFlags("", kube[env].Configpath)
	if err != nil {
		return
	}

	// Create a Kubernetes clientset using the provided kubeconfig.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	sshReq := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	if wsConn, err = InitWebsocket(c); err != nil {
		return
	}
	if executor, err = remotecommand.NewSPDYExecutor(config, "POST", sshReq.URL()); err != nil {
		// goto END
		return
	}

	handler = &streamHandler{wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	ctx := context.Background()
	if err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		fmt.Println("err messages: ", err)
		goto END
	}
	return

END:
	wsConn.WsClose()

}
