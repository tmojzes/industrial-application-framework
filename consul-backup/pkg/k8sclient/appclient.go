package k8sclient

import (
	log "github.com/sirupsen/logrus"

	appdacnokiacomv1alpha1 "github.com/nokia/industrial-application-framework/consul-backup/api/v1alpha1"
	//	appv1alpha1 "github.com/nokia/industrial-application-framework/consul-backup/api/v1alpha1"
	//	appv1alpha1 "gitlabe2.ext.net.nokia.com/Nokia_DAaaS/edge-microservices/application-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	cfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Config    *rest.Config
	k8sClient client.Client
	scheme    = runtime.NewScheme()
)

func init() {
	utilruntime.Must(corev1.AddToScheme(scheme))      //Scheme for Core V1
	utilruntime.Must(appdacnokiacomv1alpha1.AddToScheme(scheme)) //Scheme for Applications
}

func GetK8sClient() (client.Client, error) {
	if k8sClient != nil {
		return k8sClient, nil
	} else {
		return CreateK8sClient()
	}
}

func CreateK8sClient() (client.Client, error) {
	log.Info("CreateK8sClient called")

	if k8sClient != nil {
		log.Info("k8sClient already created")
		return k8sClient, nil
	}

	log.Info("no k8sClient")

	if Config == nil {
		log.Info("no config")
		Config = cfg.GetConfigOrDie()
	}

	log.Info("have config")

	cl, err := client.New(Config, client.Options{Scheme: scheme})
	if err != nil {
		log.Error(err, "failed to create client")
		return nil, err
	}
	log.Info("client created")
	return cl, nil
}
