module github.com/dergeberl/kubeteach

go 1.16

require (
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/tidwall/gjson v1.8.1
	k8s.io/api v0.20.10
	k8s.io/apimachinery v0.20.10
	k8s.io/client-go v0.20.10
	sigs.k8s.io/controller-runtime v0.8.3
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
