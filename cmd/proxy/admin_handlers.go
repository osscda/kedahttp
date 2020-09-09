package main

import (
	"encoding/json"
	"log"

	"github.com/labstack/echo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func newAdminDeleteDeploymentHandler(k8sCl *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		deployName := c.QueryParam("name")
		if deployName == "" {
			return c.String(400, "'name' query param required")
		}
		if err := deleteService(ctx, deployName, k8sCl.CoreV1().Services("cscaler")); err != nil {
			return err
		}
		if err := deleteDeployment(ctx, deployName, k8sCl.AppsV1().Deployments("cscaler")); err != nil {
			return err
		}
		// if err := deleteScaledObject(deployName); err != nil {
		// 	return err
		// }
		return nil
	}
}

func newAdminCreateDeploymentHandler(k8sCl *kubernetes.Clientset) echo.HandlerFunc {

	type reqBody struct {
		Name           string `json:"name"`
		ContainerImage string `json:"image"`
	}

	return func(c echo.Context) error {
		r := c.Request()
		w := c.Response()
		ctx := c.Request().Context()
		req := new(reqBody)
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Printf("JSON decode error (%s)", err)
			w.WriteHeader(400)
			return nil
		}

		appsCl := k8sCl.AppsV1().Deployments("cscaler")
		deployment := newDeployment(ctx, "cscaler", req.Name, req.ContainerImage)
		// TODO: watch the deployment until it reaches ready state
		if _, err := appsCl.Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
			log.Printf("Error creating new deployment (%s)", err)
			w.WriteHeader(400)
			return nil
		}

		coreCl := k8sCl.CoreV1().Services("cscaler")
		service := newService(ctx, "cscaler", req.Name)
		if _, err := coreCl.Create(ctx, service, metav1.CreateOptions{}); err != nil {
			log.Printf("Error creating new service (%s)", err)
			w.WriteHeader(400)
			return nil
		}
		// TODO: create ScaledObject and ClusterIP service

		w.WriteHeader(200)
		return nil
	}
}
