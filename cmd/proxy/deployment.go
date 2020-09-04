package main

import (
	"github.com/ericchiang/k8s"
	appsv1 "github.com/ericchiang/k8s/apis/apps/v1"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

func newDeployment(name, image string) *appsv1.Deployment {
	return &appsv1.Deployment{
		// Metadata:
		Spec: &appsv1.DeploymentSpec{
			Replicas: k8s.Int32(1),
			// Selector:
			Template: &corev1.PodTemplateSpec{
				// Metadata:
				Spec: &corev1.PodSpec{
					Containers: []*corev1.Container{
						&corev1.Container{
							Image:           k8s.String(image),
							Name:            k8s.String(name),
							ImagePullPolicy: k8s.String("Always"),
						},
					},
				},
			},
		},
	}
}
