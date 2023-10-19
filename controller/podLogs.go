package controller

import (
	"context"
	"fmt"

	"github.com/sjian_mstr/cluster-management/kubernetes"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	corev1 "k8s.io/api/core/v1"
)

// GetPodLogs retrieves logs from a pod and returns them as a string
func GetPodLogs(env string, namespace string, podName string, num int) (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kube[env].Configpath)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		// os.Exit(1)
	}
	// Create a Kubernetes clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		// os.Exit(1)
	}
	tailLines := int64(num) // Number of log lines to retrieve
	// Get pod logs
	podLogs, err := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		TailLines: &tailLines,
	}).Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	// Read the pod logs
	logs := ""
	buf := make([]byte, 1024)
	for {
		n, err := podLogs.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return "", err
		}
		logs += string(buf[:n])
	}

	return logs, nil
}
