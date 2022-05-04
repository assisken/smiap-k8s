package domain

import "k8s.io/apimachinery/pkg/runtime/schema"

type Resource schema.GroupVersionResource

var (
	Namespace        Resource = Resource{Group: "", Version: "v1", Resource: "namespaces"}
	PersistentVolume Resource = Resource{Group: "", Version: "v1", Resource: "persistentvolumes"}
	StatefulSet      Resource = Resource{Group: "apps", Version: "v1", Resource: "statefulsets"}
)
