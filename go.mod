module github.com/dergeberl/kubeteach

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/prometheus/client_golang v1.7.1
	github.com/tidwall/gjson v1.9.3
	go.uber.org/automaxprocs v1.4.0
	k8s.io/api v0.20.10
	k8s.io/apimachinery v0.20.10
	k8s.io/client-go v0.20.10
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009
	sigs.k8s.io/controller-runtime v0.8.3
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
