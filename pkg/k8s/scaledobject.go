package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func kedaGVR() *schema.GroupVersionResource {
	return &schema.GroupVersionResource{
		Group:    "keda.k8s.io",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
}

// NewScaledObjectClient returns a new dynamic client capable
// of interacting with ScaledObjects in a cluster
func NewScaledObjectClient(cl dynamic.Interface) dynamic.NamespaceableResourceInterface {
	return cl.Resource(kedaGVR())
}

// NewScaledObject creates a new ScaledObject in memory
func NewScaledObject(namespace, name, deploymentName string) *ScaledObject {
	// https://keda.sh/docs/1.5/faq/
	// https://github.com/kedacore/keda/blob/v2/api/v1alpha1/scaledobject_types.go
	return &ScaledObject{
		TypeMeta: metav1.TypeMeta{
			Kind: "ScaledObject",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels(name),
		},
		Spec: ScaledObjectSpec{
			MinReplicaCount: int32(0),
			MaxReplicaCount: int32(1000),
			PollingInterval: int32(1),
			ScaleTargetRef: &ScaleTarget{
				Name: deploymentName,
				Kind: "deployment",
			},
		},
	}
	return nil
}
