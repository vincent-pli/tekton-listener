module github.com/vincent-pli/tekton-listener

go 1.12

require (
	github.com/cloudevents/sdk-go v0.0.0-20190617162319-ccd952c41493 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/google/go-containerregistry v0.0.0-20190531175139-2687bd5ba651 // indirect
	github.com/knative/build v0.6.0 // indirect
	github.com/knative/eventing-sources v0.6.0
	github.com/knative/pkg v0.0.0-20190604144441-678bb6612d2c // indirect
	github.com/knative/serving v0.6.0
	github.com/kubernetes-sigs/controller-runtime v0.1.12
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/prometheus/common v0.2.0
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/tektoncd/pipeline v0.4.0
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.0-beta.1
	sigs.k8s.io/controller-tools v0.2.0-beta.1 // indirect
)
