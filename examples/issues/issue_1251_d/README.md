# Feature Test for Issue 1251


This folder contains files describing how to address [Issue 1251](https://github.com/kubernetes-sigs/kustomize/issues/1251)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/common
mkdir -p ${DEMO_HOME}/component1
mkdir -p ${DEMO_HOME}/component2
mkdir -p ${DEMO_HOME}/myapp
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/kustomization.yaml
resources:
- configmap.yaml

vars:
- name: ConfigMap.global.data.user
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: global
  fieldref:
    fieldpath: data.user
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component1/kustomization.yaml
resources:
- ../common
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component2/kustomization.yaml
resources:
- ../common
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/myapp/kustomization.yaml
resources:
- ../common
- ../component1
- ../component2
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: global
data:
  settings: |
     database: mydb
     port: 3000
  user: myuser
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component1/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: component1
  labels:
     app: component1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: component1
  template:
    metadata:
      labels:
        app: component1
    spec:
      containers:
      - name: component1
        image: k8s.gcr.io/busybox
        env:
        - name: APP_USER
          value: $(ConfigMap.global.data.user)
        command: [ "/bin/sh", "-c", "cat /etc/config/component1 && sleep 60" ]
        volumeMounts:
        - name: config-volume
          mountPath: /etc/config
      volumes:
      - name: config-volume
        configMap:
          name: global
          items:
          - key: settings
            path: component1
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component2/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: component2
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
        env:
        - name: APP_USER
          value: $(ConfigMap.global.data.user)
        command: [ "/bin/sh", "-c", "cat /etc/config/component2 && sleep 60" ]
        volumeMounts:
        - name: config-volume
          mountPath: /etc/config
      volumes:
      - name: config-volume
        configMap:
          name: global
          items:
          - key: settings
            path: component2
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/myapp -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_component1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: component1
  name: component1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: component1
  template:
    metadata:
      labels:
        app: component1
    spec:
      containers:
      - command:
        - /bin/sh
        - -c
        - cat /etc/config/component1 && sleep 60
        env:
        - name: APP_USER
          value: myuser
        image: k8s.gcr.io/busybox
        name: component1
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      volumes:
      - configMap:
          items:
          - key: settings
            path: component1
          name: global
        name: config-volume
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_component2.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: component2
  name: component2
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
        env:
        - name: APP_USER
          value: myuser
        image: k8s.gcr.io/busybox
        name: component2
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      volumes:
      - configMap:
          items:
          - key: settings
            path: component2
          name: global
        name: config-volume
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_global.yaml
apiVersion: v1
data:
  settings: |
    database: mydb
    port: 3000
  user: myuser
kind: ConfigMap
metadata:
  name: global
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

