# Feature Test for Issue 1301


This folder contains files describing how to address [Issue 1301](https://github.com/kubernetes-sigs/kustomize/issues/1301)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: default

resources:
- test-deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../base/

configMapGenerator:
- name: cafe-configmap
  literals:
  - FOO=BAR
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/test-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  labels:
     app: component2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: component2
  template:
    metadata:
      labels:
        app: component2
    spec:
      containers:
      - name: component2
        image: k8s.gcr.io/busybox
        command: [ "/bin/sh", "-c", "cat /etc/config/component2 && sleep 60" ]
        volumeMounts:
        - name: config-volume
          mountPath: /etc/config
      volumes:
      - name: config-volume
        configMap:
          name: cafe-configmap
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_test-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: component2
  name: test-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: component2
  template:
    metadata:
      labels:
        app: component2
    spec:
      containers:
      - command:
        - /bin/sh
        - -c
        - cat /etc/config/component2 && sleep 60
        image: k8s.gcr.io/busybox
        name: component2
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      volumes:
      - configMap:
          name: cafe-configmap-bm6m88fk92
        name: config-volume
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_cafe-configmap-bm6m88fk92.yaml
apiVersion: v1
data:
  FOO: BAR
kind: ConfigMap
metadata:
  name: cafe-configmap-bm6m88fk92
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

