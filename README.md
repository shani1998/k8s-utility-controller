# k8s-utility-controller
> This repository implements the rest endpoints to fetch all apps deployed as a deployment object on the current k8s cluster.
### API
#### /services
* `GET` : Get all services contains number of pods running in the cluster in namespace default per service and per application group.

Example:
``` sh
$ curl -X GET -H "Content-type: application/json" -H "Accept: application/json" http://localhost:8080/services
[
  {
    "name": "<service>",
    "applicationGroup": "alpha",
    "runningPodsCount": 2
  },
  {
    "name": "<service>",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  }
  ...
]
```
#### /services/:title
* `GET` : Get all services by application group contains number of running pods in the cluster in namespace `default` that are part of the same `applicationGroup`

Example:

```sh
$ curl -X GET -H "Content-type: application/json" -H "Accept: application/json" http://localhost:8080/services/alpha
GET `/services/alpha`
[
  {
    "name": "<service>",
    "applicationGroup": "alpha",
    "runningPodsCount": 2
  }
]
```
## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites
- Docker Engine (preferably 20.0.1+) is required to run `make docker-build`/`make docker-push`.
- If you are running locally using `make run-local` then make sure that kubeconfig file exist at location `~/.kube/config` and context is set to the Kubernetes cluster that you want to work on.

### Usage

#### Download:
```sh
$ git clone https://github.com/shani1998/k8s-utility-controller.git
```

#### How to run tests
To run all tests
```sh
$ make test
```

To build and push docker image

```sh
# change image registry variable in Makefile to where you have the permission.
$ make docker-build
$ make docker-push
```

To run and test locally
```sh
$ make run-local #kubeconfig file should present at ~/.kube/config
$ curl localhost:8080/services
$ curl localhost:8080/services/alpha
```

To deploy in kubernetes cluster and test
```sh
$ k apply -f deploy/
$ curl <Node-IP>:<NODE-PORT>/services
$ curl <Node-IP>:<NODE-PORT>/services/alpha
# grep NodePort using
# kubectl get service  k8s-utility-controller -oyaml | grep -i nodeport
# get the nodeIP where pod deployed
```
