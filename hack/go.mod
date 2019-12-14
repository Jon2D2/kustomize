module sigs.k8s.io/kustomize/hack

go 1.14

require (
	github.com/spf13/cobra v0.0.5
	sigs.k8s.io/kustomize/api v0.3.2
	sigs.k8s.io/kustomize/cmd/config v0.1.1
	sigs.k8s.io/kustomize/kyaml v0.1.1
)

exclude github.com/golangci/golangci-lint v1.24.0
