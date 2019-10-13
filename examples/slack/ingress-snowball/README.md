# Feature Test for Issue ingress-snowball

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
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/production
mkdir -p ${DEMO_HOME}/overlays/staging
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- ./snowball-webclient-ip.yaml
- ./snowball-webclient-certificate.yaml
- ./snowball-webclient-deployment.yaml
- ./snowball-webclient-ingress.yaml
- ./snowball-webclient-service.yaml

configurations:
- ./kustomizeconfig/ingress.yaml

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/kustomization.yaml
resources:
- ../../base
namePrefix: production-
patchesStrategicMerge:
- snowball-webclient-certificate-patch.yaml
- snowball-webclient-ingress-patch.yaml

images:
- name: gcr.io/myapp/snowball-webclient
  newTag: 2.0.0


EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/kustomization.yaml
resources:
- ../../base
namePrefix: staging-
patchesStrategicMerge:
- snowball-webclient-certificate-patch.yaml
- snowball-webclient-ingress-patch.yaml

images:
- name: gcr.io/myapp/snowball-webclient
  newTag: 3.0.0

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/ingress.yaml
nameReference:
- kind: ManagedCertificate
  group: networking.gke.io
  version: v1beta1
  fieldSpecs:
  - path: metadata/annotations/networking.gke.io\/managed-certificates
    kind: Ingress

- kind: GKEGlobalStaticIP
  fieldSpecs:
  - path: metadata/annotations/kubernetes.io\/ingress.global-static-ip-name
    kind: Ingress
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/snowball-webclient-certificate.yaml
apiVersion: networking.gke.io/v1beta1
kind: ManagedCertificate
metadata:
  name: snowball-webclient-certificate
spec:
  domains:
    - $()

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/snowball-webclient-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: snowball-webclient
  name: snowball-webclient
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: snowball-webclient
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: snowball-webclient
    spec:
      containers:
      - name: snowball-webclient
        image: gcr.io/myapp/snowball-webclient:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 3000
          protocol: TCP
        imagePullPolicy: IfNotPresent
        
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/snowball-webclient-ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: snowball-webclient-ingress
  annotations:
    kubernetes.io/ingress.global-static-ip-name: snowball-webclient-ip
    networking.gke.io/managed-certificates: snowball-webclient-certificate
  labels:
    app: snowball-webclient
spec:
  backend:
    serviceName: snowball-webclient-service
    servicePort: 80

EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/snowball-webclient-ip.yaml
kind: GKEGlobalStaticIP
metadata:
  name: snowball-webclient-ip
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/snowball-webclient-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: snowball-webclient-service
  labels:
    app: snowball-webclient-service
spec:
  type: NodePort
  selector:
    app: snowball-webclient
  ports:
    - port: 80
      targetPort: 3000

EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/snowball-webclient-certificate-patch.yaml
apiVersion: networking.gke.io/v1beta1
kind: ManagedCertificate
metadata:
  name: snowball-webclient-certificate
spec:
  domains:
    - public.prolificparc.com

EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/snowball-webclient-ingress-patch.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: snowball-webclient-ingress
spec:
  backend:
    serviceName: snowball-webclient-service
    servicePort: 80

EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/snowball-webclient-certificate-patch.yaml
apiVersion: networking.gke.io/v1beta1
kind: ManagedCertificate
metadata:
  name: snowball-webclient-certificate
spec:
  domains:
    - abc.prolificparc.com

EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/snowball-webclient-ingress-patch.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: snowball-webclient-ingress
spec:
  backend:
    serviceName: snowball-webclient-service
    servicePort: 80

EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/foo
./expected/staging.yaml
./expected/production.yaml
./README.md
./overlays/production/snowball-webclient-ingress-patch.yaml
./overlays/staging/snowball-webclient-ingress-patch.yaml
./base/snowball-webclient-deployment.yaml
./base/snowball-webclient-ingress.yaml
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/production -o ${DEMO_HOME}/actual/production.yaml
kustomize build ${DEMO_HOME}/overlays/staging -o ${DEMO_HOME}/actual/staging.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: snowball-webclient-service
  name: production-snowball-webclient-service
spec:
  ports:
  - port: 80
    targetPort: 3000
  selector:
    app: snowball-webclient
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: snowball-webclient
  name: production-snowball-webclient
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: snowball-webclient
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: snowball-webclient
    spec:
      containers:
      - image: gcr.io/myapp/snowball-webclient:2.0.0
        imagePullPolicy: IfNotPresent
        name: snowball-webclient
        ports:
        - containerPort: 3000
          protocol: TCP
---
apiVersion: networking.gke.io/v1beta1
kind: ManagedCertificate
metadata:
  name: production-snowball-webclient-certificate
spec:
  domains:
  - public.prolificparc.com
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: production-snowball-webclient-ip
    networking.gke.io/managed-certificates: production-snowball-webclient-certificate
  labels:
    app: snowball-webclient
  name: production-snowball-webclient-ingress
spec:
  backend:
    serviceName: production-snowball-webclient-service
    servicePort: 80
---
kind: GKEGlobalStaticIP
metadata:
  name: production-snowball-webclient-ip
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: snowball-webclient-service
  name: staging-snowball-webclient-service
spec:
  ports:
  - port: 80
    targetPort: 3000
  selector:
    app: snowball-webclient
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: snowball-webclient
  name: staging-snowball-webclient
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: snowball-webclient
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: snowball-webclient
    spec:
      containers:
      - image: gcr.io/myapp/snowball-webclient:3.0.0
        imagePullPolicy: IfNotPresent
        name: snowball-webclient
        ports:
        - containerPort: 3000
          protocol: TCP
---
apiVersion: networking.gke.io/v1beta1
kind: ManagedCertificate
metadata:
  name: staging-snowball-webclient-certificate
spec:
  domains:
  - abc.prolificparc.com
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: staging-snowball-webclient-ip
    networking.gke.io/managed-certificates: staging-snowball-webclient-certificate
  labels:
    app: snowball-webclient
  name: staging-snowball-webclient-ingress
spec:
  backend:
    serviceName: staging-snowball-webclient-service
    servicePort: 80
---
kind: GKEGlobalStaticIP
metadata:
  name: staging-snowball-webclient-ip
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

