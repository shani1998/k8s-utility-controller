package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

const defaultNS = "default"

func GetServices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Infof("Incomming request %s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
	fmt.Fprint(w, "GetServices!")
}

func GetServicesByAppLabel(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Infof("Incomming request %s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
	fmt.Fprint(w, "GetServicesByAppLabel!")
}
