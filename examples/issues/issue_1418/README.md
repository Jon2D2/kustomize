# Feature Test for Issue 1418


This folder contains files describing how to address [Issue 1418](https://github.com/kubernetes-sigs/kustomize/issues/1418)

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
mkdir -p ${DEMO_HOME}/overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- resources.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/resources.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: filebeat
  namespace: kube-system
  labels:
    k8s-app: filebeat
spec:
  template:
    metadata:
      labels:
        k8s-app: filebeat
    spec:
      serviceAccountName: filebeat
      terminationGracePeriodSeconds: 30
      containers:
      - name: filebeat
        image: docker.elastic.co/beats/filebeat:7.2.0
        args: [
          "-c", "/etc/filebeat.yml",
          "-e",
          "--strict.perms=false",
        ]
        securityContext:
          runAsUser: 0
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 100Mi
        volumeMounts:
        - name: config
          mountPath: /etc/filebeat.yml
          readOnly: true
          subPath: filebeat.yml
        - name: data
          mountPath: /usr/share/filebeat/data
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      volumes:
      - name: config
        configMap:
          defaultMode: "0600"
          name: filebeat-config
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
      # data folder stores a registry of read status for all files, so we don't send everything again on a Filebeat pod restart
      - name: data
        hostPath:
          path: /var/lib/filebeat-data
          type: DirectoryOrCreate
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_daemonset_filebeat.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: filebeat
  name: filebeat
  namespace: kube-system
spec:
  template:
    metadata:
      labels:
        k8s-app: filebeat
    spec:
      containers:
      - args:
        - -c
        - /etc/filebeat.yml
        - -e
        - --strict.perms=false
        image: docker.elastic.co/beats/filebeat:7.2.0
        name: filebeat
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 100Mi
        securityContext:
          runAsUser: 0
        volumeMounts:
        - mountPath: /etc/filebeat.yml
          name: config
          readOnly: true
          subPath: filebeat.yml
        - mountPath: /usr/share/filebeat/data
          name: data
        - mountPath: /var/lib/docker/containers
          name: varlibdockercontainers
          readOnly: true
      serviceAccountName: filebeat
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: "0600"
          name: filebeat-config
        name: config
      - hostPath:
          path: /var/lib/docker/containers
        name: varlibdockercontainers
      - hostPath:
          path: /var/lib/filebeat-data
          type: DirectoryOrCreate
        name: data
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

