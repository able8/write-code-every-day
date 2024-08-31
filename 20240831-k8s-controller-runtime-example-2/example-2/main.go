package main

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ApplicationReconciler reconciles a Pod object
type ApplicationReconciler struct{}

// Reconcile reads that state of the cluster for a Pod and makes changes based on the state
func (a *ApplicationReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Implement your reconciliation logic here
	log.Printf("Reconciling Pod: %s", request.NamespacedName)
	return reconcile.Result{}, nil
}

func start() {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		log.Fatalf("unable to add corev1 scheme: %v", err)
	}

	// 1. Initialize Manager
	mgr, err := manager.New(ctrl.GetConfigOrDie(), manager.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatalf("unable to create manager: %v", err)
	}

	// 2. Initialize Controller using NewControllerManagedBy
	err = ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(event event.CreateEvent) bool {
				log.Printf("Pod created: %s/%s", event.Object.GetNamespace(), event.Object.GetName())
				return true // Process create events
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				log.Printf("Pod updated: %s/%s", updateEvent.ObjectNew.GetNamespace(), updateEvent.ObjectNew.GetName())
				return true // Process update events
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				log.Printf("Pod deleted: %s/%s", deleteEvent.Object.GetNamespace(), deleteEvent.Object.GetName())
				return true // Process delete events
			},
		}).
		Complete(&ApplicationReconciler{})
	if err != nil {
		log.Fatalf("unable to create controller: %v", err)
	}

	// 3. Start the Manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatalf("unable to start manager: %v", err)
	}
}

func main() {
	start()
}
