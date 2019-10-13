# Feature Test for Issue 1713


This folder contains files describing how to address [Issue 1713](https://github.com/kubernetes-sigs/kustomize/issues/1713)

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

### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/namespace.yaml
varReference:
- path: metadata/name
  kind: Namespace
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: $(ServiceAccount.my-service-account.metadata.namespace)
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/other.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account
  namespace: some-namespace
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- namespace.yaml
- other.yaml

# This kustomization.yaml leverages auto-var feature
# instead of doing it manually
# Uncommment if you do not have access to the feature.
# vars:
# - name: ServiceAccount.my-service-account.metadata.namespace
#   objref:
#     kind: ServiceAccount
#     name: my-service-account
#     apiVersion: v1
#   fieldref:
#     fieldpath: metadata.namespace
 
# configurations:
# - kustomizeconfig/namespace.yaml
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_namespace_some-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: some-namespace
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_serviceaccount_my-service-account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account
  namespace: some-namespace
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

