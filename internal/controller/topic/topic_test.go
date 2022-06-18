/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package topic

import (
	"context"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/scalvetr/poc-crossplane-provider/apis/objects/v1alpha1"
	"github.com/scalvetr/poc-crossplane-provider/internal/client/poc"
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

// Unlike many Kubernetes projects Crossplane does not use third party testing
// libraries, per the common Go test review comments. Crossplane encourages the
// use of table driven unit tests. The tests of the crossplane-runtime project
// are representative of the testing style Crossplane encourages.
//
// https://github.com/golang/go/wiki/TestComments
// https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md#contributing-code
func instance(name string, partitions int, status string) *v1alpha1.Topic {
	return &v1alpha1.Topic{
		Spec: v1alpha1.TopicSpec{
			ForProvider: v1alpha1.TopicParameters{
				Name:       name,
				Partitions: partitions,
			},
		},
		Status: v1alpha1.TopicStatus{
			AtProvider: v1alpha1.TopicObservation{
				Status: status,
			},
		},
	}
}
func observation(exists bool) managed.ExternalObservation {
	return managed.ExternalObservation{ResourceExists: exists, ResourceUpToDate: exists, ConnectionDetails: managed.ConnectionDetails{}}
}
func buildClient() poc.POCService {
	service, _ := poc.NewService("tenant1", "http://localhost:8080/v1", []byte("ABC"))
	return service
}

func TestObserve(t *testing.T) {
	type fields struct {
		service poc.POCService
	}

	type args struct {
		ctx context.Context
		mg  resource.Managed
	}

	type want struct {
		o   managed.ExternalObservation
		err error
	}

	cases := map[string]struct {
		reason string
		fields fields
		args   args
		want   want
	}{
		"Created": {
			reason: "",
			fields: fields{
				service: buildClient(),
			},
			args: args{
				mg: instance("topic1", 4, "CREATING"),
			},
			want: want{
				o:   observation(true),
				err: nil,
			},
		},
		"Invalid": {
			reason: "",
			fields: fields{
				service: buildClient(),
			},
			args: args{
				mg: instance("topic1", 11, "CREATING"),
			},
			want: want{
				o:   observation(false),
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{service: tc.fields.service}
			got, err := e.Observe(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}
