# Feature Test for Issue 1155


This folder contains files describing how to address [Issue 1155](https://github.com/kubernetes-sigs/kustomize/issues/1155)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
configMapGenerator:
- name: test
  namespace: default
  literals:
    - key=value
- name: test
  namespace: kube-system
  literals:
    - key=value
secretGenerator:
- name: test
  namespace: default
  literals:
  - username=admin
  - password=somepw
- name: test
  namespace: kube-system
  literals:
  - username=admin
  - password=somepw
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
cat <<'EOF' >${DEMO_HOME}/expected/default_~g_v1_secret_test-h65t9hg6kc.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: test-h65t9hg6kc
  namespace: default
type: Opaque
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default_~g_v1_configmap_test-t5t4md8fdm.yaml
apiVersion: v1
data:
  key: value
kind: ConfigMap
metadata:
  name: test-t5t4md8fdm
  namespace: default
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kube-system_~g_v1_configmap_test-t5t4md8fdm.yaml
apiVersion: v1
data:
  key: value
kind: ConfigMap
metadata:
  name: test-t5t4md8fdm
  namespace: kube-system
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kube-system_~g_v1_secret_test-h65t9hg6kc.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: test-h65t9hg6kc
  namespace: kube-system
type: Opaque
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

