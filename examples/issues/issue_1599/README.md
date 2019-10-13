# Feature Test for Issue 1599


This folder contains files describing how to address [Issue 1599](https://github.com/kubernetes-sigs/kustomize/issues/1599)

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
mkdir -p ${DEMO_HOME}/namespaces
mkdir -p ${DEMO_HOME}/rbac-my-namespace
mkdir -p ${DEMO_HOME}/rbac-my-other-namespace
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
   app: my-label

resources:
- ./namespaces
- ./rbac-my-namespace
- ./rbac-my-other-namespace
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/namespaces/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ./my-namespace.yaml
- ./my-other-namespace.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-namespace/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# There is a real bug here in kustomize
# rolebinding will not be updated properly
# namePrefix: pfx1-

namespace: my-namespace

resources:
- ./service-account.yaml
- ./rolebinding.yaml
- ./role.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-other-namespace/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# There is a real bug here in kustomize
# rolebinding will not be updated properly
# namePrefix: pfx2-

namespace: my-other-namespace

resources:
- ./service-account.yaml
- ./rolebinding.yaml
- ./role.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/namespaces/my-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/namespaces/my-other-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-other-namespace
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-namespace/rolebinding.yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: my-rolebinding
subjects:
  - kind: ServiceAccount
    name: default
    namespace: will-be-replaced-because-created-automatically-with-the-namespace
  - kind: ServiceAccount
    name: my-namespace-sa
  - kind: ServiceAccount
    name: my-other-namespace-sa
    namespace: my-other-namespace
  - kind: User
    apiGroup: rbac.authorization.k8s.io
    name: user-1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: my-role
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-namespace/role.yaml
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: my-role
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-namespace/service-account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-namespace-sa
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-other-namespace/rolebinding.yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: my-rolebinding
subjects:
  - kind: ServiceAccount
    name: default
    namespace: will-be-replaced-because-created-automatically-with-the-namespace
  - kind: ServiceAccount
    name: my-other-namespace-sa
  - kind: ServiceAccount
    name: my-namespace-sa
    namespace: my-namespace
  - kind: User
    apiGroup: rbac.authorization.k8s.io
    name: user-1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: my-role
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-other-namespace/role.yaml
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: my-role
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rbac-my-other-namespace/service-account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-other-namespace-sa
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_namespace_my-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: my-label
  name: my-namespace
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_namespace_my-other-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: my-label
  name: my-other-namespace
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/my-namespace_~g_v1_serviceaccount_my-namespace-sa.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: my-label
  name: my-namespace-sa
  namespace: my-namespace
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/my-namespace_rbac.authorization.k8s.io_v1beta1_role_my-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app: my-label
  name: my-role
  namespace: my-namespace
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/my-namespace_rbac.authorization.k8s.io_v1_rolebinding_my-rolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    app: my-label
  name: my-rolebinding
  namespace: my-namespace
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: my-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: my-namespace
- kind: ServiceAccount
  name: my-namespace-sa
  namespace: my-namespace
- kind: ServiceAccount
  name: my-other-namespace-sa
  namespace: my-other-namespace
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: user-1
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/my-other-namespace_~g_v1_serviceaccount_my-other-namespace-sa.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: my-label
  name: my-other-namespace-sa
  namespace: my-other-namespace
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/my-other-namespace_rbac.authorization.k8s.io_v1beta1_role_my-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app: my-label
  name: my-role
  namespace: my-other-namespace
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/my-other-namespace_rbac.authorization.k8s.io_v1_rolebinding_my-rolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    app: my-label
  name: my-rolebinding
  namespace: my-other-namespace
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: my-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: my-other-namespace
- kind: ServiceAccount
  name: my-other-namespace-sa
  namespace: my-other-namespace
- kind: ServiceAccount
  name: my-namespace-sa
  namespace: my-namespace
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: user-1
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

