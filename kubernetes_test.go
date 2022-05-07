// +build manual_integration

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"sigs.k8s.io/yaml"
	"testing"
	"time"
)

func TestListNamespaces(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	namespaces, err := client.ListNamespaces(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, namespaces)
}

func TestListServices(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	services, err := client.ListServices(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, services)
}

func TestListPods(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	pods, err := client.ListPods(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, pods)
}

func TestListDeploy(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	deploy, err := client.ListDeploy(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, deploy)
}

func TestListConfigMaps(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	configmap, err := client.ListConfigmaps(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, configmap)
}

func TestListReplicaset(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	configmap, err := client.ListReplicaset(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, configmap)
}

func TestListEndPoint(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	endpoint, err := client.ListEndpoint(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, endpoint)
}

func TestListIngresses(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	endpoint, err := client.ListIngresses(context.Background(), "default")

	assert.Nil(t, err)
	assert.NotNil(t, endpoint)
}

func TestCreateDeploy(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	createdeply, err := client.CreateDeploy(context.Background(), "default", "deploy-name", 1, "app-name", "container-name", "nginx:1.12")

	assert.Nil(t, err)
	assert.NotNil(t, createdeply)

}

func TestDeleteDeploy(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.DeleteDeploy(context.Background(), "default", "deploy-name")

	assert.Nil(t, err)
}

func TestCreateService(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	err = client.CreateService(context.Background(), "ClusterIP", "default", "service-name", "app-name")

	assert.Nil(t, err)
}

func TestDeleteService(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	err = client.DeleteService(context.Background(), "default", "service-name")

	assert.Nil(t, err)
}

func TestCreateEndpoint(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	err = client.CreateEndpoint(context.Background(), "default", "ep-name", "5.5.5.5", "myport", 100, "TCP")
	assert.Nil(t, err)
}

func TestDeleteEndpoint(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.DeleteEndpoint(context.Background(), "default", "ep-name")

	assert.Nil(t, err)
}

func TestCreateComfigmap(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.CreateConfigmap(context.Background(), "default", "cm-name", "config-it", "config=with")

	assert.Nil(t, err)
}

func TestDeleteConfigmap(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.DeleteConfigmap(context.Background(), "default", "cm-name")
	assert.Nil(t, err)
}

func TestCreateApplicationService(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	err = client.CreateApplicationService(context.Background(), "default", "my-name", 2, "my-appname", "nginx:1.12",
		"ClusterIP", "my-config-item", "value",
		"my-db", "5.5.5.5", "portname", 1002, "TCP")

	assert.Nil(t, err)
}

func TestDeleteApplicationService(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.DeleteApplicationService(context.Background(), "default", "my-name", "my-db")

	assert.Nil(t, err)
}


//TestdeploymentYAMLConfigMapSiteCars
func TestSiteCarsCreateDynamicUnstructured(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLConfigMapSiteCars)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLConfigMap2SiteCars)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLServiceSiteCars)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLDeploymentSiteCars)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLIngressSiteCars)
	assert.Nil(t, err)
}

func TestRedisCreateDynamicUnstructured(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis1)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis2)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis3)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis4)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis5)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis6)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis7)
	assert.Nil(t, err)
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis8)
	assert.Nil(t, err)
	//stateful set
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis9)
	assert.Nil(t, err)
	//stateful set
	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAMLReadis11)
	assert.Nil(t, err)

	fmt.Printf("RedisDynamic created  \n")
}

const TestdeploymentYAMLDeploy = `
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: skornfeld
`
const TestdeploymentYAMLIngress= `
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  namespace: skornfeld
`
const TestdeploymentYAMLServiceAccount= `
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: skornfeld
`
const TestdeploymentYAMLSecret = `
apiVersion: v1
kind: Secret
metadata:
  namespace: skornfeld
`
const TestdeploymentYAMLConfigMap = `
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: skornfeld
`
const TestdeploymentYAMLService = `
apiVersion: v1
kind: Service
metadata:
  namespace: skornfeld
`
const TestdeploymentYAMLStatefulSet = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  namespace: "skornfeld"
`

func TestRedisEventsDynamicUnstructured(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)
	watchme, err := client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLServiceAccount)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	watchme, err = client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLSecret)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	watchme, err = client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLConfigMap)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	watchme, err = client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLService)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	watchme, err = client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLStatefulSet)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	watchme, err = client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLDeploy)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	watchme, err = client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAMLIngress)
	//if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)
	assert.Nil(t, err)

	fmt.Printf("EventsDynamic installed  \n")

	time.Sleep (10000*1000*time.Millisecond ) // let it go at some time for now, forever later
}


func TestCreateDynamicUnstructured(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.CreateDynamicUnstructured(context.Background(), TestdeploymentYAML)
	assert.Nil(t, err)
}

func TestDeleleDynamicUnstructured(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	err = client.DeleteDynamicUnstructured(context.Background(), TestdeploymentYAML)
	assert.Nil(t, err)
}

func TestEventsDynamicUnstructured(t *testing.T) {
	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".kube", "config"))
	var config KubeConfig
	err = yaml.Unmarshal(bytes, &config)
	bytes, err = json.Marshal(config)

	client, err := NewKubernetesClient(bytes)

	watchme, err := client.EventsDynamicUnstructured(context.Background(), TestdeploymentYAML)
    //if not err ..
	EventsInstallCallDynamicUnstructured(ProcessEventDeploy, watchme)

	assert.Nil(t, err)
}

