package main

import (
	context "context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// import(
// "github.com/ericchiang/k8s"
// appsv1 "github.com/ericchiang/k8s/apis/apps/v1"
// corev1 "github.com/ericchiang/k8s/apis/core/v1"
// )

func int32P(i int32) *int32 {
	return &i
}

func str(s string) *string {
	return &s
}

func newDeployment(ctx context.Context, namespace, name, image string) (*appsv1.Deployment, error) {
	labels := map[string]string{
		"name": name,
		"app":  fmt.Sprintf("cscaler-%s", name),
	}
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "cscaler",
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: int32P(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Image:           image,
							Name:            name,
							ImagePullPolicy: "Always",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									ContainerPort: 8080,
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "PORT",
									Value: "8080",
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment, nil
	// appsCl := clientset.AppsV1().Deployments(namespace)
	// appsCl.Create(ctx, deployment, metav1.CreateOptions{})

}

// 	return &appsv1.Deployment{
// 		// Metadata:
// 		Spec: &appsv1.DeploymentSpec{
// 			Replicas: k8s.Int32(1),
// 			// Selector:
// 			Template: &corev1.PodTemplateSpec{
// 				// Metadata:
// 				Spec: &corev1.PodSpec{
// 					Containers: []*corev1.Container{
// 						&corev1.Container{
// 							Image:           k8s.String(image),
// 							Name:            k8s.String(name),
// 							ImagePullPolicy: k8s.String("Always"),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }
