# Feature Test for Issue 0249


This folder contains files describing how to address [Issue 0249](https://github.com/kubernetes-sigs/kustomize/issues/0249)

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
commonLabels:
  environment: staging
configurations:
- ./kustomizeconfig/commonlabels.yaml
resources:
- external-service.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/external-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: external-service
spec:
  ports:
    - protocol: TCP
      port: 80
---
apiVersion: v1
kind: Endpoints
metadata:
  name: external-service
subsets:
  - addresses:
      - ip: 192.168.60.1
      - ip: 192.168.60.2
    ports:
      - port: 80

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/commonlabels.yaml
commonLabels:
- path: spec/selector
  kind: Service
  version: v1
  behavior: remove
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_endpoints_external-service.yaml
apiVersion: v1
kind: Endpoints
metadata:
  labels:
    environment: staging
  name: external-service
subsets:
- addresses:
  - ip: 192.168.60.1
  - ip: 192.168.60.2
  ports:
  - port: 80
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_external-service.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    environment: staging
  name: external-service
spec:
  ports:
  - port: 80
    protocol: TCP
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

