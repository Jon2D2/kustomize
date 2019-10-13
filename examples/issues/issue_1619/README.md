# Feature Test for Issue 1619


This folder contains files describing how to address [Issue 1619](https://github.com/kubernetes-sigs/kustomize/issues/1619)

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

namespace: new-cert-namespace

resources:
- resources.yaml

configurations:
- kustomizeconfig.yaml

EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
nameReference:
- kind: Service
  version: v1
  fieldSpecs:
  - path: spec/service
    kind: APIService
    group: apiregistration.k8s.io
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resources.yaml
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: cert-manager
  name: cert-manager-webhook
  namespace: cert-manager
spec:
  ports:
  - name: https
    port: 443
    targetPort: 6443
  selector:
    app: cert-manager
  type: ClusterIP
---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  labels:
    app: cert-manager
  name: v1beta1.admission.certmanager.k8s.io
spec:
  group: admission.certmanager.k8s.io
  groupPriorityMinimum: 1000
  service:
    name: cert-manager-webhook
    namespace: cert-manager
  version: v1beta1
  versionPriority: 15
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
cat <<'EOF' >${DEMO_HOME}/expected/apiregistration.k8s.io_v1beta1_apiservice_v1beta1.admission.certmanager.k8s.io.yaml
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  labels:
    app: cert-manager
  name: v1beta1.admission.certmanager.k8s.io
spec:
  group: admission.certmanager.k8s.io
  groupPriorityMinimum: 1000
  service:
    name: cert-manager-webhook
    namespace: new-cert-namespace
  version: v1beta1
  versionPriority: 15
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_cert-manager-webhook.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: cert-manager
  name: cert-manager-webhook
  namespace: new-cert-namespace
spec:
  ports:
  - name: https
    port: 443
    targetPort: 6443
  selector:
    app: cert-manager
  type: ClusterIP
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

