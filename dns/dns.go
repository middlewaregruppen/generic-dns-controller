package dns

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type DNSProvider interface {
	CreateRecord(name, value string) error
	CreateRecordHttp(name, value string) error
	DeleteRecord(name, value string) error
	SearchRecord(name string) (bool, error)
	SearchRecordHttp(name string) (SeaRecord, error)
	UpdateRecord(name string) error
}

type DNSController struct {
	clientset   *kubernetes.Clientset
	dnsProvider DNSProvider
}

var (
	INGRESS_NAMESPACE        = os.Getenv("INGRESS_NAMESPACE")
	INGRESS_CONTROLLER_NAME  = os.Getenv("INGRESS_CONTROLLER_NAME")
	INGRESS_CONTROLLER_LABEL = "app.kubernetes.io/name"
	INGRESS_DNS_ANNOTATION   = os.Getenv("INGRESS_DNS_ANNOTATION")
)

func NewDNSController(clientset *kubernetes.Clientset, dnsProvider DNSProvider) *DNSController {
	fmt.Println("Initialising DNS Provider")
	return &DNSController{
		clientset:   clientset,
		dnsProvider: dnsProvider,
	}
}

func (c *DNSController) Run(ctx context.Context) {
	for {
		ingressWatcher, err := c.clientset.NetworkingV1().Ingresses(metav1.NamespaceAll).Watch(ctx, metav1.ListOptions{})
		if err != nil {
			slog.Info("Error watching ingresses")
		}

		for event := range ingressWatcher.ResultChan() {
			fmt.Printf("Got event type of %v\n", event.Type)
			switch event.Type {
			case watch.Added:
				ingress := event.Object.(*v1.Ingress)
				if notmanagedByExternalDns(ingress) {
					c.CreateRecord(ingress)
					//c.CreateRecordHttp(ingress)
				}
			case watch.Deleted:
				ingress := event.Object.(*v1.Ingress)
				if notmanagedByExternalDns(ingress) {
					c.DeleteRecord(ingress)
				}
			case watch.Modified:
				//Get old object and update
				// Check for annotations.
				ingress := event.Object.(*v1.Ingress)
				c.UpdateRecord(ingress)
			default:
				slog.Info("unknown operation for DNS controller")
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (c *DNSController) CreateRecord(ingress *v1.Ingress) {
	name := ingress.Name
	class := ingress.Spec.IngressClassName
	defaultClass := "nginx"
	if class == nil {
		class = &defaultClass
	}

	if notmanagedByExternalDns(ingress) {
		slog.Info("Record not managed by an DNS controller endpoint")
		//return
	}

	rules := ingress.Spec.Rules
	for _, rule := range rules {
		host := rule.Host
		dnsName := strings.Split(host, ".")[0]
		ref, _ := c.SearchRecordHttp(dnsName)
		if len(ref.Ref) == 0 {
			slog.Info("Getting IP for controller ", *class, nil)
			ip := getIngressControllerIp(*class)
			if len(ip) == 0 {
				slog.Error("Controller %v is missing an IP address - no DNS record will be created.", *class, nil)
				return
			}
			slog.Info("DNS record %s record to be created for ingress - %v\n", dnsName, name)
			//err := c.dnsProvider.CreateRecord(name, ip)
			err := c.dnsProvider.CreateRecordHttp(dnsName, ip)
			if err != nil {
				slog.Error("Unable to create DNS record: ", err.Error(), nil)
				return
			}
			slog.Info("Created DNS record %s for %s", dnsName, name)
		}
	}
}

func (c *DNSController) CreateRecordHttp(ingress *v1.Ingress) {
	log.Info().Msg("Creating record using HTTP path")
}

func (c *DNSController) DeleteRecord(ingress *v1.Ingress) {
	name := ingress.Name
	rules := ingress.Spec.Rules
	for _, rule := range rules {
		host := rule.Host
		//This should actually be DeleteRecord(host,ipAddress)
		err := c.dnsProvider.DeleteRecord(name, host)
		if err != nil {
			slog.Error("unable to delete DNS record ", err.Error(), nil)
			return
		}
		slog.Info("DNS record deleted")
	}
}

func (c *DNSController) SearchRecord(name string) bool {
	slog.Info("Searching for ", name, nil)
	status, _ := c.dnsProvider.SearchRecord(name)
	return status
}

func (c *DNSController) SearchRecordHttp(name string) (SeaRecord, error) {
	ref, _ := c.dnsProvider.SearchRecordHttp(name)
	return ref, nil
}

func notmanagedByExternalDns(ingress *v1.Ingress) bool {
	if len(INGRESS_DNS_ANNOTATION) == 0 {
		INGRESS_DNS_ANNOTATION = "managed-by-externaldns"
	}
	_, managed := ingress.Annotations[INGRESS_DNS_ANNOTATION]
	return !managed
}

func getIngressControllerIp(name string) string {
	clientset, err := RestK8sClient()
	if err != nil {
		slog.Error("error connecting to cluster ", err.Error(), nil)
	}

	ingressControllers, err := clientset.NetworkingV1().IngressClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		slog.Error("error retrieving ingresscontrollers ", err.Error(), nil)
	}

	var loadbalancerIp string

	for _, ic := range ingressControllers.Items {
		if ic.Name == name {
			for key, val := range ic.Labels {
				if strings.Contains(key, INGRESS_CONTROLLER_LABEL) {
					slog.Info("getting service with loadbalancer for ", name, nil)
					loadbalancerIp = processServicesList(clientset, val, name)
					if len(loadbalancerIp) > 0 {
						slog.Info("Got me some Ipie: ", loadbalancerIp, nil)
					}
				}
			}
		}
	}
	return loadbalancerIp
}

func processServicesList(clientset *kubernetes.Clientset, icAppName, name string) string {
	services, err := clientset.CoreV1().Services(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: INGRESS_CONTROLLER_LABEL,
	})
	if err != nil {
		slog.Error("error getting services ", err.Error(), nil)
	}

	var loadbalancerIp string
	for _, svc := range services.Items {
		for k, v := range svc.ObjectMeta.Labels {
			if v == icAppName && k == INGRESS_CONTROLLER_LABEL {
				//Does this controller have a LoadBalancerIp?
				if len(svc.Spec.LoadBalancerIP) > 0 {
					slog.Error(name, "'s IP address is ", svc.Spec.LoadBalancerIP)
					loadbalancerIp = svc.Spec.LoadBalancerIP
				}
				if len(loadbalancerIp) == 0 && len(os.Getenv("LOADBALANCER_IP")) > 0 {
					slog.Info(name, " ingress controller has no IP address.Perhaps it is a clusterIP or NodePort service", nil)
					slog.Info("assigning a hard-coded LoadBalancerIP variable")

					loadbalancerIp = os.Getenv("LOADBALANCER_IP")
				}
			}
		}
	}
	return loadbalancerIp
}

func (c *DNSController) UpdateRecord(ingress *v1.Ingress) error {
	var name string
	rules := ingress.Spec.Rules

	for _, rule := range rules {
		name = rule.Host
	}
	if len(name) > 0 {
		err := c.dnsProvider.UpdateRecord(name)
		return err
	}
	return nil
}
