# Test CRD Register namespacetransformer


This folder demonstrates how to use namespacetransformer Kustomize Transformer

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/
```

### Preparation Step Kustomize File

<!-- @createkustomization.yaml @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- resources.yaml
transformers:
- transformer.yaml
EOF
```


### Preparation Step Resources

<!-- @createresources.yaml @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resources.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: thenamespace
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns:
  group: my.org
  version: v1alpha1
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          type: object
          properties:
            simpletext:
              type: string
            replicas:
              type: integer
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  simpletext: some simple text
  replicas: 123
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
EOF
```


### Preparation Step Transformer Config

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/transformer.yaml
apiVersion: builtin
kind: NamespaceTransformer
metadata:
  name: namespacetransformer
  namespace: tutorial-ns
fieldSpecs:
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

# Update the namespace object itself
- path: metadata/name
  kind: Namespace

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

# This a Cluster wide CRD
- path: metadata/namespace
  kind: MyCRD
  skip: true

EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/ -o ${DEMO_HOME}/actual/result.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected Result

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/result.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: tutorial-ns
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            replicas:
              type: integer
            simpletext:
              type: string
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
  namespace: tutorial-ns
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replicas: 123
  simpletext: some simple text
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

