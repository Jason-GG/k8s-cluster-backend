package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/controller"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentResponse struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Most recently observed status of the Deployment.
	// +optional
	Status appsv1.DeploymentStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

func getDeploymentListInfo(c *fiber.Ctx) error {
	// var modifiedData []DeploymentResponse

	queryParams := c.Queries()
	// Retrieve specific query parameters
	env := queryParams["env"]
	namespace := queryParams["namespace"]

	deploymentInfo, err := controller.ListDeployments(env, namespace)
	if err != nil {
		fmt.Printf("Error ListDeployments: %v\n", err)
		// os.Exit(1)
		return jsonResponse(c, fiber.StatusOK, "OK", 1, err)

	}
	for _, deployment := range deploymentInfo {
		fmt.Printf("Name: %s\n", deployment.ObjectMeta.Name)
		fmt.Println("-----------------------")
	}

	return jsonResponse(c, fiber.StatusOK, "OK", 1, deploymentInfo)
}
