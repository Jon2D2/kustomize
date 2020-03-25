// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package builtinpluginconsts

const (
	namespaceFieldSpecs = `
namespace:
# replace or add namespace field
# on all entities by default
- path: metadata/namespace
  create: true

# Update namespace if necessary
# in the subjects fields
- path: subjects
  kind: RoleBinding
- path: subjects
  kind: ClusterRoleBinding

# Would make sense to skip those
# by default but would break backward
# compatility. 
# - path: metadata/name
#   kind: Namespace

# skip those ClusterWide entities
- path: metadata/namespace
  kind: APIService
  skip: true
- path: metadata/namespace
  kind: CSIDriver
  skip: true
- path: metadata/namespace
  kind: CSINode
  skip: true
- path: metadata/namespace
  kind: CertificateSigningRequest
  skip: true
- path: metadata/namespace
  kind: ClusterRole
  skip: true
- path: metadata/namespace
  kind: ClusterRoleBinding
  skip: true
- path: metadata/namespace
  kind: ComponentStatus
  skip: true
- path: metadata/namespace
  kind: CustomResourceDefinition
  skip: true
- path: metadata/namespace
  kind: MutatingWebhookConfiguration
  skip: true
- path: metadata/namespace
  kind: Namespace
  skip: true
- path: metadata/namespace
  kind: Node
  skip: true
- path: metadata/namespace
  kind: PersistentVolume
  skip: true
- path: metadata/namespace
  kind: PodSecurityPolicy
  skip: true
- path: metadata/namespace
  kind: PriorityClass
  skip: true
- path: metadata/namespace
  kind: RuntimeClass
  skip: true
- path: metadata/namespace
  kind: SelfSubjectAccessReview
  skip: true
- path: metadata/namespace
  kind: SelfSubjectRulesReview
  skip: true
- path: metadata/namespace
  kind: StorageClass
  skip: true
- path: metadata/namespace
  kind: SubjectAccessReview
  skip: true
- path: metadata/namespace
  kind: TokenReview
  skip: true
- path: metadata/namespace
  kind: ValidatingWebhookConfiguration
  skip: true
- path: metadata/namespace
  kind: VolumeAttachment
  skip: true
`
)
