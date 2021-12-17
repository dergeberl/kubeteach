module github.com/dergeberl/kubeteach

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.6
	github.com/go-logr/logr v0.4.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/prometheus/client_golang v1.11.0
	github.com/tidwall/gjson v1.10.2
	go.uber.org/automaxprocs v1.4.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	sigs.k8s.io/controller-runtime v0.10.0
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
