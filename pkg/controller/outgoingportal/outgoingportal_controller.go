package outgoingportal

import (
	"context"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"

	"github.com/dirty49374/tesseract-operator/pkg/certs"
	"github.com/dirty49374/tesseract-operator/pkg/util"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_outgoingportal")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new OutgoingPortal Controller and adds it to the Manager. The Manager will set fields on the Controller
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
	return &ReconcileOutgoingPortal{client: mgr.GetClient(), scheme: mgr.GetScheme(), certs: certs}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("outgoingportal-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource OutgoingPortal
	err = c.Watch(&source.Kind{Type: &tesseractv1alpha1.OutgoingPortal{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Deployments and requeue the owner OutgoingPortal
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tesseractv1alpha1.OutgoingPortal{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Services and requeue the owner OutgoingPortal
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tesseractv1alpha1.OutgoingPortal{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileOutgoingPortal{}

// ReconcileOutgoingPortal reconciles a OutgoingPortal object
type ReconcileOutgoingPortal struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	certs  *certs.Certs
}

// Reconcile reads that state of the cluster for a OutgoingPortal object and makes changes based on the state read
// and what is in the OutgoingPortal.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOutgoingPortal) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling OutgoingPortal")

	// Fetch the OutgoingPortal instance
	instance := &tesseractv1alpha1.OutgoingPortal{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	hash, err := r.reconcileConfigMap(instance, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileSecret(instance, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileDeployment(instance, hash, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileService(instance, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileOutgoingPortal) reconcileConfigMap(instance *tesseractv1alpha1.OutgoingPortal, reqLogger logr.Logger) (string, error) {
	configmap, err := r.newConfigMapForCR(instance)
	if err != nil {
		return "", err
	}
	sha, err := util.JsonShaObject(configmap)
	if err != nil {
		return "", err
	}

	// Set OutgoingPortal instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, configmap, r.scheme); err != nil {
		return "", err
	}

	// Check if this ConfigMap already exists
	found := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, found)
	if err != nil {

		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configmap.Namespace, "ConfigMap.Name", configmap.Name)
			err = r.client.Create(context.TODO(), configmap)
			if err != nil {
				return "", err
			}
			// ConfigMap created successfully - don't requeue
			return sha, nil
		}

		return "", err
	}

	err = r.client.Update(context.TODO(), configmap)
	if err != nil {
		return "", err
	}

	reqLogger.Info("reconcile: updating existing ConfigMap", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	return sha, nil
}

func (r *ReconcileOutgoingPortal) reconcileSecret(instance *tesseractv1alpha1.OutgoingPortal, reqLogger logr.Logger) error {
	secret := newSecretForCR(instance, r.certs)

	// Set OutgoingPortal instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		return err
	}

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			err = r.client.Create(context.TODO(), secret)
			if err != nil {
				return err
			}
			// Secret created successfully - don't requeue
			return nil

		} else {
			return err
		}
	}

	// err = r.client.Update(context.TODO(), secret)
	// if err != nil {
	// 	return err
	// }

	reqLogger.Info("reconcile: updating existing Secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
	return nil
}

func (r *ReconcileOutgoingPortal) reconcileDeployment(instance *tesseractv1alpha1.OutgoingPortal, hash string, reqLogger logr.Logger) error {
	deployment := newDeploymentForCR(instance, hash)

	// Set OutgoingPortal instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return err
	}

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			err = r.client.Create(context.TODO(), deployment)
			if err != nil {
				return err
			}
			// Deployment created successfully - don't requeue
			return nil

		} else {
			return err
		}
	}

	if found.Annotations["config"] != hash {
		err = r.client.Update(context.TODO(), deployment)
		if err != nil {
			return err
		}
		reqLogger.Info("reconcile: updating existing Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	} else {
		reqLogger.Info("reconcile: Skipping existing Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	}

	return nil
}

func (r *ReconcileOutgoingPortal) reconcileService(instance *tesseractv1alpha1.OutgoingPortal, reqLogger logr.Logger) error {
	service := newServiceForCR(instance)

	// Set OutgoingPortal instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return err
	}

	// Check if this Service already exists
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			err = r.client.Create(context.TODO(), service)
			if err != nil {
				return err
			}
			// Service created successfully - don't requeue
			return nil

		} else {
			return err
		}
	}

	// err = r.client.Update(context.TODO(), service)
	// if err != nil {
	// 	return err
	// }

	reqLogger.Info("reconcile: Skip existing Service", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	return nil
}
