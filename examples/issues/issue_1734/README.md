# Feature Test for Issue 1734


This folder contains files describing how to address [Issue 1734](https://github.com/kubernetes-sigs/kustomize/issues/1734)

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
mkdir -p ${DEMO_HOME}/base/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- nodejsapp-service.yaml

configMapGenerator:
- name: nodejs-test-cm
  literals:
    - branch=develop
    - build=20191025
    - external_ip=10.0.0.3
    - external_port=3000

# This kustomization.yaml leverages auto-var feature
# instead of doing it manually
# Uncommment if you do not have access to the feature.
# configurations:
# - kustomizeconfig/var_references.yaml

# vars:
# - name: ConfigMap.nodejs-test-cm.data.external_ip
#   objref:
#     kind: ConfigMap
#     name: nodejs-test-cm
#     apiVersion: v1
#   fieldref:
#     fieldpath: data.external_ip
# - name: ConfigMap.nodejs-test-cm.data.external_port
#   objref:
#     kind: ConfigMap
#     name: nodejs-test-cm
#     apiVersion: v1
#   fieldref:
#     fieldpath: data.external_port
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/var_references.yaml
varReference:
- kind: Service
  path: spec/externalIPs

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/nodejsapp-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: "nodejs-test"
  labels:
    service: "nodejs-test"
spec:
  ports:
  - name: client-connect
    port: "$(ConfigMap.nodejs-test-cm.data.external_port)"
    targetPort: 3000
  selector:
    app: "nodejs-test"
  type: ClusterIP
  externalIPs:
  - "$(ConfigMap.nodejs-test-cm.data.external_ip)"
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/base -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_nodejs-test-cm-tt7btc77cb.yaml
apiVersion: v1
data:
  branch: develop
  build: "20191025"
  external_ip: 10.0.0.3
  external_port: "3000"
kind: ConfigMap
metadata:
  name: nodejs-test-cm-tt7btc77cb
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_nodejs-test.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    service: nodejs-test
  name: nodejs-test
spec:
  externalIPs:
  - 10.0.0.3
  ports:
  - name: client-connect
    port: "3000"
    targetPort: 3000
  selector:
    app: nodejs-test
  type: ClusterIP
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

