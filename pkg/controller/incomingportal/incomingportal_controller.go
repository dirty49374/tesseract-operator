package incomingportal

import (
	"context"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
	"github.com/dirty49374/tesseract-operator/pkg/certs"
	"github.com/dirty49374/tesseract-operator/pkg/util"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type Empty struct{}

var log = logf.Log.WithName("controller_incomingportal")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new IncomingPortal Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	certs, err := certs.LoadCerts("secret")
	if err != nil {
		return err
	}

	return add(mgr, newReconciler(mgr, certs))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, certs *certs.Certs) reconcile.Reconciler {
	return &ReconcileIncomingPortal{client: mgr.GetClient(), scheme: mgr.GetScheme(), certs: certs, portals: make(map[string][]int32)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("incomingportal-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IncomingPortal
	err = c.Watch(&source.Kind{Type: &tesseractv1alpha1.IncomingPortal{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner IncomingPortal
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tesseractv1alpha1.IncomingPortal{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileIncomingPortal{}

// ReconcileIncomingPortal reconciles a IncomingPortal object
type ReconcileIncomingPortal struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	certs   *certs.Certs
	portals map[string][]int32
}

// Reconcile reads that state of the cluster for a IncomingPortal object and makes changes based on the state read
// and what is in the IncomingPortal.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIncomingPortal) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	addr := request.Name + "." + request.Namespace

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling IncomingPortal")

	// Fetch the IncomingPortal instance
	instance := &tesseractv1alpha1.IncomingPortal{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		// Request object not found, could have been deleted after reconcile request.
		delete(r.portals, addr)
	} else {
		r.portals[addr] = instance.Spec.LocalPorts
		if r.portals[addr] == nil {
			r.portals[addr] = make([]int32, 0)
		}
	}

	configMap, err := r.newConfigMapForCR()
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.client.Update(context.TODO(), configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), configMap)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	}
	hash, err := util.JsonShaObject(configMap)
	if err != nil {
		return reconcile.Result{}, err
	}

	deployment := newDeployment(hash)
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), deployment)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
