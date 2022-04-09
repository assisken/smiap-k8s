package builders

import (
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

type Builder interface {
	StatefulSet() v1.StatefulSet
	PersistentVolumeClaim() apiv1.PersistentVolumeClaim
	PersistentVolume() apiv1.PersistentVolume
	Service() apiv1.Service
}
