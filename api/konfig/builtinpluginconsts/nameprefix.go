// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package builtinpluginconsts

const (
	namePrefixFieldSpecs = `
namePrefix:
- path: metadata/name
- path: metadata/name
  kind: CustomResourceDefinition
  skip: true

# Following merge PR broke backward compatility
# https://github.com/kubernetes-sigs/kustomize/pull/1526
- path: metadata/name
  kind: APIService
  group: apiregistration.k8s.io
  skip: true

# Would make sense to skip those
# by default but would break backward
# compatility
#
# - path: metadata/name
#   kind: Namespace
#   skip: true
# - path: metadata/name
#   group: storage.k8s.io
#   kind: StorageClass
#   skip: true
`
)
