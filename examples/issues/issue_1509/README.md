# Feature Test for Issue 1509


This folder contains files describing how to address [Issue 1509](https://github.com/kubernetes-sigs/kustomize/issues/1509)

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
mkdir -p ${DEMO_HOME}/cfg-file-as-transformer
mkdir -p ${DEMO_HOME}/cfg-folder-as-resource
mkdir -p ${DEMO_HOME}/cfg-folder-as-transformer
mkdir -p ${DEMO_HOME}/commonconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- deployment.yaml
- role.yaml
- service.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cfg-file-as-transformer/kustomization.yaml
resources:
- ../base
transformers:
- ../commonconfig/prefixer.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cfg-folder-as-resource/kustomization.yaml
resources:
- ../base
- ../commonconfig
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cfg-folder-as-transformer/kustomization.yaml
resources:
- ../base
transformers:
- ../commonconfig
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonconfig/kustomization.yaml
transformers:
- prefixer.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myDeployment
spec:
  template:
    metadata:
      labels:
        backend: awesome
    spec:
      containers:
      - name: whatever
        image: whatever
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/role.yaml
apiVersion: v1
kind: Role
metadata:
  name: myRole
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myService
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonconfig/prefixer.yaml
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: customPrefixer
prefix: zzz-
fieldSpecs:
- kind: Deployment
  path: metadata/name
- kind: Service
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


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

