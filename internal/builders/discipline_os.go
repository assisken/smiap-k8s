package builders

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"smiap-k8s/internal/smiap"
	"strconv"
)

type OsBuilder struct {
	container smiap.Container
}

func NewOsBuilder(container smiap.Container) *OsBuilder {
	return &OsBuilder{container: container}
}

func (d OsBuilder) StatefulSet() appsv1.StatefulSet {
	volumeMode := apiv1.PersistentVolumeFilesystem
	terminationGracePeriodSeconds := int64(10)

	return appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.container.Name,
			Labels: map[string]string{
				"app":           "smiap-k8s",
				"id":            strconv.Itoa(d.container.Id),
				"do_not_remove": strconv.FormatBool(d.container.DoNotRemove),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "smiap-k8s",
					"id":  strconv.Itoa(d.container.Id),
				},
			},
			ServiceName:     d.container.Name,
			MinReadySeconds: 60,
			PersistentVolumeClaimRetentionPolicy: &appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy{
				WhenDeleted: appsv1.DeletePersistentVolumeClaimRetentionPolicyType,
				WhenScaled:  appsv1.DeletePersistentVolumeClaimRetentionPolicyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: d.container.Name,
					Labels: map[string]string{
						"app":           "smiap-k8s",
						"id":            strconv.Itoa(d.container.Id),
						"do_not_remove": strconv.FormatBool(d.container.DoNotRemove),
					},
				},
				Spec: apiv1.PodSpec{
					RestartPolicy:                 apiv1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Hostname:                      d.container.Name,
					Containers: []apiv1.Container{
						{
							Name:  d.container.Name,
							Image: "ubuntu:22.04",
							Ports: []apiv1.ContainerPort{
								{Name: "SSH", ContainerPort: 20},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{Name: d.container.Name, ReadOnly: false, MountPath: "/"},
							},
							ImagePullPolicy: apiv1.PullAlways,
						},
					},
				},
			},
			VolumeClaimTemplates: []apiv1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: d.container.Name,
						Labels: map[string]string{
							"app":           "smiap-k8s",
							"id":            strconv.Itoa(d.container.Id),
							"do_not_remove": strconv.FormatBool(d.container.DoNotRemove),
						},
					},
					Spec: apiv1.PersistentVolumeClaimSpec{
						AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
						Resources:   apiv1.ResourceRequirements{},
						VolumeName:  d.container.Name,
						VolumeMode:  &volumeMode,
					},
				},
			},
		},
	}
}

func (d OsBuilder) PersistentVolumeClaim() apiv1.PersistentVolumeClaim {
	panic("implement me")
}

func (d OsBuilder) PersistentVolume() apiv1.PersistentVolume {
	mode := apiv1.PersistentVolumeFilesystem

	return apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.container.Name,
			Labels: map[string]string{
				"app":           "smiap-k8s",
				"id":            strconv.Itoa(d.container.Id),
				"do_not_remove": strconv.FormatBool(d.container.DoNotRemove),
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			Capacity: map[apiv1.ResourceName]resource.Quantity{
				apiv1.ResourceRequestsStorage: *resource.NewScaledQuantity(
					int64(d.container.PartitionSizeGb), resource.Giga,
				),
			},
			PersistentVolumeSource:        apiv1.PersistentVolumeSource{},
			AccessModes:                   []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimRecycle,
			VolumeMode:                    &mode,
		},
	}
}

func (d OsBuilder) Service() apiv1.Service {
	panic("implement me")
}
