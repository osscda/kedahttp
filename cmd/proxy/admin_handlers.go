package main

import (
	context "context"
	"encoding/json"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func newAdminDeployHandler() http.HandlerFunc {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting k8s config (%s)", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating k8s clientset (%s)", err)
	}

	type reqBody struct {
		Name           string `json:"name"`
		ContainerImage string `json:"image"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		req := new(reqBody)
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Printf("JSON decode error (%s)", err)
			w.WriteHeader(400)
			return
		}

		appsCl := clientset.AppsV1().Deployments("cscaler")
		deployment := newDeployment(ctx, "cscaler", req.Name, req.ContainerImage)
		// TODO: watch the deployment until it reaches ready state
		if _, err := appsCl.Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
			log.Printf("Error creating new deployment (%s)", err)
			w.WriteHeader(400)
			return
		}

		coreCl := clientset.CoreV1().Services("cscaler")
		service := newService(ctx, "cscaler", req.Name)
		if _, err := coreCl.Create(ctx, service, metav1.CreateOptions{}); err != nil {

		}
		// TODO: create ScaledObject and ClusterIP service

		w.WriteHeader(200)
	})
}
