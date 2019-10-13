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
mkdir -p ${DEMO_HOME}/combined
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/a
mkdir -p ${DEMO_HOME}/overlays/b
mkdir -p ${DEMO_HOME}/overlays/c
mkdir -p ${DEMO_HOME}/overlays/d
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- serviceaccount.yaml
- rolebinding.yaml
- clusterrolebinding.yaml
- clusterrole.yaml
namePrefix: pfx-
nameSuffix: -sfx
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/combined/kustomization.yaml
resources:
- ../overlays/a
- ../overlays/b
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/a/kustomization.yaml
resources:
- ../../base
namePrefix: a-
nameSuffix: -suffixA
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/b/kustomization.yaml
resources:
- ../../base
namePrefix: b-
nameSuffix: -suffixB
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/c/kustomization.yaml
resources:
- ../a
namePrefix: c-
nameSuffix: -suffixC
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/d/kustomization.yaml
resources:
- ../b
namePrefix: d-
nameSuffix: -suffixD
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/clusterrolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: role
subjects:
- kind: ServiceAccount
  name: serviceaccount
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: role
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/rolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: role
subjects:
- kind: ServiceAccount
  name: serviceaccount
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: serviceaccount
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/combined -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_serviceaccount_a-pfx-serviceaccount-sfx-suffixa.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: a-pfx-serviceaccount-sfx-suffixA
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_serviceaccount_b-pfx-serviceaccount-sfx-suffixb.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: b-pfx-serviceaccount-sfx-suffixB
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1beta1_clusterrolebinding_a-pfx-rolebinding-sfx-suffixa.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: a-pfx-rolebinding-sfx-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: a-pfx-role-sfx-suffixA
subjects:
- kind: ServiceAccount
  name: a-pfx-serviceaccount-sfx-suffixA
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1beta1_clusterrolebinding_b-pfx-rolebinding-sfx-suffixb.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: b-pfx-rolebinding-sfx-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: b-pfx-role-sfx-suffixB
subjects:
- kind: ServiceAccount
  name: b-pfx-serviceaccount-sfx-suffixB
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1beta1_rolebinding_a-pfx-rolebinding-sfx-suffixa.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: a-pfx-rolebinding-sfx-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: a-pfx-role-sfx-suffixA
subjects:
- kind: ServiceAccount
  name: a-pfx-serviceaccount-sfx-suffixA
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1beta1_rolebinding_b-pfx-rolebinding-sfx-suffixb.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: b-pfx-rolebinding-sfx-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: b-pfx-role-sfx-suffixB
subjects:
- kind: ServiceAccount
  name: b-pfx-serviceaccount-sfx-suffixB
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1_clusterrole_a-pfx-role-sfx-suffixa.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: a-pfx-role-sfx-suffixA
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1_clusterrole_b-pfx-role-sfx-suffixb.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: b-pfx-role-sfx-suffixB
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

