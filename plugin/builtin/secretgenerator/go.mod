module sigs.k8s.io/kustomize/plugin/builtin/secretgenerator

go 1.13

require (
	gopkg.in/yaml.v2 v2.2.7 // indirect
	sigs.k8s.io/kustomize/api v0.3.1
	sigs.k8s.io/yaml v1.1.0
)

replace sigs.k8s.io/kustomize/api => ../../../api
