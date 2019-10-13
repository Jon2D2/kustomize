# Feature Test for Issue 1372


This folder contains files describing how to address [Issue 1372](https://github.com/kubernetes-sigs/kustomize/issues/1372)

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
mkdir -p ${DEMO_HOME}/environment
mkdir -p ${DEMO_HOME}/environment/base
mkdir -p ${DEMO_HOME}/environment/dev
mkdir -p ${DEMO_HOME}/environment/prd
mkdir -p ${DEMO_HOME}/permutations
mkdir -p ${DEMO_HOME}/permutations/scenario1-dev-fallback-v1
mkdir -p ${DEMO_HOME}/permutations/scenario1-dev-normal-v1
mkdir -p ${DEMO_HOME}/permutations/scenario1-prd-fallback-v1
mkdir -p ${DEMO_HOME}/permutations/scenario1-prd-normal-v1
mkdir -p ${DEMO_HOME}/permutations/scenario2-dev-fallback-v1
mkdir -p ${DEMO_HOME}/permutations/scenario2-dev-normal-v1
mkdir -p ${DEMO_HOME}/permutations/scenario2-prd-fallback-v1
mkdir -p ${DEMO_HOME}/permutations/scenario2-prd-normal-v1
mkdir -p ${DEMO_HOME}/processor
mkdir -p ${DEMO_HOME}/processor/base
mkdir -p ${DEMO_HOME}/processor/fallback
mkdir -p ${DEMO_HOME}/processor/normal
mkdir -p ${DEMO_HOME}/scenario
mkdir -p ${DEMO_HOME}/scenario/base
mkdir -p ${DEMO_HOME}/scenario/scenario1
mkdir -p ${DEMO_HOME}/scenario/scenario2
mkdir -p ${DEMO_HOME}/version
mkdir -p ${DEMO_HOME}/version/base
mkdir -p ${DEMO_HOME}/version/v1
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../common

# placeholder for environment common patches
patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/dev/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/prd/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario1-dev-fallback-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario1-dev-fallback-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario1 # apply changes necessary for scenario1
- ../../environment/dev    # apply changes for dev
- ../../processor/fallback # apply changes for fallback
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario1-dev-normal-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario1-dev-normal-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario1 # apply changes necessary for scenario1
- ../../environment/dev    # apply changes for dev
- ../../processor/normal   # apply changes for normal
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario1-prd-fallback-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario1-prd-fallback-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario1 # apply changes necessary for scenario1
- ../../environment/prd    # apply changes for prd
- ../../processor/fallback # apply changes for fallback
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario1-prd-normal-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario1-prd-normal-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario1 # apply changes necessary for scenario1
- ../../environment/prd    # apply changes for prd
- ../../processor/normal   # apply changes for normal
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario2-dev-fallback-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario2-dev-fallback-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario2 # apply changes necessary for scenario2
- ../../environment/dev    # apply changes for dev
- ../../processor/fallback # apply changes for fallback
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario2-dev-normal-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario2-dev-normal-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario2 # apply changes necessary for scenario2
- ../../environment/dev    # apply changes for dev
- ../../processor/normal   # apply changes for normal
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile10

<!-- @createKustomizationFile10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario2-prd-fallback-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario2-prd-fallback-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario2 # apply changes necessary for scenario2
- ../../environment/prd    # apply changes for prd
- ../../processor/fallback # apply changes for fallback
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile11

<!-- @createKustomizationFile11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/permutations/scenario2-prd-normal-v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: scenario2-prd-normal-v1-
resources:
- ../../common             # Has Deployment
- ../../scenario/scenario2 # apply changes necessary for scenario2
- ../../environment/prd    # apply changes for prd
- ../../processor/normal   # apply changes for normal
- ../../version/v1         # Apply v1
EOF
```


### Preparation Step KustomizationFile12

<!-- @createKustomizationFile12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/processor/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../common

commonLabels:
   my-label: my-app

# placeholder for processor common patches
patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile13

<!-- @createKustomizationFile13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/processor/fallback/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step KustomizationFile14

<!-- @createKustomizationFile14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/processor/normal/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step KustomizationFile15

<!-- @createKustomizationFile15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/scenario/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../common

# placeholder for scenario common patches
patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile16

<!-- @createKustomizationFile16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/scenario/scenario1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step KustomizationFile17

<!-- @createKustomizationFile17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/scenario/scenario2/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step KustomizationFile18

<!-- @createKustomizationFile18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/version/base/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../common

# placeholder for version common patches
patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile19

<!-- @createKustomizationFile19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/version/v1/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../base

patchesStrategicMerge:
- patch.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/dev/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 1
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/prd/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 3
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/processor/fallback/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9090
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/processor/normal/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/scenario/scenario1/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/scenario/scenario2/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      dnsPolicy: Default
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/version/v1/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image:v1
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/permutations/scenario1-dev-fallback-v1 -o ${DEMO_HOME}/actual/scenario1-dev-fallback-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario1-dev-normal-v1 -o ${DEMO_HOME}/actual/scenario1-dev-normal-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario1-prd-fallback-v1 -o ${DEMO_HOME}/actual/scenario1-prd-fallback-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario1-prd-normal-v1 -o ${DEMO_HOME}/actual/scenario1-prd-normal-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario2-dev-fallback-v1 -o ${DEMO_HOME}/actual/scenario2-dev-fallback-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario2-dev-normal-v1 -o ${DEMO_HOME}/actual/scenario2-dev-normal-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario2-prd-fallback-v1 -o ${DEMO_HOME}/actual/scenario2-prd-fallback-v1.yaml
kustomize build ${DEMO_HOME}/permutations/scenario2-prd-normal-v1 -o ${DEMO_HOME}/actual/scenario2-prd-normal-v1.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario1-dev-fallback-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario1-dev-fallback-v1-my-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9090
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario1-dev-normal-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario1-dev-normal-v1-my-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario1-prd-fallback-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario1-prd-fallback-v1-my-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9090
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario1-prd-normal-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario1-prd-normal-v1-my-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario2-dev-fallback-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario2-dev-fallback-v1-my-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9090
        name: my-deployment
      dnsPolicy: Default
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario2-dev-normal-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario2-dev-normal-v1-my-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: Default
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario2-prd-fallback-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario2-prd-fallback-v1-my-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9090
        name: my-deployment
      dnsPolicy: Default
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/scenario2-prd-normal-v1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    my-label: my-app
  name: scenario2-prd-normal-v1-my-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      my-label: my-app
  template:
    metadata:
      labels:
        my-label: my-app
    spec:
      containers:
      - image: my-image:v1
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: Default
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

