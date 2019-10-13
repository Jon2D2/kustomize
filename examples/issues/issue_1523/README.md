# Feature Test for Issue 1523


This folder contains files describing how to address [Issue 1523](https://github.com/kubernetes-sigs/kustomize/issues/1523)

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
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/psp
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- cluster-role.yaml
commonLabels:
  kustomize.component: argo

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/psp/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
patchesJson6902:
- target:
    group: rbac.authorization.k8s.io
    version: v1beta1
    kind: ClusterRole
    name: argo
  path: patch.yaml

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/cluster-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    app: argo
  name: argo
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - pods/exec
  verbs:
  - create
  - get
  - list
  - watch
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - watch
  - list
- apiGroups:
  - ""
  resources:
  - persistentvolumeclaims
  verbs:
  - create
  - delete
- apiGroups:
  - argoproj.io
  resources:
  - workflows
  verbs:
  - get
  - list
  - watch
  - update
  - patch

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/psp/patch.yaml
- op: add
  path: /rules/0
  value:
    apiGroups:
      - extensions
    resourceNames:
      - privileged
    resources:
      - podsecuritypolicies
    verbs:
      - use
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/psp -o ${DEMO_HOME}/actual/result.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/result.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    app: argo
    kustomize.component: argo
  name: argo
rules:
- apiGroups:
  - extensions
  resourceNames:
  - privileged
  resources:
  - podsecuritypolicies
  verbs:
  - use
- apiGroups:
  - ""
  resources:
  - pods
  - pods/exec
  verbs:
  - create
  - get
  - list
  - watch
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - watch
  - list
- apiGroups:
  - ""
  resources:
  - persistentvolumeclaims
  verbs:
  - create
  - delete
- apiGroups:
  - argoproj.io
  resources:
  - workflows
  verbs:
  - get
  - list
  - watch
  - update
  - patch
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

