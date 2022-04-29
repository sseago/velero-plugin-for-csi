package backup

import (
	volumesnapmoverv1alpha1 "github.com/konveyor/volume-snapshot-mover/api/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero-plugin-for-csi/internal/util"
	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// DataMoverBackupBackupItemAction is a backup item action plugin to backup
// DataMoverBackup objects using Velero
type DataMoverBackupItemAction struct {
	Log logrus.FieldLogger
}

// AppliesTo returns information indicating that the DataMoverBackupBackupItemAction should be invoked to backup DataMoverBackups.
func (p *DataMoverBackupItemAction) AppliesTo() (velero.ResourceSelector, error) {
	p.Log.Info("DataMoverBackupItemAction AppliesTo")

	return velero.ResourceSelector{
		IncludedResources: []string{"datamoverbackups.pvc.oadp.openshift.io"},
	}, nil
}

// Execute backs up a DataMoverBackup object with a completely filled status
func (p *DataMoverBackupItemAction) Execute(item runtime.Unstructured, backup *velerov1api.Backup) (runtime.Unstructured, []velero.ResourceIdentifier, error) {
	p.Log.Infof("Executing DataMoverBackupItemAction")

	var dmb volumesnapmoverv1alpha1.DataMoverBackup
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.UnstructuredContent(), &dmb); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	dmbNew, err := util.GetDataMoverbackupWithCompletedStatus(dmb.Namespace, dmb.Name, p.Log)

	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	p.Log.Infof("Value of dmbNew is : %v", dmbNew)

	dmbMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&dmbNew)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	p.Log.Infof("Returning DMB map : %v", dmbMap)
	return &unstructured.Unstructured{Object: dmbMap}, nil, nil
}
