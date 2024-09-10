package dns

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func RestK8sClient() (*kubernetes.Clientset, error) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var clientset *kubernetes.Clientset

	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "/home/allank/.kube/config")
	flag.Parse()

	if len(strings.TrimSpace(*kubeconfig)) > 0 {
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			slog.Error("Error creating k8s config: ", "%v", err.Error())
		}

		if config == nil {
			config, err = rest.InClusterConfig()
			if err != nil {
				slog.Error("Error creating a rest client:", "%v", err)
				return nil, err
			}
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			slog.Error("Error creating kubernetes client")
			return nil, err
		}
	}
	return clientset, nil
}
