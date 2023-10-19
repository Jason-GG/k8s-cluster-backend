package controller

import (
	"context"
	"fmt"

	"github.com/sjian_mstr/cluster-management/kubernetes"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	v1 "k8s.io/api/core/v1" // Import corev1 from here
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListPodsInfo(env string, namespace string) ([]v1.Pod, error) {
	// Path to kubeconfig file (default location is $HOME/.kube/config)

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

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func GetDeploymentPods(env string, deploymentName string, namespace string) ([]v1.Pod, error) {
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

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	selector := deployment.Spec.Selector
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(selector),
	})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func DeletePod(env string, namespace string, podName string) {
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

	clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})

}
