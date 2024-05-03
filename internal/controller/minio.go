package controller

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

func (r *DbManageReconciler) InitialMinioClient(dbManage *operatorcodehorsecomv1beta1.DbManage) (*minio.Client, error) {
	minioClient, err := minio.New(dbManage.Spec.Destination.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(dbManage.Spec.Destination.AccessKey, dbManage.Spec.Destination.AccessSecret, ""),
		Secure: false,
	})
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("初始化MinIO客户端失败, 原因: %s", err.Error())
		return nil, err
	}
	return minioClient, nil
}
