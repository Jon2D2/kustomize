# Feature Test for Issue 1782


This folder contains files describing how to address [Issue 1782](https://github.com/kubernetes-sigs/kustomize/issues/1782)

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
- job.yaml

vars:
- name: BACKEND_SERVICE_NAME
  objref:
    kind: Job
    name: nginx-job
    apiVersion: batch/v1
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: nginx-job
spec:
  template:
    spec:
      containers:
        - name: nginx
          env:
            - name: BACKEND_HOST
              value: $(BACKEND_SERVICE_NAME)
      initContainers:
        - name: init-nginx-wait-backend
          image: busybox:1.31
          command: ['sh', '-c', 'until nslookup $(BACKEND_SERVICE_NAME); do sleep 2; done;']
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
cat <<'EOF' >${DEMO_HOME}/expected/batch_v1_job_nginx-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: nginx-job
spec:
  template:
    spec:
      containers:
      - env:
        - name: BACKEND_HOST
          value: nginx-job
        name: nginx
      initContainers:
      - command:
        - sh
        - -c
        - until nslookup nginx-job; do sleep 2; done;
        image: busybox:1.31
        name: init-nginx-wait-backend
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

