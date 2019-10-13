# Feature Test for Issue update-env-variable


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
- ./deployment.yaml

patchesStrategicMerge:
- ./patch.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: strimzi-topic-operator
spec:
  template:
    spec:
      containers:
      - name: container1
        env:
        - name: STRIMZI_RESOURCE_LABELS
          value: "strimzi.io/cluster=my-cluster"
        - name: STRIMZI_KAFKA_BOOTSTRAP_SERVERS
          value: my-cluster-kafka-bootstrap:9092
        - name: STRIMZI_ZOOKEEPER_CONNECT
          value: my-cluster-zookeeper-client:2181
      - name: container2
        env:
        - name: STRIMZI_RESOURCE_LABELS
          value: "strimzi.io/cluster=my-cluster"
        - name: STRIMZI_KAFKA_BOOTSTRAP_SERVERS
          value: my-cluster-kafka-bootstrap:9092
        - name: STRIMZI_ZOOKEEPER_CONNECT
          value: my-cluster-zookeeper-client:2181

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: strimzi-topic-operator
spec:
  template:
    spec:
      containers:
      - name: container2
        env:
        - name: STRIMZI_KAFKA_BOOTSTRAP_SERVERS
          value: updated-kafka-bootstrap:8888

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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_strimzi-topic-operator.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: strimzi-topic-operator
spec:
  template:
    spec:
      containers:
      - env:
        - name: STRIMZI_RESOURCE_LABELS
          value: strimzi.io/cluster=my-cluster
        - name: STRIMZI_KAFKA_BOOTSTRAP_SERVERS
          value: my-cluster-kafka-bootstrap:9092
        - name: STRIMZI_ZOOKEEPER_CONNECT
          value: my-cluster-zookeeper-client:2181
        name: container1
      - env:
        - name: STRIMZI_RESOURCE_LABELS
          value: strimzi.io/cluster=my-cluster
        - name: STRIMZI_KAFKA_BOOTSTRAP_SERVERS
          value: updated-kafka-bootstrap:8888
        - name: STRIMZI_ZOOKEEPER_CONNECT
          value: my-cluster-zookeeper-client:2181
        name: container2
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

