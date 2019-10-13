# Feature Test for Issue 0519_b

This folder contains files describing how to address [Issue 0519](https://github.com/kubernetes-sigs/kustomize/issues/0519)
Original kubernetes files have imported from [here](https://github.com/DockbitExamples/kubernetes)

This example is using either:
- issue_0519_b: using complete transformer config replacement
- issue_0519_c: skip option to select which components are changed by common transformers.
- issue_0519_d: multibase/composition to select which the components changed by transformers.

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
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/base/kustomizeconfig
mkdir -p ${DEMO_HOME}/canary
mkdir -p ${DEMO_HOME}/canary/kustomizeconfig
mkdir -p ${DEMO_HOME}/production
mkdir -p ${DEMO_HOME}/production/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ./namespace.yaml
- ./mycrd.yaml
- ./ingress.yaml
- ./service.yaml
- ./deployment.yaml

transformers:
- ./kustomizeconfig/namespacetransformer.yaml
- ./kustomizeconfig/commonlabelstransformer.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/canary/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../base

transformers:
- ./kustomizeconfig/imagetransformer.yaml
- ./kustomizeconfig/namesuffixtransformer.yaml
- ./kustomizeconfig/commonlabelstransformer.yaml
- ./kustomizeconfig/patchstrategicmergetransformer.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../base

transformers:
- ./kustomizeconfig/imagetransformer.yaml
- ./kustomizeconfig/namesuffixtransformer.yaml
- ./kustomizeconfig/commonlabelstransformer.yaml
- ./kustomizeconfig/patchstrategicmergetransformer.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: kubeapp
spec:
  replicas: 1
  template:
    metadata:
      name: kubeapp
      labels:
        app: kubeapp
    spec:
      containers:
      - name: kubeapp
        image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
        ports:
        - name: kubeapp
          containerPort: 8081
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/ingress.yaml
kind: Ingress
apiVersion: networking.k8s.io/v1beta1
metadata:
  name: kubeapp
spec:
  backend:
    serviceName: kubeapp
    servicePort: 80
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/commonlabelstransformer.yaml
apiVersion: builtin
kind: LabelTransformer
metadata:
  name: labeltransformer
labels:
  app: kubeapp
fieldSpecs:
- path: metadata/labels
  create: true
  kind: Service
- path: spec/selector
  create: true
  kind: Service
- path: metadata/labels
  create: true
  kind: Deployment
- path: spec/selector/matchLabels
  create: true
  kind: Deployment
- path: spec/template/metadata/labels
  create: true
  kind: Deployment
- path: metadata/labels
  create: true
  kind: Ingress
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/namespacetransformer.yaml
apiVersion: builtin
kind: NamespaceTransformer
metadata:
  name: namespacetransformer
  namespace: kubeapp-ns
fieldSpecs:
- path: metadata/namespace
  create: true
  kind: Service
- path: metadata/namespace
  create: true
  kind: Deployment
- path: metadata/namespace
  create: true
  kind: Ingress
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/mycrd.yaml
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
            replica:
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
  replica: 123
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: kubeapp
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: 8081
    protocol: TCP
  selector:
    app: kubeapp
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/canary/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
  rules:
  - host: canary.foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-canary
          servicePort: 80
  - host: foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-production
          servicePort: 80
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/canary/kustomizeconfig/commonlabelstransformer.yaml
apiVersion: builtin
kind: LabelTransformer
metadata:
  name: labeltransformer
labels:
  env: canary
fieldSpecs:
- path: metadata/labels
  create: true
  kind: Service
- path: spec/selector
  create: true
  kind: Service
- path: metadata/labels
  create: true
  kind: Deployment
- path: spec/selector/matchLabels
  create: true
  kind: Deployment
- path: spec/template/metadata/labels
  create: true
  kind: Deployment
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/canary/kustomizeconfig/imagetransformer.yaml
apiVersion: builtin
kind: ImageTagTransformer
metadata:
  name: imagetagtransformer
imageTag:
  name: hack4easy/kubesim_health-amd64
  newTag: 0.1.9
# fieldSpecs is left empty since `containers` and `initContainers`
# of *ANY* kind in *ANY* path are builtin supported in code
fieldSpecs:
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/canary/kustomizeconfig/namesuffixtransformer.yaml
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: customPrefixer
suffix: -canary
fieldSpecs:
- kind: Deployment
  path: metadata/name
- kind: Service
  path: metadata/name
- kind: Ingress
  path: spec/backend/serviceName

EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/canary/kustomizeconfig/patchstrategicmergetransformer.yaml
apiVersion: builtin
kind: PatchStrategicMergeTransformer
metadata:
  name: patchstrategicmergetransformer
paths:
- ingress.yaml
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/kustomizeconfig/commonlabelstransformer.yaml
apiVersion: builtin
kind: LabelTransformer
metadata:
  name: labeltransformer
labels:
  env: production
fieldSpecs:
- path: metadata/labels
  create: true
  kind: Service
- path: spec/selector
  create: true
  kind: Service
- path: metadata/labels
  create: true
  kind: Deployment
- path: spec/selector/matchLabels
  create: true
  kind: Deployment
- path: spec/template/metadata/labels
  create: true
  kind: Deployment
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/kustomizeconfig/imagetransformer.yaml
apiVersion: builtin
kind: ImageTagTransformer
metadata:
  name: imagetagtransformer
imageTag:
  name: hack4easy/kubesim_health-amd64
  newTag: 0.1.0
# fieldSpecs is left empty since `containers` and `initContainers`
# of *ANY* kind in *ANY* path are builtin supported in code
fieldSpecs:
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/kustomizeconfig/namesuffixtransformer.yaml
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: customPrefixer
suffix: -production
fieldSpecs:
- kind: Deployment
  path: metadata/name
- kind: Service
  path: metadata/name
- kind: Ingress
  path: spec/backend/serviceName

EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/kustomizeconfig/patchstrategicmergetransformer.yaml
apiVersion: builtin
kind: PatchStrategicMergeTransformer
metadata:
  name: patchstrategicmergetransformer
paths:
- ingress.yaml
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
mkdir ${DEMO_HOME}/actual/production
mkdir ${DEMO_HOME}/actual/canary
kustomize build ${DEMO_HOME}/production -o ${DEMO_HOME}/actual/production --enable_alpha_plugins
kustomize build ${DEMO_HOME}/canary -o ${DEMO_HOME}/actual/canary --enable_alpha_plugins
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
mkdir ${DEMO_HOME}/expected/production
mkdir ${DEMO_HOME}/expected/canary
```
### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/canary/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
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
            replica:
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
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/canary/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/canary/~g_v1_namespace_kubeapp-ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/canary/kubeapp-ns_apps_v1_deployment_kubeapp-canary.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Deployment
metadata:
  labels:
    app: kubeapp
    env: canary
  name: kubeapp-canary
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
      env: canary
  template:
    metadata:
      labels:
        app: kubeapp
        env: canary
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:0.1.9
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: kubeapp
        ports:
        - containerPort: 8081
          name: kubeapp
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/canary/kubeapp-ns_networking.k8s.io_v1beta1_ingress_kubeapp.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
  rules:
  - host: canary.foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-canary
          servicePort: 80
  - host: foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-production
          servicePort: 80
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/canary/kubeapp-ns_~g_v1_service_kubeapp-canary.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
    env: canary
  name: kubeapp-canary
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
    env: canary
  type: LoadBalancer
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
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
            replica:
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
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/~g_v1_namespace_kubeapp-ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/kubeapp-ns_apps_v1_deployment_kubeapp-production.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: kubeapp
    env: production
  name: kubeapp-production
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
      env: production
  template:
    metadata:
      labels:
        app: kubeapp
        env: production
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:0.1.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: kubeapp
        ports:
        - containerPort: 8081
          name: kubeapp
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/kubeapp-ns_networking.k8s.io_v1beta1_ingress_kubeapp.yaml
apiVersion: apps/v1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/kubeapp-ns_~g_v1_service_kubeapp-production.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
    env: production
  name: kubeapp-production
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
    env: production
  type: LoadBalancer
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

