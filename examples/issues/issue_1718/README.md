# Feature Test for Issue 1718


This folder contains files describing how to address [Issue 1718](https://github.com/kubernetes-sigs/kustomize/issues/1718)

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
## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build github.com/keleustes/kustomize/examples/springboot/overlays/staging?ref=allinone -o ${DEMO_HOME}/actual/staging.yaml
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
  application.properties: |
    app.name=Staging Kinflate Demo
    spring.jpa.hibernate.ddl-auto=update
    spring.datasource.url=jdbc:mysql://<staging_db_ip>:3306/db_example
    spring.datasource.username=root
    spring.datasource.password=admin
  foo: bar
  staging: ""
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: staging-demo-configmap-6m4cgm6h26
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: sbdemo
  name: staging-sbdemo
spec:
  ports:
  - port: 8080
  selector:
    app: sbdemo
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: sbdemo
  name: staging-sbdemo
spec:
  selector:
    matchLabels:
      app: sbdemo
  template:
    metadata:
      labels:
        app: sbdemo
    spec:
      containers:
      - image: jingfang/sbdemo
        name: sbdemo
        ports:
        - containerPort: 8080
        volumeMounts:
        - mountPath: /config
          name: demo-config
      volumes:
      - configMap:
          name: staging-demo-configmap-6m4cgm6h26
        name: demo-config
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

