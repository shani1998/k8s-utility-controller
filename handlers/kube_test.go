package handlers

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

var fakeError = errors.New("failed to connect with api server 0.0.0.0:6443")

var fakeDeploymentSpec = &appv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{Name: testServiceName, Namespace: defaultNS, Labels: map[string]string{appGroup: testAppGrp}},
	Spec:       appv1.DeploymentSpec{},
	Status:     appv1.DeploymentStatus{ReadyReplicas: 1},
}

func createFakeDeployment() {
	_, _ = kubeClient.AppsV1().Deployments(defaultNS).Create(context.TODO(), fakeDeploymentSpec, metav1.CreateOptions{})
}

func deleteFakeDeployment() {
	_ = kubeClient.AppsV1().Deployments(defaultNS).Delete(context.TODO(), testServiceName, metav1.DeleteOptions{})
}

var errorReaction = func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
	return true, &appv1.DeploymentList{}, fakeError
}

func TestListDeployments(t *testing.T) {
	// create the fake client.
	kubeClient = fake.NewSimpleClientset()

	tests := []struct {
		name    string
		ctx     context.Context
		opts    metav1.ListOptions
		want    *appv1.DeploymentList
		wantErr error
	}{
		{
			name:    "failure, api service unreachable",
			ctx:     context.TODO(),
			opts:    metav1.ListOptions{},
			want:    &appv1.DeploymentList{},
			wantErr: fakeError,
		},
		{
			name: "success, get one deployment",
			ctx:  context.Background(),
			opts: metav1.ListOptions{},
			want: &appv1.DeploymentList{
				Items: []appv1.Deployment{*fakeDeploymentSpec},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare test scenario
			if strings.Contains(tt.name, "failure") {
				// return a fake error getting the deployment  list
				kubeClient.(*fake.Clientset).Fake.PrependReactor("list", "deployments", errorReaction)
				defer func() {
					// reset client set from error state
					kubeClient = fake.NewSimpleClientset()
				}()
			}
			if strings.Contains(tt.name, "success") {
				// deploy one fake service
				createFakeDeployment()
				// clean fake deployment
				defer deleteFakeDeployment()
			}

			got, err := ListDeployments(tt.ctx, tt.opts)
			if err != tt.wantErr {
				t.Errorf("ListDeployments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Items, tt.want.Items) {
				t.Errorf(" ListDeployments() \n got = %v,\n want %v", got, tt.want)
			}
		})
	}
}
