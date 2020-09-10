package main

import (
	"encoding/json"

	"github.com/arschles/containerscaler/pkg/k8s"
	echo "github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func newAdminDeleteDeploymentHandler(
	k8sCl *kubernetes.Clientset,
	dynCl dynamic.Interface,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		deployName := c.QueryParam("name")
		scaledObjectName := c.QueryParam("name")
		scaledObjectCl := k8s.NewScaledObjectClient(dynCl)
		if deployName == "" {
			return c.String(400, "'name' query param required")
		}
		if err := k8s.DeleteService(ctx, deployName, k8sCl.CoreV1().Services("cscaler")); err != nil {
			return err
		}
		if err := k8s.DeleteDeployment(ctx, deployName, k8sCl.AppsV1().Deployments("cscaler")); err != nil {
			return err
		}
		if err := k8s.DeleteScaledObject(ctx, scaledObjectName, scaledObjectCl); err != nil {
			return err
		}
		return nil
	}
}

func newAdminCreateDeploymentHandler(
	k8sCl *kubernetes.Clientset,
	dynCl dynamic.Interface,
) echo.HandlerFunc {

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
			return echo.NewHTTPError(400, "Could not decode request")
		}

		appsCl := k8sCl.AppsV1().Deployments("cscaler")
		deployment := k8s.NewDeployment(ctx, "cscaler", req.Name, req.ContainerImage)
		// TODO: watch the deployment until it reaches ready state
		if _, err := appsCl.Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
			return echo.NewHTTPError(500, "Error creating the new deployment (%s)", err)
		}

		coreCl := k8sCl.CoreV1().Services("cscaler")
		service := k8s.NewService("cscaler", req.Name)
		if _, err := coreCl.Create(ctx, service, metav1.CreateOptions{}); err != nil {
			return echo.NewHTTPError(500, "Error creating the new service (%s)", err)
		}
		scaledObjectCl := k8s.NewScaledObjectClient(dynCl)
		_, err := scaledObjectCl.Namespace("cscaler").Create(ctx, k8s.NewScaledObject(
			"cscaler",
			req.Name,
			req.Name,
		), metav1.CreateOptions{})
		if err != nil {
			return echo.NewHTTPError(500, "Creating the new scaled object (%s)", err)
		}

		w.WriteHeader(200)
		return nil
	}
}
