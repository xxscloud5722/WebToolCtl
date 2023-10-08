package console

import v1 "k8s.io/api/core/v1"

func Copy(sourceConfigPath, sourceNamespace, targetConfigPath, targetNamespace string) {
	// 备份
	backupPath, err := Backup(sourceConfigPath, "./copy.zip", func(namespace v1.Namespace) bool {
		return namespace.ObjectMeta.Name == sourceNamespace
	})
	if err != nil {
		return
	}
	// 恢复
	err = Restore(targetConfigPath, targetNamespace, *backupPath)
	if err != nil {
		return
	}
}
