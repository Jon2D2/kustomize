# Feature Test for Issue 0552


This folder contains files describing how to address [Issue 0552](https://github.com/kubernetes-sigs/kustomize/issues/0552)

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
mkdir -p ${DEMO_HOME}/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
namespace: thenamespace
resources:
- resources.yaml
configurations:
- ./kustomizeconfig/mycrd.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/mycrd.yaml
namespace:
- path: metadata/namespace
  version: v1alpha1
  group: my.org
  kind: MyCRD
  skip: true
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resources.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: thenamespace
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec: |-
  some simple text
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
          type: string
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build $DEMO_HOME -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
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
          type: string
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec: some simple text
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_namespace_thenamespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: thenamespace
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/thenamespace_~g_v1_configmap_cm1.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
  namespace: thenamespace
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

