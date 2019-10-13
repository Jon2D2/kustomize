# Feature Test for Issue 1264


This folder contains files describing how to address [Issue 1264](https://github.com/kubernetes-sigs/kustomize/issues/1264)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/apps
mkdir -p ${DEMO_HOME}/apps/prometheus
mkdir -p ${DEMO_HOME}/ca-1-dev
mkdir -p ${DEMO_HOME}/common
mkdir -p ${DEMO_HOME}/common/istio
mkdir -p ${DEMO_HOME}/common/istio/istio
mkdir -p ${DEMO_HOME}/common/prometheus
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/apps/prometheus/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: prometheus

namespace: prometheus

resources:
  - rbac.yaml

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/ca-1-dev/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization


commonLabels:
  cluster: ca-1-dev

resources:
  - ../common/istio/istio
  - ../common/prometheus
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/istio/istio/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - istio.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/prometheus/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: prometheus
commonLabels:
  app: prometheus

resources:
  - ../../apps/prometheus

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/apps/prometheus/rbac.yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus


  # If this is removed, kustomize works
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: prometheus
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus
subjects:
  - kind: ServiceAccount
    name: prometheus
    namespace: prometheus
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/istio/istio/istio.yaml
# Source: istio/charts/prometheus/templates/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus
  namespace: istio-system
  labels:
    app: prometheus
    chart: prometheus
    heritage: Tiller
    release: istio
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/ca-1-dev -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/istio-system_~g_v1_serviceaccount_prometheus.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: prometheus
    chart: prometheus
    cluster: ca-1-dev
    heritage: Tiller
    release: istio
  name: prometheus
  namespace: istio-system
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus_~g_v1_serviceaccount_prometheus.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: prometheus
    cluster: ca-1-dev
  name: prometheus
  namespace: prometheus
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1beta1_clusterrolebinding_prometheus.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    app: prometheus
    cluster: ca-1-dev
  name: prometheus
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus
subjects:
- kind: ServiceAccount
  name: prometheus
  namespace: prometheus
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

