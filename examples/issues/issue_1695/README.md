# Feature Test for Issue 1695


This folder contains files describing how to address [Issue 1695](https://github.com/kubernetes-sigs/kustomize/issues/1695)

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
resources:
- deployment.yaml
- secret.yaml
# nameSuffix: -service
transformers:
- suffixer.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: group
spec:
  selector:
    matchLabels:
      app: group
  template:
    metadata:
      labels:
        app: group
    spec:
      containers:
        - name: group
          image: image
          envFrom:
            - secretRef:
                name: group
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: group
data:
  secretKey: c2VjcmV0VmFsdWUK
type: Opaque
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/suffixer.yaml
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: suffix
suffix: -service
fieldSpecs:
- kind: Secret
  path: metadata/name
- kind: Deployment
  path: metadata/name
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_group-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: group-service
spec:
  selector:
    matchLabels:
      app: group
  template:
    metadata:
      labels:
        app: group
    spec:
      containers:
      - envFrom:
        - secretRef:
            name: group-service
        image: image
        name: group
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_group-service.yaml
apiVersion: v1
data:
  secretKey: c2VjcmV0VmFsdWUK
kind: Secret
metadata:
  name: group-service
type: Opaque
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

