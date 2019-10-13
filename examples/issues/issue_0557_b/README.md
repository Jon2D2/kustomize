# Feature Test for Issue 0557


This folder contains files describing how to address [Issue 0557](https://github.com/kubernetes-sigs/kustomize/issues/0557)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/both
mkdir -p ${DEMO_HOME}/overlays/demo
mkdir -p ${DEMO_HOME}/overlays/integration
```

### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - serviceaccount.yml
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/serviceaccount.yml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: faros
EOF
```


### Preparation Step Other2

<!-- @createOther2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/both/kustomization.yml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - ../demo
  - ../integration
EOF
```


### Preparation Step Other3

<!-- @createOther3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/demo/crb.yml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: faros-clustergittrackobjects-viewer--apps-demo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: faros-clustergittrackobjects-viewer
subjects:
  - kind: ServiceAccount
    name: faros
    namespace: apps-demo
EOF
```


### Preparation Step Other4

<!-- @createOther4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/demo/kustomization.yml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: apps-demo

bases:
  - ../../base
resources:
  - crb.yml
EOF
```


### Preparation Step Other5

<!-- @createOther5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/integration/crb.yml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: faros-clustergittrackobjects-viewer--apps-integration
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: faros-clustergittrackobjects-viewer
subjects:
  - kind: ServiceAccount
    name: faros
    namespace: apps-integration
EOF
```


### Preparation Step Other6

<!-- @createOther6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/integration/kustomization.yml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: apps-integration

bases:
  - ../../base
resources:
  - crb.yml
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/both -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps-demo_~g_v1_serviceaccount_faros.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: faros
  namespace: apps-demo
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps-integration_~g_v1_serviceaccount_faros.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: faros
  namespace: apps-integration
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1_clusterrolebinding_faros-clustergittrackobjects-viewer--apps-demo.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: faros-clustergittrackobjects-viewer--apps-demo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: faros-clustergittrackobjects-viewer
subjects:
- kind: ServiceAccount
  name: faros
  namespace: apps-demo
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1_clusterrolebinding_faros-clustergittrackobjects-viewer--apps-integration.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: faros-clustergittrackobjects-viewer--apps-integration
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: faros-clustergittrackobjects-viewer
subjects:
- kind: ServiceAccount
  name: faros
  namespace: apps-integration
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

