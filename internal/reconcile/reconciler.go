package reconcile

import (
	"context"
	"fmt"
	"smiap-k8s/internal/builders"
	"smiap-k8s/internal/differ"
	"smiap-k8s/internal/smiap"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Reconciler interface {
	Reconcile(ctx context.Context) error
}

type Client interface {
	Create(ctx context.Context, obj interface{}, opts v1.CreateOptions) (interface{}, error)
	Update(ctx context.Context, obj interface{}, opts v1.UpdateOptions) (interface{}, error)
	Delete(ctx context.Context, obj interface{}, opts v1.DeleteOptions) (interface{}, error)
	List(ctx context.Context, opts interface{}) ([]interface{}, error)
}

func NewReconciler(fetcher *smiap.Fetcher, client *kubernetes.Clientset, actionRunner ActionRunner) reconciler {
	return reconciler{fetcher, client, actionRunner}
}

func (r reconciler) Reconcile(ctx context.Context) error {
	containers, err := r.fetcher.FetchContainers()
	if err != nil {
		return err
	}

	fmt.Println(containers)

	wantNamespaces := map[string]differ.DiffObject{}
	wantStatefulSets := map[string]differ.DiffObject{}
	wantPersistentVolumes := map[string]differ.DiffObject{}

	haveNamespaces := map[string]differ.DiffObject{}
	haveStatefulSets := map[string]differ.DiffObject{}
	havePersistentVolumes := map[string]differ.DiffObject{}

	options := metav1.ListOptions{
		LabelSelector: "app=smiap-k8s",
	}

	for _, container := range containers {
		builder := builders.NewOsBuilder(container)

		ns := builder.Namespace()
		wantNamespaces[ns.Name] = &ns

		namespaceList, err := r.client.CoreV1().Namespaces().List(ctx, options)
		if err != nil {
			return errors.Errorf("At getting namespaces: %v", err)
		}

		for _, ns := range namespaceList.Items {
			haveNamespaces[ns.Name] = &ns
		}

		sts := builder.StatefulSet()
		wantStatefulSets[sts.Name] = &sts

		statefulSetList, err := r.client.AppsV1().StatefulSets(sts.Namespace).List(ctx, options)
		if err != nil {
			return errors.Errorf("At getting statefulsets: %v", err)
		}

		for _, sts := range statefulSetList.Items {
			haveStatefulSets[sts.Name] = &sts
		}

		pv := builder.PersistentVolume()
		wantPersistentVolumes[pv.Name] = &pv

		persistentVolumeList, err := r.client.CoreV1().PersistentVolumes().List(ctx, options)
		if err != nil {
			return errors.Errorf("At getting persistent volumes: %v", err)
		}

		for _, pv := range persistentVolumeList.Items {
			havePersistentVolumes[pv.Name] = &pv
		}
	}

	wants := []map[string]differ.DiffObject{wantNamespaces, wantPersistentVolumes, wantStatefulSets}
	haves := []map[string]differ.DiffObject{haveNamespaces, havePersistentVolumes, haveStatefulSets}

	for i := range wants {
		toCreate, toDelete, toUpdate, err := differ.Diff(wants[i], haves[i])
		if err != nil {
			return err
		}

		for _, create := range toCreate {
			fmt.Printf("Creating %s. Reason: %s\n", create.Id, create.Reason)

			if err := r.actionRunner.Create(ctx, create); err != nil {
				return err
			}
		}

		for _, delete := range toDelete {
			fmt.Printf("Deleting %s. Reason: %s\n", delete.Id, delete.Reason)

			if err := r.actionRunner.Delete(ctx, delete); err != nil {
				return err
			}
		}
		for _, update := range toUpdate {
			fmt.Printf("Modifying %s. Reason: %s\n", update.Id, update.Reason)

			if err := r.actionRunner.Update(ctx, update); err != nil {
				return err
			}
		}
	}

	return nil
}

type reconciler struct {
	fetcher      *smiap.Fetcher
	client       *kubernetes.Clientset
	actionRunner ActionRunner
}
