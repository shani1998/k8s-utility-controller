package handlers

import (
	"context"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeClient kubernetes.Interface

// InitKubeClient reads the current cluster config and initializes api client
// and returns errors if unable to retrieve the token.
func InitKubeClient() error {
	log.Infof("intializing kube client")

	// read config based on the service account token
	conf, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("error getting in cluster config: %v", err)
		// If no in-cluster config, try the default location in the user's home directory
		conf, err = getOutClusterConfig()
		if err != nil {
			log.Errorf("error getting OutClusterConfigs %v", err)
			return err
		}
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Errorf("error getting kube clinet: %v", err)
		return err
	}
	kubeClient = clientset
	log.Infof("successfully initialized kube client")

	return nil
}

// getOutClusterConfig
func getOutClusterConfig() (*rest.Config, error) {
	var kubeconfig *string

	log.Infof("fetching out cluster config")
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = pflag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	}

	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
}

func ListDeployments(ctx context.Context, opts metav1.ListOptions) (*appv1.DeploymentList, error) {
	log.Infof("fetching list of deployments with label %s", opts.LabelSelector)
	listDeployCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return kubeClient.AppsV1().Deployments(defaultNS).List(listDeployCtx, opts)

}
