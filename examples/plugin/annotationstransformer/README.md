# Test CRD Register annotationstransformer


This folder demonstrates how to use annotationstransformer Kustomize Transformer

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
kind: AnnotationsTransformer
metadata:
  name: annotationstransformer
annotations:
  app: myApp
  greeting/morning: a string with blanks
fieldSpecs:
- path: metadata/annotations
  create: true
- path: spec/template/metadata/annotations
  create: true
  version: v1
  kind: ReplicationController
- path: spec/template/metadata/annotations
  create: true
  kind: Deployment
- path: spec/template/metadata/annotations
  create: true
  kind: ReplicaSet
- path: spec/template/metadata/annotations
  create: true
  kind: DaemonSet
- path: spec/template/metadata/annotations
  create: true
  kind: StatefulSet
- path: spec/template/metadata/annotations
  create: true
  group: batch
  kind: Job
- path: spec/jobTemplate/metadata/annotations
  create: true
  group: batch
  kind: CronJob
- path: spec/jobTemplate/spec/template/metadata/annotations
  create: true
  group: batch
  kind: CronJob
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
  annotations:
    app: myApp
    greeting/morning: a string with blanks
  name: thenamespace
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    app: myApp
    greeting/morning: a string with blanks
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
  annotations:
    app: myApp
    greeting/morning: a string with blanks
  name: cm1
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  annotations:
    app: myApp
    greeting/morning: a string with blanks
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

