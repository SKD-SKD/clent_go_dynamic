package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
    netv1 "k8s.io/api/networking/v1beta1"
	//v1beta1 "k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/restmapper"
	//"k8s.io/client-go/tools/record"


)

type YAML = string

var api corev1.CoreV1Interface
var clientSet *kubernetes.Clientset

type KubernetesClient struct {
	clientset        kubernetes.Clientset
	dynamicinterface dynamic.Interface
	discoveryclient  *discovery.DiscoveryClient
}

func NewKubernetesClient(configBytes []byte) (*KubernetesClient, error) {
	config, err := clientcmd.NewClientConfigFromBytes(configBytes)
	if err != nil {
		return nil, err
	}
	clientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	dynamicinterface, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	discoveryclient, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &KubernetesClient{clientset: *clientset, dynamicinterface: dynamicinterface, discoveryclient: discoveryclient}, nil
}

func (c *KubernetesClient) ListNamespaces(ctx context.Context) (*v1.NamespaceList, error) {
	return c.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (c *KubernetesClient) CreateNamespace(ctx context.Context, name string) (*v1.Namespace, error) {
	namespace, err := c.clientset.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	_, err = c.clientset.CoreV1().ResourceQuotas(name).Create(context.TODO(), &v1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ResourceQuota"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: name,
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: v1.ResourceList{
				v1.ResourceLimitsCPU:      resource.MustParse("128000m"),
				v1.ResourceLimitsMemory:   resource.MustParse("200Gi"),
				v1.ResourceRequestsCPU:    resource.MustParse("128000m"),
				v1.ResourceRequestsMemory: resource.MustParse("150Gi"),
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (c *KubernetesClient) DeleteNamespace(ctx context.Context, name string) error {
	var gracePeriodSeconds = int64(0)
	return c.clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds})
}

func (c *KubernetesClient) ListServices(ctx context.Context, namespace string) (*v1.ServiceList, error) {
	return c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
}

func (c *KubernetesClient) ListPods(ctx context.Context, namespace string) (*v1.PodList, error) {
	return c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

func (c *KubernetesClient) ListDeploy(ctx context.Context, namespace string) (*appsv1.DeploymentList, error) {
	return c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
}

func (c *KubernetesClient) ListConfigmaps(ctx context.Context, namespace string) (*v1.ConfigMapList, error) {
	return c.clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
}

func (c *KubernetesClient) ListReplicaset(ctx context.Context, namespace string) (*appsv1.ReplicaSetList, error) {
	return c.clientset.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
}

func (c *KubernetesClient) ListEndpoint(ctx context.Context, namespace string) (*v1.EndpointsList, error) {
	return c.clientset.CoreV1().Endpoints(namespace).List(ctx, metav1.ListOptions{})
}

//k8s.io/api/networking/v1
func (c *KubernetesClient) ListIngresses(ctx context.Context, namespace string) (*netv1.IngressList, error) {
	return c.clientset.NetworkingV1beta1().Ingresses(namespace).List(ctx, metav1.ListOptions{} )
}

func (c *KubernetesClient) CreateDeploy(ctx context.Context, namespace string, deployname string, replicas uint32, appname string, containername string, imagetag string) (*unstructured.Unstructured, error) { //(*appsv1.Deployment, error) {

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": deployname,
			},
			"spec": map[string]interface{}{
				"replicas": replicas,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": appname, //for now
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": appname, //for now
						},
					},
					//portest []map[string]interface :=  { "udp", "UDP", 80 }
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":            containername,
								"image":           imagetag,
								"imagePullPolicy": "Always",
								"ports": []map[string]interface{}{
									//{
									//  portest
									//},
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 80,
									},
									{"name": "tcp",
										"protocol":      "TCP",
										"containerPort": 8080,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := c.dynamicinterface.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())

	return result, err
}

func (c *KubernetesClient) DeleteDeploy(ctx context.Context, namespace string, deployname string) error { //(*appsv1.Deployment, error) {
	fmt.Println("Deleting deployment...")

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	var err error
	if err = c.dynamicinterface.Resource(deploymentRes).Namespace(namespace).Delete(context.TODO(), deployname, deleteOptions); err != nil {
		panic(err)
	}

	return err
}

func (c *KubernetesClient) CreateService(ctx context.Context, servicetype v1.ServiceType, namespace string, servicename string, appname string) error { //(*appsv1.Deployment, error) {
	//selector,  template, container ports - port target port

	coreV1Client := c.clientset.CoreV1()

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: servicename,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": appname,
			},
			Type: servicetype, //"ClusterIP" "NodePort" "LoadBalancer" "ExternalName"
			Ports: []v1.ServicePort{
				{
					Protocol: "TCP",
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
				//{
				//	Protocol: "UDP",
				//	Port: 81,
				//	TargetPort: intstr.IntOrString{
				//		Type:   intstr.Int,
				//		IntVal: 80,
				//	},
				//},
			},
		},
	}
	result, err := coreV1Client.Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %s\n", result.ObjectMeta.Name)
	return err

}

func (c *KubernetesClient) DeleteService(ctx context.Context, namespace string, servicename string) error { //(*appsv1.Deployment, error) {
	fmt.Println("Deleting service...")

	coreV1Client := c.clientset.CoreV1()
	var err error
	err = coreV1Client.Services(namespace).Delete(context.TODO(), servicename, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted service \n")

	return err
}

func (c *KubernetesClient) CreateEndpoint(ctx context.Context, namespace string, endpointname string, ip string, portname string, port int32, protocol v1.Protocol) error { //(*appsv1.Deployment, error) {

	coreV1Client := c.clientset.CoreV1()

	endpoints := &v1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name: endpointname,
		},
		Subsets: []v1.EndpointSubset{
			{Addresses: []v1.EndpointAddress{{IP: ip}},
				Ports: []v1.EndpointPort{{Name: portname, Port: port, Protocol: protocol}},
			},
			//{			Addresses: []v1.EndpointAddress{{IP: "4.4.4.4"}},
			//	        Ports:     []v1.EndpointPort{{Name: "4444", Port: 400 , Protocol: "TCP"}},
			//},
		},
	}
	result, err := coreV1Client.Endpoints(namespace).Create(context.TODO(), endpoints, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Endponts %s\n", result.ObjectMeta.Name)
	return err

}

func (c *KubernetesClient) DeleteEndpoint(ctx context.Context, namespace string, ednpontname string) error { //(*appsv1.Deployment, error) {
	fmt.Println("Deleting Endponts...")

	coreV1Client := c.clientset.CoreV1()
	var err error
	err = coreV1Client.Endpoints(namespace).Delete(context.TODO(), ednpontname, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted Endponts \n")

	return nil
}

func (c *KubernetesClient) CreateConfigmap(ctx context.Context, namespace string, configmapname string, configit string, configurewith string) error {

	coreV1Client := c.clientset.CoreV1()

	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configmapname,
		},
		Data: map[string]string{configit: configurewith},
	}
	result, err := coreV1Client.ConfigMaps(namespace).Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created ConfigMap %s\n", result.ObjectMeta.Name)
	return nil

}

func (c *KubernetesClient) DeleteConfigmap(ctx context.Context, namespace string, configmapname string) error {
	fmt.Println("Deleting Configmap...")

	coreV1Client := c.clientset.CoreV1()
	//var err error
	err := coreV1Client.ConfigMaps(namespace).Delete(context.TODO(), configmapname, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted Comfigmap \n")
	return nil
}

func (c *KubernetesClient) CreateApplicationService(ctx context.Context, namespace string, name string, replicas uint32, appname string, imagetag string,
	servicetype v1.ServiceType, configit string, configurewith string,
	db_endpoint string, ip string, portname string, port int32, protocol v1.Protocol) error {
	var err error

	err = c.CreateService(ctx, servicetype, namespace, name, appname)
	err = c.CreateConfigmap(ctx, namespace, name, configit, configurewith) //one config
	_, err = c.CreateDeploy(ctx, namespace, name, replicas, appname, appname, imagetag)
	err = c.CreateEndpoint(ctx, namespace, db_endpoint, ip, portname, port, protocol) //  one endpoint name for now

	fmt.Printf("creaed CG service \n")
	return err
}

func (c *KubernetesClient) DeleteApplicationService(ctx context.Context, namespace string, name string, db_endpoint string) error {
	var err error
	err = c.DeleteService(ctx, namespace, name)
	err = c.DeleteConfigmap(ctx, namespace, name) //one config
	err = c.DeleteDeploy(ctx, namespace, name)
	err = c.DeleteEndpoint(ctx, namespace, db_endpoint) //  one endpoint name for now
	fmt.Printf("deleted CG service \n")
	return err
}

var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

func (c *KubernetesClient) CreateDynamicUnstructured(ctx context.Context, yaml string) error {

	dc := c.discoveryclient
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	dyn := c.dynamicinterface

	//	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(yaml), nil, obj)
	if err != nil {
		panic(err)
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		panic(err)
	}
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "cg-controller",
	})
	if err != nil {
		panic(err)
	}

	//time.Sleep (500*1000*time.Millisecond )
	return err
}

type runStruct struct {
	Kind        string
}



func ProcessEventDeploy(watchmy watch.Interface) {
	for {
		event := <-watchmy.ResultChan()
		fmt.Printf("\n\nafter event := <-ch  %v \n", event )
		fmt.Printf("after event := <-ch  ===Type %v  \n", event.Type )
		////fmt.Printf("after event := <-ch  ===Object %v  \n", event.Object.GetObjectKind() )
        // some outputs to discover
		//fmt.Printf("after event := <-ch  ===Object %v  \n", event.Object.GetObjectKind())
		//fmt.Printf("after event := <-ch  ===: %v  \n", event.Object.DeepCopyObject().GetObjectKind())
		//fmt.Printf("after event := <-ch  ===> %v  \n", event.Object.DeepCopyObject())
		//myType :=  event.Type
		//myObject := event.Object
		//mysomethingC := myObject.DeepCopyObject()
		//mysomething := mysomethingC
		//fmt.Printf("something ---------: %v \n", mysomething)

		p := event.Object.GetObjectKind()
		// convert map to json
		jsonString, _ := json.Marshal(p)

		check := runStruct{}
		json.Unmarshal(jsonString, &check)
		fmt.Printf( "event Kind -> : %s \n",  check.Kind)
		fmt.Printf("@@ event is Kind of  -%s- \n", check.Kind)

		//going for .....
		if ( check.Kind == "Deployment" ) {
			fmt.Printf("is Kind  -%s- Deployment ? \n", check.Kind)
			s := appsv1.Deployment{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace, " spec replicas:", *(s.Spec.Replicas), " status Replicas:", s.Status.Replicas, " status Replicas:",  s.Status.ReadyReplicas)
		}
		if ( check.Kind == "StatefulSet" ) {
			fmt.Printf("is Kind  -%s- StatefulSet ? \n", check.Kind)
			s := appsv1.StatefulSet{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace, " spec replicas:", *(s.Spec.Replicas), " status Replicas:", s.Status.Replicas, " status Replicas:",  s.Status.ReadyReplicas)
		}
		if ( check.Kind == "DaemonSet" ) {
			fmt.Printf("is Kind  -%s- DaemonSet ? \n", check.Kind)
			s := appsv1.StatefulSet{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace, " spec replicas:", *(s.Spec.Replicas), " status Replicas:", s.Status.Replicas, " status Replicas:",  s.Status.ReadyReplicas)
		}
		if ( check.Kind == "Service" ) {
			fmt.Printf("is Kind  -%s- Service ? \n", check.Kind)
			s := v1.Service{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace )
		}
		if ( check.Kind == "ConfigMap" ) {
			fmt.Printf("is Kind  -%s- ConfigMap ? \n", check.Kind)
			s := v1.ConfigMap{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace )
		}
		if ( check.Kind == "Secret" ) {
			fmt.Printf("is Kind  -%s- Secret ? \n", check.Kind)
			s := v1.Secret{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace )
		}
		if ( check.Kind == "ServiceAccount" ) {
			fmt.Printf("is Kind  -%s- Secret ? \n", check.Kind)
			s := v1.ServiceAccount{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace )
		}
		if ( check.Kind == "Ingress" ) {
			fmt.Printf("is Kind  -%s- Ingress ? \n", check.Kind)
			s := netv1.Ingress{}
			json.Unmarshal(jsonString, &s)
			fmt.Println("Name:",s.Name, " namespace: ", s.Namespace, " status: ", s.Status , " spec: ", s.Spec/*.Rules*/ )
		}

	}
}


func (c *KubernetesClient) DeleteDynamicUnstructured(ctx context.Context, yaml string) error {

	dc := c.discoveryclient
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	dyn := c.dynamicinterface

	//	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(yaml), nil, obj)
	if err != nil {
		panic(err)
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		panic(err)
	}
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	obj.GetName()
	err = dr.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	return err
}


func (c *KubernetesClient) EventsDynamicUnstructured(ctx context.Context, yaml string/*, run func(watch.Interface)*/) (watch.Interface, error) {

	dc := c.discoveryclient
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	dyn := c.dynamicinterface

	//	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(yaml), nil, obj)
	if err != nil {
		panic(err)
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		panic(err)
	}
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	var watchmy watch.Interface
	watchmy, err = dr.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("watchmy  Watch fine \n")

	//var Recorder record
	//if err := Recorder.Record(obj); err != nil {
	//	klog.V(4).Infof("error recording current command: %v", err)
	//}

	return watchmy, err
}



func EventsInstallCallDynamicUnstructured( run func(watch.Interface), watchmy watch.Interface ) {

	go run (watchmy)
	///time.Sleep (5000*1000*time.Millisecond ) // let it go at some time for now

	return
}
