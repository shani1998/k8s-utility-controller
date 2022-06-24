package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shani1998/k8s-utility-controller/models"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNS = "default"
	appGroup  = "applicationGroup"
)

func responseWriter(w http.ResponseWriter, respBytes []byte, code int) {
	w.WriteHeader(code)
	_, err := w.Write(respBytes)
	if err != nil {
		log.Errorf("failed to write response %v", err)
	}
	return
}

// GetServices handler accepts incoming requests for list services, and it fetches
// the service information from the cluster and writes response back to the client.
func GetServices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Infof("Incomming request %s %s %s", r.Method, r.RequestURI, r.RemoteAddr)

	// list deployments for given namespace with context
	deployments, err := ListDeployments(r.Context(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("error listing deployments %v", err)
		responseWriter(w, []byte("unable to fetch services"), http.StatusServiceUnavailable)
	}

	// prepare response body
	var response []models.Service
	for _, deploy := range deployments.Items {
		log.Infof("name:%s, appGroup:%s, runningPodCount: %d/%d", deploy.GetName(), deploy.GetLabels()[appGroup],
			deploy.Status.ReadyReplicas, deploy.Status.Replicas)
		svc := models.Service{
			Name:             deploy.GetName(),
			ApplicationGroup: deploy.GetLabels()[appGroup],
			RunningPodsCount: int(deploy.Status.ReadyReplicas),
		}
		response = append(response, svc)
	}

	// encode response to byte object
	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Errorf("error marshaling response %v", err)
		responseWriter(w, []byte("failed to encode response"), http.StatusServiceUnavailable)
	}

	responseWriter(w, respBytes, http.StatusOK)
}

func GetServicesByAppLabel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Infof("Incomming request %s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, params.ByName(appGroup))

	// get all deployments for given app label
	listOptions := metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", appGroup, params.ByName(appGroup))}
	deployments, err := ListDeployments(r.Context(), listOptions)
	if err != nil {
		log.Errorf("error listing deployments %v", err)
		responseWriter(w, []byte("unable to fetch services"), http.StatusServiceUnavailable)
	}

	// prepare response body
	var response []models.Service
	for _, deploy := range deployments.Items {
		svc := models.Service{
			Name:             deploy.GetName(),
			ApplicationGroup: deploy.GetLabels()[appGroup],
			RunningPodsCount: int(deploy.Status.ReadyReplicas),
		}
		response = append(response, svc)
	}

	log.Debugf("got resp %+v", response)

	// encode response to byte object
	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Errorf("error marshaling response %v", err)
		responseWriter(w, []byte("failed to encode response"), http.StatusServiceUnavailable)
	}

	responseWriter(w, respBytes, http.StatusOK)
}
