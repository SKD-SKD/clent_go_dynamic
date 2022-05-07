package kubernetes

type Kind = string

const (
	ConfigMapKind          = "ConfigMap"
	ClusterRoleKind        = "ClusterRole"
	ClusterRoleBindingKind = "ClusterRoleBinding"
	DeploymentKind         = "Deployment"
	IngressKind            = "Ingress"
	JobKind                = "Job"
	NamespaceKind          = "Namespace"
	PodKind                = "Pod"
	RoleKind               = "Role"        // Not currently used
	RoleBindingKind        = "RoleBinding" // Not currently used
	SecretKind             = "Secret"
	ServiceKind            = "Service"
	ServiceMonitorKind     = "ServiceMonitor"
	StatefulSetKind        = "StatefulSet"
)

var Kinds = map[Kind]struct{}{
	ConfigMapKind:          {},
	ClusterRoleKind:        {},
	ClusterRoleBindingKind: {},
	DeploymentKind:         {},
	IngressKind:            {},
	JobKind:                {},
	NamespaceKind:          {},
	PodKind:                {},
	RoleKind:               {},
	RoleBindingKind:        {},
	SecretKind:             {},
	ServiceKind:            {},
	ServiceMonitorKind:     {},
	StatefulSetKind:        {},
}
