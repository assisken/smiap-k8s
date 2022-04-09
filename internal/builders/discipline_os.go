package builders

import (
	"fmt"
	"smiap-k8s/internal/smiap"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OsBuilder struct {
	container smiap.Container
}

func NewOsBuilder(container smiap.Container) *OsBuilder {
	return &OsBuilder{container: container}
}

func (d OsBuilder) Namespace() apiv1.Namespace {
	return apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "os",
			Labels: map[string]string{
				"app": "smiap-k8s",
				"id":  "os",
			},
		},
	}
}

func (d OsBuilder) StatefulSet() appsv1.StatefulSet {
	volumeMode := apiv1.PersistentVolumeFilesystem
	terminationGracePeriodSeconds := int64(10)
	storageClassName := ""

	return appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.container.Name,
			Namespace: "os",
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
							Name:    d.container.Name,
							Image:   "ubuntu:22.04",
							Command: []string{"sleep", "infinity"},
							Ports: []apiv1.ContainerPort{
								{Name: "ssh", ContainerPort: 20},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{Name: d.container.Name, ReadOnly: false, MountPath: "/home"},
							},
							ImagePullPolicy: apiv1.PullAlways,
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceMemory: *resource.NewScaledQuantity(
										int64(d.container.MemoryGb), resource.Giga,
									),
									apiv1.ResourceCPU: *resource.NewScaledQuantity(
										int64(d.container.Cores)*1000, resource.Milli,
									),
								},
								Limits: apiv1.ResourceList{
									apiv1.ResourceMemory: *resource.NewScaledQuantity(
										int64(d.container.MemoryGb)*2, resource.Giga,
									),
									apiv1.ResourceCPU: *resource.NewScaledQuantity(
										int64(d.container.Cores)*2000, resource.Milli,
									),
								},
							},
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
						StorageClassName: &storageClassName,
						AccessModes:      []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
						VolumeName:       d.container.Name,
						VolumeMode:       &volumeMode,
						Resources: apiv1.ResourceRequirements{
							Requests: apiv1.ResourceList{
								apiv1.ResourceStorage: *resource.NewScaledQuantity(
									int64(d.container.PartitionSizeGb), resource.Giga,
								),
							},
						},
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
			Name:      d.container.Name,
			Namespace: "os",
			Labels: map[string]string{
				"app":           "smiap-k8s",
				"id":            strconv.Itoa(d.container.Id),
				"do_not_remove": strconv.FormatBool(d.container.DoNotRemove),
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			StorageClassName: "",
			Capacity: map[apiv1.ResourceName]resource.Quantity{
				apiv1.ResourceStorage: *resource.NewScaledQuantity(
					int64(d.container.PartitionSizeGb), resource.Giga,
				),
			},
			AccessModes:                   []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimRecycle,
			VolumeMode:                    &mode,
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: fmt.Sprintf("/tmp/%s", d.container.Name),
				},
			},
		},
	}
}

func (d OsBuilder) Service() apiv1.Service {
	panic("implement me")
}
