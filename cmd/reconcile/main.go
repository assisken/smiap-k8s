package main

import (
	"context"
	"fmt"
	"smiap-k8s/internal/actions"
	"smiap-k8s/internal/reconcile"
	"smiap-k8s/internal/smiap"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fetcher := smiap.NewFetcher()

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	autoRunner := actions.NewActionRunner(clientset)
	reconciler := reconcile.NewReconciler(fetcher, clientset, autoRunner)

	ctx := context.Background()

	for {
		err := reconciler.Reconcile(ctx)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(5 * time.Second)
	}
}
