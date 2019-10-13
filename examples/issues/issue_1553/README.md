# Feature Test for Issue 1553


This folder contains files describing how to address [Issue 1553](https://github.com/kubernetes-sigs/kustomize/issues/1553)

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
mkdir -p ${DEMO_HOME}/overlays/staging
mkdir -p ${DEMO_HOME}/overlays/staging/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: sentry

resources:
- ./deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base

secretGenerator:
- name: core
  type: Opaque
  envs:
  - secrets.txt

patchesStrategicMerge:
- patch.yaml

# ##########################################################
# Uncomment this section if you don't have this PR: 
# Automatic Creation of 'vars:' and 'varReferences:' sections #1217
# https://github.com/kubernetes-sigs/kustomize/pull/1217
# ##########################################################
 
# vars:
# - name: Secret.core.metadata.name
#   objref:
#     kind: Secret
#     name: core
#     apiVersion: v1
#   fieldref:
#     fieldpath: metadata.name
# 
# configurations:
# - ./kustomizeconfig/varreference.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: core
spec:
  template:
    spec:
      containers:
      - name: sentry
        imagePullPolicy: Always
        image: sentryimage:1.0
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/kustomizeconfig/varreference.yaml
varReference:
- path: spec/template/spec/containers[]/env[]/valueFrom/secretKeyRef/name
  kind: Deployment
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: core
spec:
  template:
    spec:
      containers:
      - name: sentry
        env:
        - name: SENTRY_DNS
          valueFrom:
            secretKeyRef:
              name: $(Secret.core.metadata.name)
              key: sentryDNS
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/secrets.txt
username=theuser
password=thepassword
sentryDNS=127.0.0.1
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/staging -o ${DEMO_HOME}/actual/staging.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
data:
  password: dGhlcGFzc3dvcmQ=
  sentryDNS: MTI3LjAuMC4x
  username: dGhldXNlcg==
kind: Secret
metadata:
  name: core-966hh9bc88
type: Opaque
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: sentry
  name: core
spec:
  selector:
    matchLabels:
      app: sentry
  template:
    metadata:
      labels:
        app: sentry
    spec:
      containers:
      - env:
        - name: SENTRY_DNS
          valueFrom:
            secretKeyRef:
              key: sentryDNS
              name: core-966hh9bc88
        image: sentryimage:1.0
        imagePullPolicy: Always
        name: sentry
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

