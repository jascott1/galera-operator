// Copyright 2017 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2eutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	api "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
	"github.com/coreos/etcd-operator/pkg/client"
	"github.com/coreos/etcd-operator/pkg/util/k8sutil"
	"github.com/coreos/etcd-operator/pkg/util/retryutil"

	"github.com/aws/aws-sdk-go/service/s3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

type StorageCheckerOptions struct {
	S3Cli          *s3.S3
	S3Bucket       string
	DeletedFromAPI bool
}

func CreateCluster(t *testing.T, crClient client.EtcdClusterCR, namespace string, cl *api.EtcdCluster) (*api.EtcdCluster, error) {
	cl.Namespace = namespace
	res, err := crClient.Create(context.TODO(), cl)
	if err != nil {
		return nil, err
	}
	t.Logf("creating etcd cluster: %s", res.Name)

	return res, nil
}

func UpdateCluster(crClient client.EtcdClusterCR, cl *api.EtcdCluster, maxRetries int, updateFunc k8sutil.EtcdClusterCRUpdateFunc) (*api.EtcdCluster, error) {
	return AtomicUpdateClusterCR(crClient, cl.Name, cl.Namespace, maxRetries, updateFunc)
}

func AtomicUpdateClusterCR(crClient client.EtcdClusterCR, name, namespace string, maxRetries int, updateFunc k8sutil.EtcdClusterCRUpdateFunc) (*api.EtcdCluster, error) {
	result := &api.EtcdCluster{}
	err := retryutil.Retry(1*time.Second, maxRetries, func() (done bool, err error) {
		etcdCluster, err := crClient.Get(context.TODO(), namespace, name)
		if err != nil {
			return false, err
		}

		updateFunc(etcdCluster)

		result, err = crClient.Update(context.TODO(), etcdCluster)
		if err != nil {
			if apierrors.IsConflict(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return result, err
}

func DeleteCluster(t *testing.T, crClient client.EtcdClusterCR, kubeClient kubernetes.Interface, cl *api.EtcdCluster) error {
	t.Logf("deleting etcd cluster: %v", cl.Name)
	err := crClient.Delete(context.TODO(), cl.Namespace, cl.Name)
	if err != nil {
		return err
	}
	return waitResourcesDeleted(t, kubeClient, cl)
}

func DeleteClusterAndBackup(t *testing.T, crClient client.EtcdClusterCR, kubecli kubernetes.Interface, cl *api.EtcdCluster, checkerOpt StorageCheckerOptions) error {
	err := DeleteCluster(t, crClient, kubecli, cl)
	if err != nil {
		return err
	}
	t.Logf("waiting backup deleted of cluster (%v)", cl.Name)
	err = WaitBackupDeleted(kubecli, cl, checkerOpt)
	if err != nil {
		return fmt.Errorf("fail to wait backup deleted: %v", err)
	}
	return nil
}
