# Feature Test for Issue 0710


This folder contains files describing how to address [Issue 0710](https://github.com/kubernetes-sigs/kustomize/issues/0710)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base1
mkdir -p ${DEMO_HOME}/base2
mkdir -p ${DEMO_HOME}/overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base1/kustomization.yaml
resources:
  - app-deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base2/kustomization.yaml
namespace: appName
namePrefix: prefix-
commonLabels:
  app: appName

resources:
  - ../base1
  - app-other-service.yaml

patchesStrategicMerge:
  - app-deployment-patch.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
resources:
- ../base2

namespace: appName
commonLabels:
  app: appName

patchesStrategicMerge:
- app-other-deployment-patch.yaml

configMapGenerator:
- name: config
  files:
  - config-file.json
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base1/app-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
        - name: app
          #[REDACTED]
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base2/app-deployment-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
      - name: app
        env:
        - name: ANOTHERENV
          value: ANOTHERVALUE
      - name: anothercontainer
        image: anotherimage
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base2/app-other-service.yaml
apiVersion: apps/v1
kind: Service
metadata:
  name: app
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/app-other-deployment-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
        - name: app
          #[REDACTED]
          volumeMounts:
            - name: config-volume
              mountPath: "/path/to/file"
      volumes:
        - name: config-volume
          configMap:
            name: config
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/config-file.json
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlay -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_prefix-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: appName
  name: prefix-app
  namespace: appName
spec:
  selector:
    matchLabels:
      app: appName
  template:
    metadata:
      labels:
        app: appName
    spec:
      containers:
      - env:
        - name: ANOTHERENV
          value: ANOTHERVALUE
        name: app
        volumeMounts:
        - mountPath: /path/to/file
          name: config-volume
      - image: anotherimage
        name: anothercontainer
      volumes:
      - configMap:
          name: config-79tktd9hkb
        name: config-volume
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_service_prefix-app.yaml
apiVersion: apps/v1
kind: Service
metadata:
  labels:
    app: appName
  name: prefix-app
  namespace: appName
spec:
  selector:
    app: appName
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_config-79tktd9hkb.yaml
apiVersion: v1
data:
  config-file.json: ""
kind: ConfigMap
metadata:
  labels:
    app: appName
  name: config-79tktd9hkb
  namespace: appName
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

