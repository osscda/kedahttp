package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/arschles/containerscaler/pkg/k8s"
	echo "github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// NAMESPACE NAME
const NAMESPACE string = "cscaler"

func newAdminDeleteAppHandler(
	k8sCl *kubernetes.Clientset,
	dynCl dynamic.Interface,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Logger()
		ctx := c.Request().Context()
		deployName := c.QueryParam("name")
		scaledObjectCl := k8s.NewScaledObjectClient(dynCl).Namespace(NAMESPACE)
		if deployName == "" {
			logger.Errorf("'name' query param not found")
			return c.String(400, "'name' query param required")
		}
		if err := k8s.DeleteService(ctx, deployName, k8sCl.CoreV1().Services(NAMESPACE)); err != nil {
			logger.Errorf("Deleting service %s (%s)", err)
			return c.String(500, "deleting service")
		}
		if err := k8s.DeleteIngress(ctx, deployName, k8sCl.ExtensionsV1beta1().Ingresses(NAMESPACE)); err != nil {
			logger.Errorf("Deleting ingress %s (%s)", err)
			return c.String(500, "deleting ingress")
		}
		if err := k8s.DeleteDeployment(ctx, deployName, k8sCl.AppsV1().Deployments(NAMESPACE)); err != nil {
			logger.Errorf("Deleting deployment %s (%s)", deployName, err)
			return c.String(500, "deleting deployment")
		}
		if err := k8s.DeleteScaledObject(ctx, deployName, scaledObjectCl); err != nil {
			logger.Errorf("Deleting scaledobject %s (%s)", deployName, err)
			return c.String(500, "deleting scaledobject")
		}
		c.String(200, "deleted")
		return nil
	}
}

func newAdminCreateAppHandler(
	k8sCl *kubernetes.Clientset,
	dynCl dynamic.Interface,
	scalerAddress string,
) echo.HandlerFunc {

	type reqBody struct {
		Name           string `json:"name"`
		ContainerImage string `json:"image"`
		Port           string `json:"port"`
	}

	return func(c echo.Context) error {
		r := c.Request()
		ctx := c.Request().Context()
		req := new(reqBody)
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			c.Logger().Errorf("Decoding request (%s)", err)
			return c.String(400, "decoding request")
		}

		portInt, err := strconv.Atoi(req.Port)
		if err != nil {
			c.Logger().Errorf("Invalid port %s (%s)", req.Port, err)
			return c.String(400, "invalid port")
		}

		appsCl := k8sCl.AppsV1().Deployments(NAMESPACE)
		deployment := k8s.NewDeployment(ctx, NAMESPACE, req.Name, req.ContainerImage, int32(portInt))
		// TODO: watch the deployment until it reaches ready state
		if _, err := appsCl.Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
			c.Logger().Errorf("Creating deployment (%s)", err)
			return c.String(500, "creating deployment")
		}

		coreCl := k8sCl.CoreV1().Services(NAMESPACE)
		service := k8s.NewService(NAMESPACE, req.Name, int32(portInt))
		if _, err := coreCl.Create(ctx, service, metav1.CreateOptions{}); err != nil {
			c.Logger().Errorf("Creating service (%s)", err)
			return c.String(500, "creating service")
		}

		coreIng := k8sCl.ExtensionsV1beta1().Ingresses(NAMESPACE)
		dnsZoneName := ""
		cscalerAdminIngress, err := coreIng.Get(ctx, "cscaler-admin", metav1.GetOptions{});
		if err != nil {
			c.Logger().Errorf("Getting proxy DNS zone name (%s)", err)
			return c.String(500, "getting proxy DNS Zone Name")
		}
		dnsZoneName = strings.Join(strings.Split(cscalerAdminIngress.Spec.Rules[0].Host, ".")[1:], ".") // Get HAR hostname (DNS Zone Name)

		ingress := k8s.NewIngress(NAMESPACE, req.Name, dnsZoneName)
		if _, err := coreIng.Create(ctx, ingress, metav1.CreateOptions{}); err != nil {
			c.Logger().Errorf("Creating ingress (%s)", err)
			return c.String(500, "creating ingress")
		}

		scaledObjectCl := k8s.NewScaledObjectClient(dynCl)
		_, err = scaledObjectCl.Namespace(NAMESPACE).Create(ctx, k8s.NewScaledObject(
			NAMESPACE,
			req.Name,
			req.Name,
			scalerAddress,
		), metav1.CreateOptions{})
		if err != nil {
			c.Logger().Errorf("Creating scaledobject (%s)", err)
			return c.String(500, "creating scaledobject")
		}

		return nil
	}
}
