package actions

import (
	"context"
	"fmt"
	"reflect"
	"smiap-k8s/internal/differ"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func NewActionRunner(clientset *kubernetes.Clientset) *actionRunner {
	return &actionRunner{clientset}
}

func (r actionRunner) Create(ctx context.Context, item differ.DiffItem) error {
	switch item.Object.(type) {
	case *apiv1.Namespace:
		return createNamespace(ctx, r.clientset, item)
	case *apiv1.PersistentVolume:
		return createPersistentVolume(ctx, r.clientset, item)
	case *appsv1.StatefulSet:
		return createStatefulSet(ctx, r.clientset, item)
	}

	return fmt.Errorf("cannot create unsupported object: %s", reflect.TypeOf(item.Object))
}

func (r actionRunner) Delete(ctx context.Context, item differ.DiffItem) error {
	switch item.Object.(type) {
	case *apiv1.Namespace:
		return deleteNamespace(ctx, r.clientset, item)
	case *apiv1.PersistentVolume:
		return deletePersistentVolume(ctx, r.clientset, item)
	case *appsv1.StatefulSet:
		return deleteStatefulSet(ctx, r.clientset, item)
	}

	return fmt.Errorf("cannot delete unsupported object: %s", reflect.TypeOf(item.Object))
}

func (r actionRunner) Update(ctx context.Context, item differ.DiffItem) error {
	switch item.Object.(type) {
	case *apiv1.Namespace:
		return updateNamespace(ctx, r.clientset, item)
	case *apiv1.PersistentVolume:
		return updatePersistentVolume(ctx, r.clientset, item)
	case *appsv1.StatefulSet:
		return updateStatefulSet(ctx, r.clientset, item)
	}

	return fmt.Errorf("cannot update unsupported object: %s", reflect.TypeOf(item.Object))
}

type actionRunner struct {
	clientset *kubernetes.Clientset
}

func createNamespace(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	ns := item.Object.(*apiv1.Namespace)

	_, err := clientset.CoreV1().Namespaces().Create(ctx, ns, v1.CreateOptions{})
	return err
}

func deleteNamespace(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	ns := item.Object.(*apiv1.Namespace)

	return clientset.CoreV1().Namespaces().Delete(ctx, ns.Name, v1.DeleteOptions{})
}

func updateNamespace(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	ns := item.Object.(*apiv1.Namespace)

	_, err := clientset.CoreV1().Namespaces().Update(ctx, ns, v1.UpdateOptions{})
	return err
}

func createPersistentVolume(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	pv := item.Object.(*apiv1.PersistentVolume)

	_, err := clientset.CoreV1().PersistentVolumes().Create(ctx, pv, v1.CreateOptions{})
	return err
}
func deletePersistentVolume(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	pv := item.Object.(*apiv1.PersistentVolume)

	return clientset.CoreV1().PersistentVolumes().Delete(ctx, pv.Name, v1.DeleteOptions{})
}
func updatePersistentVolume(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	pv := item.Object.(*apiv1.PersistentVolume)

	_, err := clientset.CoreV1().PersistentVolumes().Update(ctx, pv, v1.UpdateOptions{})
	return err
}

func createStatefulSet(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	sts := item.Object.(*appsv1.StatefulSet)

	_, err := clientset.AppsV1().StatefulSets(sts.Namespace).Create(ctx, sts, v1.CreateOptions{})
	return err
}

func deleteStatefulSet(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	sts := item.Object.(*appsv1.StatefulSet)

	return clientset.AppsV1().StatefulSets(sts.Namespace).Delete(ctx, sts.Name, v1.DeleteOptions{})
}

func updateStatefulSet(ctx context.Context, clientset *kubernetes.Clientset, item differ.DiffItem) error {
	sts := item.Object.(*appsv1.StatefulSet)

	_, err := clientset.AppsV1().StatefulSets(sts.Namespace).Update(ctx, sts, v1.UpdateOptions{})
	return err
}
