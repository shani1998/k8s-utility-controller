package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shani1998/k8s-utility-controller/models"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	testServiceName = "fake-test-service"
	testAppGrp      = "alpha"
)

func TestGetServices(t *testing.T) {
	// create the fake client.
	kubeClient = fake.NewSimpleClientset()
	req := httptest.NewRequest("GET", "http://test-service.com/services", nil)
	go func() {
		for {
			// consume test errors
			<-HealthChan
		}
	}()

	tests := []struct {
		name     string
		rw       http.ResponseWriter
		rq       *http.Request
		params   httprouter.Params
		want     []models.Service
		wantErr  error
		wantCode int
	}{
		{
			name:     "Failure, api service unavailable",
			rq:       req,
			want:     []models.Service{},
			wantErr:  errors.New("failed to list services"),
			wantCode: http.StatusServiceUnavailable,
		},
		{
			name:     "Success, no services deployed",
			rq:       req,
			want:     []models.Service{},
			wantErr:  nil,
			wantCode: http.StatusOK,
		},
		{
			name: "Success, one service deployed",
			rq:   req,
			want: []models.Service{
				{
					Name:             "fake-test-service",
					ApplicationGroup: "alpha",
					RunningPodsCount: 1,
				}},
			wantErr:  nil,
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare test scenario
			if strings.Contains(tt.name, "Failure") {
				// return a fake error getting the deployment  list
				kubeClient.(*fake.Clientset).Fake.PrependReactor("list", "deployments", errorReaction)
				defer func() {
					// reset client set from error state
					kubeClient = fake.NewSimpleClientset()
				}()
			}
			if strings.Contains(tt.name, "one service deployed") {
				// deploy one fake service
				createFakeDeployment()
				defer deleteFakeDeployment()
			}

			w := httptest.NewRecorder()
			GetServices(w, tt.rq, tt.params)

			if tt.wantErr != nil {
				// assert on expected status code
				if tt.wantCode != w.Code {
					t.Errorf("mismatched status code: want=%v, got=%v", tt.wantCode, w.Code)
				}
				// assert on expected error value
				if !strings.Contains(w.Body.String(), tt.wantErr.Error()) {
					t.Errorf("mismatched error: want=%v, got=%v", tt.wantErr, w.Body)
				}
				// avoid unmarshalling in case of errors
				return
			}

			var gotResp []models.Service
			if err := json.Unmarshal(w.Body.Bytes(), &gotResp); err != nil {
				t.Errorf("failed to unmarshal response %v", err)
			}
			// assert on total number of response received
			if len(tt.want) != len(gotResp) {
				t.Errorf("mismatched response count: want=%v, got=%v", len(tt.want), len(gotResp))
				return
			}
			// assert on value received in response
			if len(gotResp) > 0 && !reflect.DeepEqual(tt.want, gotResp) {
				t.Errorf("want %v, got %v", tt.want, gotResp)
			}
		})
	}
}

func TestGetServicesByAppLabel(t *testing.T) {
	// create the fake client.
	kubeClient = fake.NewSimpleClientset()
	req := httptest.NewRequest("GET", "http://test.service.com/services/invalid<label>", nil)
	testParams := httprouter.Params{httprouter.Param{Key: appGroup, Value: "alpha"}}
	go func() {
		for {
			// consume test errors
			<-HealthChan
		}
	}()

	tests := []struct {
		name     string
		rw       http.ResponseWriter
		rq       *http.Request
		params   httprouter.Params
		want     []models.Service
		wantErr  error
		wantCode int
	}{
		{
			name:     "Failure, api service unavailable",
			rq:       req,
			params:   testParams,
			want:     []models.Service{},
			wantErr:  errors.New("failed to list services"),
			wantCode: http.StatusServiceUnavailable,
		},
		{
			name:     "Success, no services deployed",
			rq:       req,
			params:   testParams,
			want:     []models.Service{},
			wantErr:  nil,
			wantCode: http.StatusOK,
		},
		{
			name:     "Success, one service deployed and passed invalid label",
			rq:       req,
			params:   httprouter.Params{httprouter.Param{Key: appGroup, Value: "wrongLabel"}},
			want:     []models.Service{},
			wantErr:  nil,
			wantCode: http.StatusServiceUnavailable,
		},
		{
			name:   "Success, one service deployed and passed valid label",
			rq:     req,
			params: testParams,
			want: []models.Service{
				{
					Name:             "fake-test-service",
					ApplicationGroup: "alpha",
					RunningPodsCount: 1,
				}},
			wantErr:  nil,
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//prepare test scenario
			if strings.Contains(tt.name, "Failure") {
				// return a fake error getting the deployment  list
				kubeClient.(*fake.Clientset).Fake.PrependReactor("list", "deployments", errorReaction)
				defer func() {
					// reset client set from error state
					kubeClient = fake.NewSimpleClientset()
				}()
			}
			if strings.Contains(tt.name, "one service deployed") {
				// deploy one fake service
				createFakeDeployment()
				// filter by label alpha
				tt.rq = httptest.NewRequest("GET", "http://test-service.com/services/alpha", nil)
				defer deleteFakeDeployment()
			}

			w := httptest.NewRecorder()
			GetServicesByAppLabel(w, tt.rq, tt.params)

			if tt.wantErr != nil {
				// assert on expected status code
				if tt.wantCode != w.Code {
					t.Errorf("mismatched status code: want=%v, got=%v", tt.wantCode, w.Code)
				}
				// assert on expected error value
				if !strings.Contains(w.Body.String(), tt.wantErr.Error()) {
					t.Errorf("mismatched error: want=%v, got=%v", tt.wantErr, w.Body)
				}
				// avoid unmarshalling in case of errors
				return
			}

			var gotResp []models.Service
			if err := json.Unmarshal(w.Body.Bytes(), &gotResp); err != nil {
				t.Errorf("failed to unmarshal response %v", err)
			}
			// assert on total number of response received
			if len(tt.want) != len(gotResp) {
				t.Errorf("mismatched response count: want=%v, got=%v", len(tt.want), len(gotResp))
				return
			}
			// assert on value received in response
			if len(gotResp) > 0 && !reflect.DeepEqual(tt.want, gotResp) {
				t.Errorf("want %v, got %v", tt.want, gotResp)
			}
		})
	}
}
