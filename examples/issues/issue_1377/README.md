# Feature Test for Issue 1377


This folder contains files describing how to address [Issue 1377](https://github.com/kubernetes-sigs/kustomize/issues/1377)

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

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: foo
resources:
- test.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test.yaml
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: crb
subjects:
- kind: ServiceAccount
  name: sa
  namespace: bar
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sa
  namespace: bar
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
  namespace: foo # this configmap has the same namespace as the transformer target
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_cm.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
  namespace: foo
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_serviceaccount_sa.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sa
  namespace: foo
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1beta1_clusterrolebinding_crb.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: crb
subjects:
- kind: ServiceAccount
  name: sa
  namespace: foo
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

