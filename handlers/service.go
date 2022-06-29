package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shani1998/k8s-utility-controller/models"
	log "github.com/sirupsen/logrus"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNS = "default"
	appGroup  = "applicationGroup"
)

var HealthChan = make(chan error)

func responseWriter(w http.ResponseWriter, respBytes []byte, code int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(respBytes)
	if err != nil {
		log.Errorf("failed to write response %v", err)
	}
	if code == http.StatusOK {
		//clear the error from the healthCheck variable
		HealthChan <- nil
	} else {
		// update health if response is not success
		HealthChan <- fmt.Errorf("%s", respBytes)
	}
}

func getResponseBytes(deployments *appv1.DeploymentList) ([]byte, error) {
	response := make([]models.Service, 0)

	// traverse through all deployments
	for _, deploy := range deployments.Items {
		svc := models.Service{
			Name:             deploy.GetName(),
			ApplicationGroup: deploy.GetLabels()[appGroup],
			RunningPodsCount: int(deploy.Status.ReadyReplicas),
		}
		response = append(response, svc)
	}

	// encode response to byte object
	return json.Marshal(response)

}

// GetServices handler accepts incoming requests for list services, and it fetches
// the service information from the cluster and writes response back to the client.
func GetServices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Infof("Incomming request %s %s %s", r.Method, r.RequestURI, r.RemoteAddr)

	// list deployments for given namespace with context
	deployments, err := ListDeployments(r.Context(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("error listing deployments %v", err)
		responseWriter(w, []byte("failed to list services"), http.StatusServiceUnavailable)
		return
	}

	// prepare response with fetched services
	respBytes, err := getResponseBytes(deployments)
	if err != nil {
		log.Errorf("error marshaling response %v", err)
		responseWriter(w, []byte("failed to list services"), http.StatusServiceUnavailable)
		return
	}

	responseWriter(w, respBytes, http.StatusOK)
	log.Infof("successfully written response")
}

// GetServicesByAppLabel handler fetches list of deployments with given app group in default ns
// and write response back to the client
func GetServicesByAppLabel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Infof("Incomming request %s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, params.ByName(appGroup))

	// get all deployments for given app label
	listOptions := metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", appGroup, params.ByName(appGroup))}
	deployments, err := ListDeployments(r.Context(), listOptions)
	if err != nil {
		log.Errorf("error listing deployments %v", err)
		responseWriter(w, []byte("failed to list services"), http.StatusServiceUnavailable)
		return
	}

	// prepare response with fetched services
	respBytes, err := getResponseBytes(deployments)
	if err != nil {
		log.Errorf("error marshaling response %v", err)
		responseWriter(w, []byte("failed to list services"), http.StatusServiceUnavailable)
		return
	}

	responseWriter(w, respBytes, http.StatusOK)
	log.Infof("successfully written response")
}
