# Feature Test for Issue 1368


This folder contains files describing how to address [Issue 1368](https://github.com/kubernetes-sigs/kustomize/issues/1368)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: kubeapp-ns

commonLabels:
  app: kubeapp

resources:
- ./namespace.yaml
- ./mycrd.yaml
- ./ingress.yaml
- ./service.yaml
- ./deployment.yaml
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
apiVersion: apps/v1
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
cat <<'EOF' >${DEMO_HOME}/base/kubectlapplyordertransformer.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: KubectlApplyOrderTransformer
metadata:
  name: kubectlapplyordertransformer
kindorder:
- CustomResourceDefinition
- Namespace
- Deployment
- Service
- Ingress
- MyCRD
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kubectldeleteordertransformer.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: KubectlDeleteOrderTransformer
metadata:
  name: kubectldeleteordertransformer
kindorder:
- MyCRD
- Ingress
- Service
- Deployment
- Namespace
- CustomResourceDefinition
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
  scope: Namespaced
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

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/base -o ${DEMO_HOME}/actual/kubectl-apply-order.yaml --reorder=kubectlapply --enable_alpha_plugins
kustomize build ${DEMO_HOME}/base -o ${DEMO_HOME}/actual/kubectl-delete-order.yaml --reorder=kubectldelete --enable_alpha_plugins
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kubectl-apply-order.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  labels:
    app: kubeapp
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Namespaced
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
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: kubeapp
  name: kubeapp-ns
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
  template:
    metadata:
      labels:
        app: kubeapp
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:latest
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
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp
    servicePort: 80
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  labels:
    app: kubeapp
  name: my-crd
  namespace: kubeapp-ns
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kubectl-delete-order.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  labels:
    app: kubeapp
  name: my-crd
  namespace: kubeapp-ns
spec:
  replica: 123
  simpletext: some simple text
---
apiVersion: apps/v1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp
    servicePort: 80
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
  template:
    metadata:
      labels:
        app: kubeapp
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:latest
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
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: kubeapp
  name: kubeapp-ns
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  labels:
    app: kubeapp
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Namespaced
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


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

