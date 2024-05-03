/*
Copyright 2024.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DbManageSpec defines the desired state of DbManage
type DbManageSpec struct {
	// 是否开启任务
	Enable bool `json:"enable"`
	// 任务类型标识(0: 备份任务 1: 巡检任务 ...)
	Flag int `json:"flag"`
	// 任务开始时间
	StartTime string `json:"startTime"`
	// 任务间隔时间(分钟)
	Period int `json:"period"`
	// 数据库源地址
	Origin *Origin `json:"origin"`
	// 目标地址(MinIO)
	Destination *Destination `json:"destination"`
}

// 数据库源地址
type Origin struct {
	// 数据库IP地址
	Host string `json:"host"`
	// 数据库端口
	Port int32 `json:"port"`
	// 数据库用户名
	Username string `json:"username"`
	// 数据库密码
	Password string `json:"password"`
}

// 目标地址(MinIO)
type Destination struct {
	// MinIO地址
	Endpoint string `json:"endpoint"`
	// MinIO访问的key
	AccessKey string `json:"accessKey"`
	// MinIO访问的密钥
	AccessSecret string `json:"accessSecret"`
	// MinIO令牌桶名称
	BucketName string `json:"bucketName"`
}

// DbManageStatus defines the observed state of DbManage
type DbManageStatus struct {
	// 任务状态是否为活动
	Active bool `json:"active"`
	// 任务下次启动时间(时间戳)
	NextTime int `json:"nextTime"`
	// 最后一次任务运行情况
	LastTaskResult string `json:"lastTaskResult"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DbManage is the Schema for the dbmanages API
type DbManage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DbManageSpec   `json:"spec,omitempty"`
	Status DbManageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DbManageList contains a list of DbManage
type DbManageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DbManage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DbManage{}, &DbManageList{})
}
