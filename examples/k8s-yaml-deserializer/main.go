package main

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const AcceptedK8sTypes = `(Deployment|Service)`

var yamlString = `apiVersion: v1
kind: Service
metadata:
  name: my-nginx-svc
  labels:
    app: nginx
spec:
  type: LoadBalancer
  ports:
    - port: 80
  selector:
    app: nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-nginx
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80`

var schemeBuilders = runtime.SchemeBuilder{
	appsv1.AddToScheme,
	corev1.AddToScheme,
}

func main() {
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	_ = schemeBuilders.AddToScheme(scheme)

	// parse  `---` separated yaml file convert it
	// to corresponding k8s objects structure
	acceptedK8sTypes := regexp.MustCompile(AcceptedK8sTypes)
	yamlFiles := strings.Split(yamlString, "---")
	k8sObjects := make([]runtime.Object, 0, len(yamlFiles))

	for _, yamlFile := range yamlFiles {
		if yamlFile == "\n" || yamlFile == "" {
			continue // ignore empty cases
		}

		obj, groupVersionKind, err := codecs.UniversalDeserializer().Decode([]byte(yamlFile), nil, nil)
		if err != nil {
			log.Println(fmt.Sprintf("error while decoding YAML object. %v", err))
			continue
		}

		if !acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
			log.Errorf("unsupported K8s object types: `%s` ", groupVersionKind.String())
			continue
		}
		k8sObjects = append(k8sObjects, obj)
	}

	for _, obj := range k8sObjects {

		//  switch over the type of the object
		switch obj := obj.(type) {
		case *appsv1.Deployment:
			log.Infof("got deployment obj %+v", obj)
		case *corev1.Service:
			log.Infof("got deployment obj %+v", obj)
		default:
			log.Warnf("unsupported object Kind %s", obj.GetObjectKind())
			//o is unknown for us
		}
		log.Infof("==================")
	}

}
