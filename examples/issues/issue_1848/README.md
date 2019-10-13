# Feature Test for Issue 1848


This folder contains files describing how to address [Issue 1848](https://github.com/kubernetes-sigs/kustomize/issues/1848)

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
mkdir -p ${DEMO_HOME}/base/proxy
mkdir -p ${DEMO_HOME}/base/proxy/config
mkdir -p ${DEMO_HOME}/environments
mkdir -p ${DEMO_HOME}/environments/test
mkdir -p ${DEMO_HOME}/environments/test/proxy
mkdir -p ${DEMO_HOME}/environments/test/proxy/config
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
  - proxy
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/proxy/kustomization.yaml
resources:
  - deployment.yaml
  - service.yaml

configMapGenerator:
  - name: proxy-config
    files:
      - config/dhparam.pem
      - config/chained.crt
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environments/test/kustomization.yaml
namespace: test

resources:
  - ../../base
  - proxy
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environments/test/proxy/kustomization.yaml
resources:
- ../../../base/proxy

configMapGenerator:
  - name: proxy-config
    behavior: merge
    files:
      - config/nginx.conf
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
  # testing environment specializations
  - environments/test
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/proxy/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    metadata:
      labels:
        backend: awesome
    spec:
      containers:
      - name: whatever
        image: whatever
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/proxy/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/proxy/config/chained.crt
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/proxy/config/dhparam.pem
EOF
```


### Preparation Step Other2

<!-- @createOther2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environments/test/proxy/config/nginx.conf
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_my-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
  namespace: test
spec:
  template:
    metadata:
      labels:
        backend: awesome
    spec:
      containers:
      - image: whatever
        name: whatever
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_proxy-config-7t5mhk6kg7.yaml
apiVersion: v1
data:
  chained.crt: ""
  dhparam.pem: ""
  nginx.conf: ""
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: proxy-config-7t5mhk6kg7
  namespace: test
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_my-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: test
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

