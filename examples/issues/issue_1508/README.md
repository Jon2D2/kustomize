# Feature Test for Issue 1508


This folder contains files describing how to address [Issue 1508](https://github.com/kubernetes-sigs/kustomize/issues/1508)

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
- role.yaml
- service.yaml
transformers:
- prefixer.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
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
cat <<'EOF' >${DEMO_HOME}/prefixer.yaml
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


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/role.yaml
apiVersion: v1
kind: Role
metadata:
  name: myRole
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myService
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME} --enable_alpha_plugins -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_zzz-mydeployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zzz-myDeployment
spec:
  template:
    metadata:
      labels:
        backend: awesome
    spec:
      containers:
      - image: whatever
        name: whatever
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_role_myrole.yaml
apiVersion: v1
kind: Role
metadata:
  name: myRole
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_zzz-myservice.yaml
apiVersion: v1
kind: Service
metadata:
  name: zzz-myService
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

