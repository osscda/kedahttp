package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kubernetes-client/go/kubernetes/client"
	"github.com/kubernetes-client/go/kubernetes/config"
)

func newAdminDeployHandler() http.HandlerFunc {
	cfg, err := config.LoadKubeConfig()
	if err != nil {
		log.Fatalf("Error creating kube client (%s)", err)
	}

	cl := client.NewAPIClient(cfg)

	type reqBody struct {
		Name           string `json:"name"`
		ContainerImage string `json:"image"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := new(reqBody)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Printf("JSON decode error (%s)", err)
			w.WriteHeader(400)
			return
		}

		deployment := client.V1Deployment{
			Metadata: &client.V1ObjectMeta{
				Name:      req.Name,
				Namespace: "cscaler",
			},
			Spec: &client.V1DeploymentSpec{
				Replicas: 1,
				Template: &client.V1PodTemplateSpec{
					Metadata: &client.V1ObjectMeta{
						Labels: map[string]string{
							"name": req.Name,
							"app":  "cscaler-" + req.Name,
						},
					},
					Spec: &client.V1PodSpec{
						Containers: []client.V1Container{
							client.V1Container{
								Image:           req.ContainerImage,
								Name:            req.Name,
								ImagePullPolicy: "Always",
								Ports: []client.V1ContainerPort{
									client.V1ContainerPort{
										ContainerPort: 8080,
									},
								},
								Env: []client.V1EnvVar{
									client.V1EnvVar{
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
		_, _, err = cl.AppsV1ApiService.CreateNamespacedDeployment(
			context.Background(),
			"cscaler",
			deployment,
			nil,
		)
		if err != nil {
			log.Printf("Error creating deployment (%s)", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	})
}
