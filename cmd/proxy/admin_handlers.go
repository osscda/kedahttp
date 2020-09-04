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
	appsCl := clientset.AppsV1().Deployments("cscaler")

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

		deployment, err := newDeployment(ctx, "cscaler", req.Name, req.ContainerImage)
		if err != nil {
			log.Printf("Error filling out new deployment (%s)", err)
			w.WriteHeader(400)
			return
		}
		// TODO: watch the deployment until it reaches ready state
		if _, err := appsCl.Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
			log.Printf("Error creating new deployment (%s)", err)
			w.WriteHeader(400)
			return
		}
		// TODO: create Service, then ScaledObject
		// use ExternalDNS to set up dynamic DNS's
		// https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/cloudflare.md

		w.WriteHeader(200)
	})
}
