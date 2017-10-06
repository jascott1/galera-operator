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

package v1beta2

import "errors"

type BackupStorageType string

const (
	BackupStorageTypeDefault          = ""
	BackupStorageTypePersistentVolume = "PersistentVolume"
	BackupStorageTypeS3               = "S3"
	BackupStorageTypeABS              = "ABS"

	AWSSecretCredentialsFileName = "credentials"
	AWSSecretConfigFileName      = "config"

	// ABSStorageAccount defines the key for the Azure Storage Account value in the ABS Kubernetes secret
	ABSStorageAccount = "storage-account"
	// ABSStorageKey defines the key for the Azure Storage Key value in the ABS Kubernetes secret
	ABSStorageKey = "storage-key"
)

var errPVZeroSize = errors.New("PV backup should not have 0 size volume")

type BackupPolicy struct {
	// Pod defines the policy to create the backup pod.
	Pod *PodPolicy `json:"pod,omitempty"`

	// StorageType specifies the type of storage device to store backup files.
	// If it's not set by user, the default is "PersistentVolume".
	StorageType BackupStorageType `json:"storageType"`

	StorageSource `json:",inline"`

	// BackupIntervalInSecond specifies the interval between two backups.
	// The default interval is 1800 seconds.
	BackupIntervalInSecond int `json:"backupIntervalInSecond"`

	// If greater than 0, MaxBackups is the maximum number of backup files to retain.
	// If equal to 0, it means unlimited backups.
	// Otherwise, it is invalid.
	MaxBackups int `json:"maxBackups"`

	// AutoDelete tells whether to cleanup backup data if cluster is deleted.
	// By default (false), operator will keep the backup data.
	AutoDelete bool `json:"autoDelete"`
}

func (bp *BackupPolicy) Validate() error {
	if bp.MaxBackups < 0 {
		return errors.New("MaxBackups value should be >= 0")
	}
	if bp.StorageType == BackupStorageTypePersistentVolume {
		if pv := bp.StorageSource.PV; pv == nil || pv.VolumeSizeInMB <= 0 {
			return errPVZeroSize
		}
	}
	return nil
}

type StorageSource struct {
	// PV represents a Persistent Volume resource, operator will claim the
	// required size before creating the etcd cluster for backup purpose.
	// If the snapshot size is larger than the size specified operator would
	// kill the cluster and report failure condition in status.
	PV *PVSource `json:"pv,omitempty"`
	S3 *S3Source `json:"s3,omitempty"`
	// ABS represents an Azure Blob Storage resource for storing etcd backups
	ABS *ABSSource `json:"abs,omitempty"`
}

// TODO: support per cluster S3 Source configuration.
type S3Source struct {
	// The name of the AWS S3 bucket to store backups in.
	//
	// S3Bucket overwrites the default etcd operator wide bucket.
	S3Bucket string `json:"s3Bucket,omitempty"`

	// Prefix is the S3 prefix used to prefix the bucket path.
	// It's the prefix at the beginning.
	// After that, it will have version and cluster specific paths.
	Prefix string `json:"prefix,omitempty"`

	// The name of the secret object that stores the AWS credential and config files.
	// The file name of the credential MUST be 'credentials'.
	// The file name of the config MUST be 'config'.
	// The profile to use in both files will be 'default'.
	//
	// AWSSecret overwrites the default etcd operator wide AWS credential and config.
	AWSSecret string `json:"awsSecret,omitempty"`
}

// ABSSource represents an Azure Blob Storage (ABS) backup storage source
type ABSSource struct {
	// ABSContainer is the name of the ABS container to store backups in.
	ABSContainer string `json:"absContainer,omitempty"`

	// ABSSecret is the name of the secret object that stores the ABS credentials.
	//
	// Within the secret object, the following fields MUST be provided:
	// 'storage-account' holding the Azure Storage account name
	// 'storage-key' holding the Azure Storage account key
	ABSSecret string `json:"absSecret,omitempty"`
}

type BackupServiceStatus struct {
	// RecentBackup is status of the most recent backup created by
	// the backup service
	RecentBackup *BackupStatus `json:"recentBackup,omitempty"`

	// Backups is the totoal number of existing backups
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
