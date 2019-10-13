# Feature Test for TestCase multibases-conflict


This folder contains files for old test-case multibases-conflict

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- serviceaccount.yaml
- rolebinding.yaml
namePrefix: base-
nameSuffix: -suffix
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
namePrefix: a-
nameSuffix: -suffixA

resources:
- ../../base/
- serviceaccount.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/b/kustomization.yaml
resources:
- ../../base/

namePrefix: b-
nameSuffix: -suffixB
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/rolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: role
subjects:
- kind: ServiceAccount
  name: serviceaccount
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: serviceaccount
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/a/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: serviceaccount
EOF
```

## Execution

<!-- @build @test -->
```bash
# kustomize build ${DEMO_HOME}/combined -o ${DEMO_HOME}/actual.yaml
```

## Verification


<!-- @compareActualToExpected @test -->
```bash
```

