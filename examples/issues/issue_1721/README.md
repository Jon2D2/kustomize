# Feature Test for Issue 1721


This folder contains files describing how to address [Issue 1721](https://github.com/kubernetes-sigs/kustomize/issues/1721)

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
mkdir -p ${DEMO_HOME}/overlay
mkdir -p ${DEMO_HOME}/overlay/prod
mkdir -p ${DEMO_HOME}/overlay/staging
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ingress.yaml
- service.yaml
- deployment.yaml
- horizontal_pod_autoscaler.yaml

configurations:
- kustomizeconfig/name_references.yaml
# Leverage auto-var feature instead doing it manually
# Uncommment following section to do it manually
# - kustomizeconfig/var_references.yaml

# Leverage auto-var feature instead doing it manually
# Uncommment following section to do it manually
# vars:
# - name: HorizontalPodAutoscaler.puppetserver.spec.minReplicas
#   objref:
#     kind: HorizontalPodAutoscaler
#     name: puppetserver
#     apiVersion: autoscaling/v2beta2
#   fieldref:
#     fieldpath: spec.minReplicas
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/prod/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

nameprefix: prodpfx-
namespace: prodns

bases:
- ../../base

patchesStrategicMerge:
- horizontal_pod_autoscaler.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

nameprefix: stagingpfx-
namespace: stagingns

bases:
- ../../base

patchesStrategicMerge:
- horizontal_pod_autoscaler.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: puppetserver
  labels:
    app: puppetserver
spec:
  selector:
    matchLabels:
      app: puppetserver
  replicas: $(HorizontalPodAutoscaler.puppetserver.spec.minReplicas)
  template:
    metadata:
      labels:
        app: puppetserver
    spec:
      containers:
      - name: main
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
        - name: pupperserver
          containerPort: 8081
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/horizontal_pod_autoscaler.yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: puppetserver
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: puppetserver
  minReplicas: 1
  maxReplicas: 15
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
  - type: Pods
    pods:
      metric:
        name: packets-per-second
      target:
        type: AverageValue
        averageValue: 1k
  - type: Object
    object:
      metric:
        name: requests-per-second
      describedObject:
        apiVersion: networking.k8s.io/v1beta1
        kind: Ingress
        name: main-route
        # namespace: default
      target:
        type: Value
        value: 10k
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/ingress.yaml
kind: Ingress
apiVersion: networking.k8s.io/v1beta1
metadata:
  name: main-route
spec:
  backend:
    serviceName: puppetserver
    servicePort: 80
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/name_references.yaml
nameReference:
- kind: Ingress
  group: networking.k8s.io
  fieldSpecs:
  - path: spec/metrics/object/describedObject/name
    kind: HorizontalPodAutoscaler
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/var_references.yaml
varReference:
- path: spec/replicas
  kind: Deployment
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: puppetserver
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: 8081
    protocol: TCP
  selector:
    app: puppetserver
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/prod/horizontal_pod_autoscaler.yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: puppetserver
spec:
  maxReplicas: 15
  minReplicas: 10
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/horizontal_pod_autoscaler.yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: puppetserver
spec:
  maxReplicas: 6
  minReplicas: 3
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlay/staging -o ${DEMO_HOME}/actual/staging.yaml
kustomize build ${DEMO_HOME}/overlay/prod -o ${DEMO_HOME}/actual/prod.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod.yaml
apiVersion: v1
kind: Service
metadata:
  name: prodpfx-puppetserver
  namespace: prodns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: puppetserver
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: puppetserver
  name: prodpfx-puppetserver
  namespace: prodns
spec:
  replicas: 10
  selector:
    matchLabels:
      app: puppetserver
  template:
    metadata:
      labels:
        app: puppetserver
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: main
        ports:
        - containerPort: 8081
          name: pupperserver
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
---
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: prodpfx-puppetserver
  namespace: prodns
spec:
  maxReplicas: 15
  metrics:
  - resource:
      name: cpu
      target:
        averageUtilization: 50
        type: Utilization
    type: Resource
  - pods:
      metric:
        name: packets-per-second
      target:
        averageValue: 1k
        type: AverageValue
    type: Pods
  - object:
      describedObject:
        apiVersion: networking.k8s.io/v1beta1
        kind: Ingress
        name: prodpfx-main-route
      metric:
        name: requests-per-second
      target:
        type: Value
        value: 10k
    type: Object
  minReplicas: 10
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: prodpfx-puppetserver
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: prodpfx-main-route
  namespace: prodns
spec:
  backend:
    serviceName: prodpfx-puppetserver
    servicePort: 80
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
kind: Service
metadata:
  name: stagingpfx-puppetserver
  namespace: stagingns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: puppetserver
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: puppetserver
  name: stagingpfx-puppetserver
  namespace: stagingns
spec:
  replicas: 3
  selector:
    matchLabels:
      app: puppetserver
  template:
    metadata:
      labels:
        app: puppetserver
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: main
        ports:
        - containerPort: 8081
          name: pupperserver
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
---
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: stagingpfx-puppetserver
  namespace: stagingns
spec:
  maxReplicas: 6
  metrics:
  - resource:
      name: cpu
      target:
        averageUtilization: 50
        type: Utilization
    type: Resource
  - pods:
      metric:
        name: packets-per-second
      target:
        averageValue: 1k
        type: AverageValue
    type: Pods
  - object:
      describedObject:
        apiVersion: networking.k8s.io/v1beta1
        kind: Ingress
        name: stagingpfx-main-route
      metric:
        name: requests-per-second
      target:
        type: Value
        value: 10k
    type: Object
  minReplicas: 3
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: stagingpfx-puppetserver
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: stagingpfx-main-route
  namespace: stagingns
spec:
  backend:
    serviceName: stagingpfx-puppetserver
    servicePort: 80
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

