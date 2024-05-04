package controller

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/minio/minio-go/v7"
	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

// 数据库备份任务
func (r *DbManageReconciler) DbCheckTask(dbManage *operatorcodehorsecomv1beta1.DbManage) error {
	// 数据库状态检查
	dbCheckDir := "/tmp/dbcheck"
	_, err := os.Stat(dbCheckDir)
	if err != nil {
		if errx := os.MkdirAll(dbCheckDir, 0700); errx == nil {
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s文件夹创建成功", dbCheckDir)
		} else {
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s文件夹创建失败, 原因: %s", dbCheckDir, errx.Error())
			return errx
		}
	}
	fileName := time.Now().Format("0102150405") + ".txt"
	dbcheckCmd := fmt.Sprintf(`mysql -u%s -p%s -h%s -P%d -e"show global status\G" > %s/%s`,
		dbManage.Spec.Origin.Username,
		dbManage.Spec.Origin.Password,
		dbManage.Spec.Origin.Host,
		dbManage.Spec.Origin.Port,
		dbCheckDir,
		fileName)
	_, err = exec.Command("bash", "-c", dbcheckCmd).Output()
	if err != nil {
		return err
	}
	// 同步备份文件到MinIO
	minioClient, err := r.InitialMinioClient(dbManage)
	if err != nil {
		return err
	}
	f, err := os.Open(fmt.Sprintf("%s/%s", dbCheckDir, fileName))
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("打开数据库状态检查文件%s失败", fileName)
		return err
	}
	_, err = minioClient.PutObject(
		context.TODO(),
		dbManage.Spec.Destination.BucketName,
		fmt.Sprintf("%s/%s", dbCheckDir, fileName),
		f,
		-1,
		minio.PutObjectOptions{})
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("上传数据库状态检查文件%s到minio失败", fileName)
		return err
	}
	return nil
}
