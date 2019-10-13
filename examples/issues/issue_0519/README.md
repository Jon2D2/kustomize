# Feature Test for Issue 0519


This folder contains files describing how to address [Issue 0519](https://github.com/kubernetes-sigs/kustomize/issues/0519)

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
mkdir -p ${DEMO_HOME}/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
namePrefix: p1-
nameSuffix: -s1
resources:
- resources.yaml
configurations:
- ./kustomizeconfig/skip-secrets-prefix.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/skip-secrets-prefix.yaml
namePrefix:
- path: metadata/name
  kind: CustomResourceDefinition
  skip: true
- path: metadata/name
  apiVersion: v1
  kind: Secret
  skip: true
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resources.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
---
apiVersion: v1
kind: Secret
metadata:
  name: secret
---
kind: CustomResourceDefinition
metadata:
  name: my-crd
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME} -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_p1-cm1-s1.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: p1-cm1-s1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: secret
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_~v_customresourcedefinition_my-crd.yaml
kind: CustomResourceDefinition
metadata:
  name: my-crd
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

