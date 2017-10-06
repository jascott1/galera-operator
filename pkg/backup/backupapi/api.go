// Copyright 2016 The etcd-operator Authors
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

package backupapi

import "path"

const (
	APIV1 = "/v1"
	// S3V1 indicates the version 1 of
	// S3 backup format: <s3Bucket>/<s3Prefix>/"v1"/<namespace>/<clusterName>
	S3V1 = "v1"
)

type ServiceStatus struct {
	// RecentBackup is status of the most recent backup created by
	// the backup service
	RecentBackup *BackupStatus `json:"recentBackup,omitempty"`

	// Backups is the totoal number of existing backups.
	Backups int `json:"backups"`

	// BackupSize is the total size of existing backups in MB.
	BackupSize float64 `json:"backupSize"`
}

type BackupStatus struct {
	// Creation time of the backup.
	CreationTime string `json:"creationTime"`

	// Size is the size of the backup in MB.
	Size float64 `json:"size"`

	// Revision is the revision of the backup.
	Revision int64 `json:"revision"`

	// Version is the version of the backup cluster.
	Version string `json:"version"`

	// TimeTookInSecond is the total time took to create the backup.
	TimeTookInSecond int `json:"timeTookInSecond"`
}

// ToS3Prefix concatenates s3Prefix, S3V1, namespace, clusterName to a single s3 prefix.
// the concatenated prefix determines the location of S3 backup files.
func ToS3Prefix(s3Prefix, namespace, clusterName string) string {
	return path.Join(s3Prefix, S3V1, namespace, clusterName)
}
