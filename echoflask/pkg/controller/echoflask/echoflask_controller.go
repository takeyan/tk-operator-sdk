package echoflask

import (
	"context"
	"reflect"

        appsv1 "k8s.io/api/apps/v1"            // add for Depoyment
        "k8s.io/apimachinery/pkg/util/intstr"  // add for TargetPort
	swallowlabv1alpha1 "echoflask/pkg/apis/swallowlab/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_echoflask")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new EchoFlask Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileEchoFlask{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("echoflask-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource EchoFlask
	err = c.Watch(&source.Kind{Type: &swallowlabv1alpha1.EchoFlask{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}



// WATCH BLOCK BEGIN
        // TODO(user): Modify this to be the types you create that are owned by the primary resource
        // Watch for changes to secondary resource Pods and requeue the owner EchoFlask


        err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
                IsController: true,
                OwnerType: &swallowlabv1alpha1.EchoFlask{},
        })
        if err != nil {
                return err
        }

        err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
                IsController: true,
                OwnerType: &swallowlabv1alpha1.EchoFlask{},
        })
        if err != nil {
                return err
        }

// WATCH BLOCK END
	return nil
}

// blank assignment to verify that ReconcileEchoFlask implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileEchoFlask{}

// ReconcileEchoFlask reconciles a EchoFlask object
type ReconcileEchoFlask struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a EchoFlask object and makes changes based on the state read
// and what is in the EchoFlask.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileEchoFlask) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling EchoFlask")

	// Fetch the EchoFlask instance
	instance := &swallowlabv1alpha1.EchoFlask{}
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


// RECONCILE BLOCK BEGIN

        // Define Deployment name and Service name
        dep_name := instance.Name + "-deployment"
        svc_name := instance.Name + "-svc"

        // Check if the deployment already exists, if not create a new one
        depfound := &appsv1.Deployment{}
        err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep_name, Namespace: instance.Namespace}, depfound)
        if err != nil && errors.IsNotFound(err) {
            // Define a new deployment
            reqLogger.Info("Defining a new Deployment for: " + instance.Name)
            dep := r.newDeploymentForCR(instance, dep_name)
            reqLogger.Info("Creating a App Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
            err = r.client.Create(context.TODO(), dep)
            if err != nil {
                reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
                return reconcile.Result{}, err
            }   
            // Deployment created successfully - return and requeue
            return reconcile.Result{Requeue: true}, nil
        } else if err != nil {
            reqLogger.Error(err, "Failed to get Deployment")
            return reconcile.Result{}, err
        }

        // デプロイメントのReplicasをCRのspecのsizeと同じになるように調整する
        size := instance.Spec.Size
        if *depfound.Spec.Replicas != size {
            depfound.Spec.Replicas = &size
            err = r.client.Update(context.TODO(), depfound)
            if err != nil {
                reqLogger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", depfound.Namespace, "Deployment.Name", depfound.Name)
                return reconcile.Result{}, err
            }
        }

        // Check if the service already exists, if not create a new one
        svcfound := &corev1.Service{}
        err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc_name, Namespace: instance.Namespace}, svcfound)
        if err != nil && errors.IsNotFound(err) {
            // Define a new service
            reqLogger.Info("Defining a new Service for: " + instance.Name)
            svc := r.newServiceForCR(instance, svc_name)
            reqLogger.Info("Creating a App Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
            err = r.client.Create(context.TODO(), svc)
            if err != nil {
                reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
                return reconcile.Result{}, err
            }
        } else if err != nil {
            reqLogger.Error(err, "Failed to get Service")
            return reconcile.Result{}, err
        }


    // Update the CR status with the pod names
        // List the pods for this CR's deployment
        podList := &corev1.PodList{}
        listOpts := []client.ListOption{
            client.InNamespace(instance.Namespace),
            client.MatchingLabels(newLabelsForCR(instance.Name)),
        }
        err = r.client.List(context.TODO(), podList, listOpts...)
        if err != nil {
            reqLogger.Error(err, "Failed to list pods.", "CR.Namespace", instance.Namespace, "CR.Name", instance.Name)
            return reconcile.Result{}, err
        }
        podNames := getPodNames(podList.Items)

        // Update status.Nodes if needed
        if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
            instance.Status.Nodes = podNames
            err := r.client.Status().Update(context.TODO(), instance)
            if err != nil {
                reqLogger.Error(err, "Failed to update CR status.")
                return reconcile.Result{}, err
            }
        }


        // Deployment and Service already exist - don't requeue
        reqLogger.Info("Skip reconcile: Deployment and Service already exists", "Deployment.Name", depfound.Name, "Service.Name", svcfound.Name)
        return reconcile.Result{}, nil

        }       



// newDeploymentForCR returns a busybox pod with the same name/namespace as the cr
func (r *ReconcileEchoFlask) newDeploymentForCR(cr *swallowlabv1alpha1.EchoFlask, dep_name string) *appsv1.Deployment {
    labels := newLabelsForCR(cr.Name)
    dep := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: dep_name,
            Namespace: cr.Namespace,
            Labels: labels,
        },
        Spec: appsv1.DeploymentSpec{
            Selector: &metav1.LabelSelector{
                MatchLabels: labels,
            },
          Replicas: &cr.Spec.Size,
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{Labels: labels },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name: "echoflask",
                           Image: "takeyan/flask:0.0.3",
                            Ports: []corev1.ContainerPort{{
                                ContainerPort: 5000,
                            }},
                            Env: []corev1.EnvVar{
                                {
                                    Name: "K8S_NODE_NAME",
                                    ValueFrom: &corev1.EnvVarSource{ FieldRef: &corev1.ObjectFieldSelector{ FieldPath: "spec.nodeName" }},
                              },
                                {
                                    Name: "K8S_POD_NAME",
                                    ValueFrom: &corev1.EnvVarSource{ FieldRef: &corev1.ObjectFieldSelector{ FieldPath: "metadata.name" }},
                                },
                                {
                                    Name: "K8S_POD_IP",
                                    ValueFrom: &corev1.EnvVarSource{ FieldRef: &corev1.ObjectFieldSelector{ FieldPath: "status.podIP" }},
                                },
                            },
                        },
                   },      
                },
            },
        },
    }
    controllerutil.SetControllerReference(cr, dep, r.scheme)
    return dep
}


func (r *ReconcileEchoFlask) newServiceForCR(cr *swallowlabv1alpha1.EchoFlask, svc_name string) *corev1.Service {
    labels := newLabelsForCR(cr.Name)
    svc := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name: svc_name,
            Namespace: cr.Namespace,
        },
        Spec: corev1.ServiceSpec{
            Ports: []corev1.ServicePort{{
                Protocol: "TCP",
                Port: 5000,
                TargetPort: intstr.FromInt(5000),
            }},
            Type: corev1.ServiceTypeNodePort,
            Selector: labels,
        },
    }
    controllerutil.SetControllerReference(cr, svc, r.scheme)
    return svc    
}

// newLabelsForCR returns the labels for selecting the resources
// belonging to the given CR name.
func newLabelsForCR(name string) map[string]string {
	return map[string]string{"app": "echoflask", "custom_resource": name }
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
    var podNames []string
    for _, pod := range pods {
        podNames = append(podNames, pod.Name)
    }
    return podNames
}

// RECONCILE BLOCK END


