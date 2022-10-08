package main

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

var dynClient dynamic.Interface
var mapper meta.RESTMapper

func initDynamicClient() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	dClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	dynClient = dClient
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		log.Fatal(err)
	}

	mapper = restmapper.NewDiscoveryRESTMapper(gr)
}
