package main

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
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

// start initializes the controller manager and starts the reconciliation process
func start() {
	scheme := runtime.NewScheme()
	// Add corev1 scheme to the scheme
	if err := corev1.AddToScheme(scheme); err != nil {
		log.Fatalf("Error adding to scheme: %v", err)
	}

	// 1. Initialize the Manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatalf("Error starting manager: %v", err)
	}

	// 2. Initialize Reconciler (Controller)
	err = ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(&ApplicationReconciler{})
	if err != nil {
		log.Fatalf("Error creating controller: %v", err)
	}

	// 3. Start the Manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatalf("Error starting manager: %v", err)
	}
}

// main function to run the application
func main() {
	start()
}
