package k8s

import (
	context "context"
	"strings"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8sextensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

func DeleteIngress(ctx context.Context, name string, cl k8sextensionsv1beta1.IngressInterface) error {
	return cl.Delete(ctx, name, metav1.DeleteOptions{})
}

func NewIngress(namespace, name string, dnsZoneName string) *extensionsv1beta1.Ingress {
	return &extensionsv1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind: "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "cscaler",
			Labels:    labels(name),
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "addon-http-application-routing",
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: strings.Join([]string{name, dnsZoneName}, "."),
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "cscaler-proxy",
										ServicePort: intstr.FromString("web"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
