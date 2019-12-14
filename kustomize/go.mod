module sigs.k8s.io/kustomize/kustomize/v3

go 1.13

require (
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	k8s.io/client-go v0.17.0
	sigs.k8s.io/kustomize/api v0.3.2
	sigs.k8s.io/yaml v1.1.0
)

replace sigs.k8s.io/kustomize/api => ../api
