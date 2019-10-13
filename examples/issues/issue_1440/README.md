# Feature Test for Issue 1440


This folder contains files describing how to address [Issue 1440](https://github.com/kubernetes-sigs/kustomize/issues/1440)

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
mkdir -p ${DEMO_HOME}/intermediate-base
mkdir -p ${DEMO_HOME}/overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- deploy.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/intermediate-base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: frontend-web-
namespace: sandbox

resources:
- ../base

patchesStrategicMerge:
- patch-deploy-env.yaml

configMapGenerator:
- envs:
  - env.properties
  namespace: sandbox
  name: env
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../intermediate-base

namePrefix: production-
namespace: sandbox

configMapGenerator:
- envs:
  - env.properties
  name: env
  namespace: sandbox
  behavior: merge
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deploy.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 1
  template:
    spec:
      containers:
      - args: []
        command: []
        image: container-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
        name: main
        ports:
        - containerPort: 8000
          name: http
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/intermediate-base/patch-deploy-env.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
        - name: main
          envFrom:
          - configMapRef:
              name: env
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/intermediate-base/env.properties
DEBUG=False
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/env.properties
DEBUG=True
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/base -o ${DEMO_HOME}/actual/base.yaml
kustomize build ${DEMO_HOME}/intermediate-base -o ${DEMO_HOME}/actual/intermediate-base.yaml
kustomize build ${DEMO_HOME}/overlay -o ${DEMO_HOME}/actual/overlay.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/base.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 1
  template:
    spec:
      containers:
      - args: []
        command: []
        image: container-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
        name: main
        ports:
        - containerPort: 8000
          name: http
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/intermediate-base.yaml
apiVersion: v1
data:
  DEBUG: "False"
kind: ConfigMap
metadata:
  name: frontend-web-env-28m45kmmm8
  namespace: sandbox
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend-web-app
  namespace: sandbox
spec:
  replicas: 1
  template:
    spec:
      containers:
      - args: []
        command: []
        envFrom:
        - configMapRef:
            name: frontend-web-env-28m45kmmm8
        image: container-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
        name: main
        ports:
        - containerPort: 8000
          name: http
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/overlay.yaml
apiVersion: v1
data:
  DEBUG: "True"
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: production-frontend-web-env-59bmdgdmkh
  namespace: sandbox
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: production-frontend-web-app
  namespace: sandbox
spec:
  replicas: 1
  template:
    spec:
      containers:
      - args: []
        command: []
        envFrom:
        - configMapRef:
            name: production-frontend-web-env-59bmdgdmkh
        image: container-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
        name: main
        ports:
        - containerPort: 8000
          name: http
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

