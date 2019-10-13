# Mixin with Kustomize using Patches and Vars.

Demonstrate how to leverage kustomize transformers, patches and vars
to achieve mixin.

##  Layout

.
├── app
│   ├── base
│   │   ├── kustomization.yaml
│   │   └── values.yaml
│   ├── dev
│   │   ├── base
│   │   │   ├── kustomization.yaml
│   │   │   ├── overlayvalues.yaml
│   │   │   └── values.yaml
│   │   └── clienta
│   │       ├── kustomization.yaml
│   │       └── overlayvalues.yaml
│   └── prod
│       ├── base
│       │   ├── kustomization.yaml
│       │   ├── overlayvalues.yaml
│       │   └── values.yaml
│       ├── clienta
│       │   ├── kustomization.yaml
│       │   └── overlayvalues.yaml
│       └── clientb
│           ├── kustomization.yaml
│           └── overlayvalues.yaml
└── components
    ├── appdeployment
    │   ├── kustomization.yaml
    │   ├── service.yaml
    │   └── values.yaml
    ├── mysql
    │   ├── kustomization.yaml
    │   ├── service.yaml
    │   └── values.yaml
    ├── persistencelayer
    │   ├── kustomization.yaml
    │   ├── service.yaml
    │   └── values.yaml
    ├── rediscache
    │   ├── kustomization.yaml
    │   ├── service.yaml
    │   └── values.yaml
    ├── redissession
    │   ├── kustomization.yaml
    │   ├── service.yaml
    │   └── values.yaml
    └── varnish
        ├── kustomization.yaml
        ├── service.yaml
        └── values.yaml

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/app
mkdir -p ${DEMO_HOME}/app/base
mkdir -p ${DEMO_HOME}/app/dev
mkdir -p ${DEMO_HOME}/app/dev/base
mkdir -p ${DEMO_HOME}/app/dev/clienta
mkdir -p ${DEMO_HOME}/app/prod
mkdir -p ${DEMO_HOME}/app/prod/base
mkdir -p ${DEMO_HOME}/app/prod/clienta
mkdir -p ${DEMO_HOME}/app/prod/clientb
mkdir -p ${DEMO_HOME}/components
mkdir -p ${DEMO_HOME}/components/appdeployment
mkdir -p ${DEMO_HOME}/components/mysql
mkdir -p ${DEMO_HOME}/components/persistencelayer
mkdir -p ${DEMO_HOME}/components/rediscache
mkdir -p ${DEMO_HOME}/components/redissession
mkdir -p ${DEMO_HOME}/components/varnish
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../components/appdeployment
- ../../components/mysql
- ../../components/persistencelayer
- ../../components/rediscache
- ../../components/redissession
- ../../components/varnish
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/dev/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../base
- values.yaml

patchesStrategicMerge:
- overlayvalues.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/dev/clienta/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- overlayvalues.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../base
- values.yaml

patchesStrategicMerge:
- overlayvalues.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/clienta/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- overlayvalues.yaml
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/clientb/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- overlayvalues.yaml
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/appdeployment/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./service.yaml
- ./values.yaml
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/mysql/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./service.yaml
- ./values.yaml
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/persistencelayer/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./service.yaml
- ./values.yaml
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/rediscache/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./service.yaml
- ./values.yaml
EOF
```


### Preparation Step KustomizationFile10

<!-- @createKustomizationFile10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/redissession/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./service.yaml
- ./values.yaml
EOF
```


### Preparation Step KustomizationFile11

<!-- @createKustomizationFile11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/varnish/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./service.yaml
- ./values.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/base/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: app-base
  namespace: build
spec:
  field1: value1
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/dev/base/overlayvalues.yaml
apiVersion: v1
kind: Values
metadata:
  name: appdeployment
  namespace: build
spec:
  port: 8501
---
apiVersion: v1
kind: Values
metadata:
  name: mysql
  namespace: build
spec:
  targetPort: 9502
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/dev/base/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: app-dev-base
  namespace: build
spec:
  field1: value1
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/dev/clienta/overlayvalues.yaml
apiVersion: v1
kind: Values
metadata:
  name: mysql
  namespace: build
spec:
  targetPort: 12502
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/base/overlayvalues.yaml
apiVersion: v1
kind: Values
metadata:
  name: rediscache
  namespace: build
spec:
  port: 8704
---
apiVersion: v1
kind: Values
metadata:
  name: redissession
  namespace: build
spec:
  targetPort: 9705
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/base/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: app-prod-base
  namespace: build
spec:
  field1: value1
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/clienta/overlayvalues.yaml
apiVersion: v1
kind: Values
metadata:
  name: persistencelayer
  namespace: build
spec:
  port: 28903
  targetPort: 29903
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/prod/clientb/overlayvalues.yaml
apiVersion: v1
kind: Values
metadata:
  name: redissession
  namespace: build
spec:
  port: 18203
  targetPort: 19203
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/appdeployment/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: appdeployment
spec:
  selector:
    app: appdeployment
  ports:
  - name: web
    port: $(Values.appdeployment.spec.port)
    targetPort: $(Values.appdeployment.spec.targetPort)
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/appdeployment/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: appdeployment
  namespace: build
spec:
  port: 8001
  targetPort: 9001
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/mysql/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  selector:
    app: mysql
  ports:
  - name: web
    port: $(Values.mysql.spec.port)
    targetPort: $(Values.mysql.spec.targetPort)
EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/mysql/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: mysql
  namespace: build
spec:
  port: 8002
  targetPort: 9002
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/persistencelayer/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: persistencelayer
spec:
  selector:
    app: persistencelayer
  ports:
  - name: web
    port: $(Values.persistencelayer.spec.port)
    targetPort: $(Values.persistencelayer.spec.targetPort)
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/persistencelayer/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: persistencelayer
  namespace: build
spec:
  port: 8003
  targetPort: 9003
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/rediscache/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: rediscache
spec:
  selector:
    app: rediscache
  ports:
  - name: web
    port: $(Values.rediscache.spec.port)
    targetPort: $(Values.rediscache.spec.targetPort)
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/rediscache/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: rediscache
  namespace: build
spec:
  port: 8004
  targetPort: 9004
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/redissession/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: redissession
spec:
  selector:
    app: redissession
  ports:
  - name: web
    port: $(Values.redissession.spec.port)
    targetPort: $(Values.redissession.spec.targetPort)
EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/redissession/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: redissession
  namespace: build
spec:
  port: 8005
  targetPort: 9005
EOF
```


### Preparation Step Resource18

<!-- @createResource18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/varnish/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: varnish
spec:
  selector:
    app: varnish
  ports:
  - name: web
    port: $(Values.varnish.spec.port)
    targetPort: $(Values.varnish.spec.targetPort)
EOF
```


### Preparation Step Resource19

<!-- @createResource19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/components/varnish/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: varnish
  namespace: build
spec:
  port: 8006
  targetPort: 9006
EOF
```

## Execution

### app-dev-clienta

<!-- @buildDevClientA @test -->
```bash
mkdir -p ${DEMO_HOME}/actual/dev-clienta
kustomize build ${DEMO_HOME}/app/dev/clienta -o ${DEMO_HOME}/actual/dev-clienta
```

### app-prod-clienta

<!-- @buildProdClientA @test -->
```bash
mkdir -p ${DEMO_HOME}/actual/prod-clienta
kustomize build ${DEMO_HOME}/app/prod/clienta -o ${DEMO_HOME}/actual/prod-clienta
```

### app-prod-clientb

<!-- @buildProdClientB @test -->
```bash
mkdir -p ${DEMO_HOME}/actual/prod-clientb
kustomize build ${DEMO_HOME}/app/prod/clientb -o ${DEMO_HOME}/actual/prod-clientb
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/dev-clienta
mkdir -p ${DEMO_HOME}/expected/prod-clienta
mkdir -p ${DEMO_HOME}/expected/prod-clientb
```

### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev-clienta/default_~g_v1_service_appdeployment.yaml
apiVersion: v1
kind: Service
metadata:
  name: appdeployment
spec:
  ports:
  - name: web
    port: 8501
    targetPort: 9001
  selector:
    app: appdeployment
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev-clienta/default_~g_v1_service_mysql.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  ports:
  - name: web
    port: 8002
    targetPort: 12502
  selector:
    app: mysql
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev-clienta/default_~g_v1_service_persistencelayer.yaml
apiVersion: v1
kind: Service
metadata:
  name: persistencelayer
spec:
  ports:
  - name: web
    port: 8003
    targetPort: 9003
  selector:
    app: persistencelayer
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev-clienta/default_~g_v1_service_rediscache.yaml
apiVersion: v1
kind: Service
metadata:
  name: rediscache
spec:
  ports:
  - name: web
    port: 8004
    targetPort: 9004
  selector:
    app: rediscache
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev-clienta/default_~g_v1_service_redissession.yaml
apiVersion: v1
kind: Service
metadata:
  name: redissession
spec:
  ports:
  - name: web
    port: 8005
    targetPort: 9005
  selector:
    app: redissession
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev-clienta/default_~g_v1_service_varnish.yaml
apiVersion: v1
kind: Service
metadata:
  name: varnish
spec:
  ports:
  - name: web
    port: 8006
    targetPort: 9006
  selector:
    app: varnish
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clienta/default_~g_v1_service_appdeployment.yaml
apiVersion: v1
kind: Service
metadata:
  name: appdeployment
spec:
  ports:
  - name: web
    port: 8001
    targetPort: 9001
  selector:
    app: appdeployment
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clienta/default_~g_v1_service_mysql.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  ports:
  - name: web
    port: 8002
    targetPort: 9002
  selector:
    app: mysql
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clienta/default_~g_v1_service_persistencelayer.yaml
apiVersion: v1
kind: Service
metadata:
  name: persistencelayer
spec:
  ports:
  - name: web
    port: 28903
    targetPort: 29903
  selector:
    app: persistencelayer
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clienta/default_~g_v1_service_rediscache.yaml
apiVersion: v1
kind: Service
metadata:
  name: rediscache
spec:
  ports:
  - name: web
    port: 8704
    targetPort: 9004
  selector:
    app: rediscache
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clienta/default_~g_v1_service_redissession.yaml
apiVersion: v1
kind: Service
metadata:
  name: redissession
spec:
  ports:
  - name: web
    port: 8005
    targetPort: 9705
  selector:
    app: redissession
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clienta/default_~g_v1_service_varnish.yaml
apiVersion: v1
kind: Service
metadata:
  name: varnish
spec:
  ports:
  - name: web
    port: 8006
    targetPort: 9006
  selector:
    app: varnish
EOF
```


### Verification Step Expected12

<!-- @createExpected12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clientb/default_~g_v1_service_appdeployment.yaml
apiVersion: v1
kind: Service
metadata:
  name: appdeployment
spec:
  ports:
  - name: web
    port: 8001
    targetPort: 9001
  selector:
    app: appdeployment
EOF
```


### Verification Step Expected13

<!-- @createExpected13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clientb/default_~g_v1_service_mysql.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  ports:
  - name: web
    port: 8002
    targetPort: 9002
  selector:
    app: mysql
EOF
```


### Verification Step Expected14

<!-- @createExpected14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clientb/default_~g_v1_service_persistencelayer.yaml
apiVersion: v1
kind: Service
metadata:
  name: persistencelayer
spec:
  ports:
  - name: web
    port: 8003
    targetPort: 9003
  selector:
    app: persistencelayer
EOF
```


### Verification Step Expected15

<!-- @createExpected15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clientb/default_~g_v1_service_rediscache.yaml
apiVersion: v1
kind: Service
metadata:
  name: rediscache
spec:
  ports:
  - name: web
    port: 8704
    targetPort: 9004
  selector:
    app: rediscache
EOF
```


### Verification Step Expected16

<!-- @createExpected16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clientb/default_~g_v1_service_redissession.yaml
apiVersion: v1
kind: Service
metadata:
  name: redissession
spec:
  ports:
  - name: web
    port: 18203
    targetPort: 19203
  selector:
    app: redissession
EOF
```


### Verification Step Expected17

<!-- @createExpected17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod-clientb/default_~g_v1_service_varnish.yaml
apiVersion: v1
kind: Service
metadata:
  name: varnish
spec:
  ports:
  - name: web
    port: 8006
    targetPort: 9006
  selector:
    app: varnish
EOF
```


<!-- @compareActualToExpectedDevClientA @test -->
```bash
rm -f $DEMO_HOME/actual/dev-clienta/build_*.yaml
test 0 == \
$(diff -r $DEMO_HOME/actual/dev-clienta $DEMO_HOME/expected/dev-clienta | wc -l); \
echo $?
```

<!-- @compareActualToExpectedProdClientA @test -->
```bash
rm -f $DEMO_HOME/actual/prod-clienta/build_*.yaml
test 0 == \
$(diff -r $DEMO_HOME/actual/prod-clienta $DEMO_HOME/expected/prod-clienta | wc -l); \
echo $?
```

<!-- @compareActualToExpectedProdClientB @test -->
```bash
rm -f $DEMO_HOME/actual/prod-clientb/build_*.yaml
test 0 == \
$(diff -r $DEMO_HOME/actual/prod-clientb $DEMO_HOME/expected/prod-clientb | wc -l); \
echo $?
```





