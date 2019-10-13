# Feature Test for Issue 1243


This folder contains files describing how to address [Issue 1243](https://github.com/kubernetes-sigs/kustomize/issues/1243)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base-basens
mkdir -p ${DEMO_HOME}/base-default
mkdir -p ${DEMO_HOME}/base-nons
mkdir -p ${DEMO_HOME}/overlay-basens-basens
mkdir -p ${DEMO_HOME}/overlay-basens-stagingns
mkdir -p ${DEMO_HOME}/overlay-default-nons
mkdir -p ${DEMO_HOME}/overlay-default-stagingns
mkdir -p ${DEMO_HOME}/overlay-nons-default
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base-basens/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base-default/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base-nons/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-basens-basens/kustomization.yaml
resources:
- ../base-basens
- deployment.yaml

patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-basens-stagingns/kustomization.yaml
resources:
- ../base-basens
- deployment.yaml

patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-default-nons/kustomization.yaml
resources:
- ../base-default
- deployment.yaml

patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-default-stagingns/kustomization.yaml
resources:
- ../base-default
- deployment.yaml

patchesStrategicMerge:
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-nons-default/kustomization.yaml
resources:
- ../base-nons
- deployment.yaml

patchesStrategicMerge:
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base-basens/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  labels:
    app: dply1
  namespace: base
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base-default/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  labels:
    app: dply1
  namespace: default
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base-nons/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  labels:
    app: dply1
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-basens-basens/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
  namespace: base
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-basens-stagingns/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  namespace: staging
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-default-nons/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-default-stagingns/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  namespace: staging
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay-nons-default/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
  namespace: default
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```

## Execution

<!-- @createActualDir @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/basens-basens
mkdir -p ${DEMO_HOME}/actual/basens-stagingns
mkdir -p ${DEMO_HOME}/actual/default-nons
mkdir -p ${DEMO_HOME}/actual/default-stagingns
mkdir -p ${DEMO_HOME}/actual/nons-default
```

### Case 1: Base NS set to "base" , Overlay NS set to "staging"

Let's build into a file

<!-- @buildAsFileStaging @basensStagingNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-basens-stagingns -o $DEMO_HOME/actual/basens-stagingns.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @basensStagingNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/basens-stagingns
kustomize build $DEMO_HOME/overlay-basens-stagingns -o $DEMO_HOME/actual/basens-stagingns/
```

### Case 2: Base NS set to "base" , Overlay NS set to "base"

Let's build into a file

<!-- @buildAsFileStaging @baseNsBaseNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-basens-basens -o $DEMO_HOME/actual/basens-basens.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @baseNsBaseNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/basens-basens
kustomize build $DEMO_HOME/overlay-basens-basens -o $DEMO_HOME/actual/basens-basens/
```

### Case 3: Base NS set to "default" , Overlay NS not set

Let's build into a file

<!-- @buildAsFileStaging @defaultNoNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-default-nons -o $DEMO_HOME/actual/default-nons.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @defaultNoNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/default-nons
kustomize build $DEMO_HOME/overlay-default-nons -o $DEMO_HOME/actual/default-nons/
```

## Case 4: Base NS set to "default" , Overlay NS set to staging

Let's build into a file

<!-- @buildAsFileStaging @defaultStagingNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-default-stagingns -o $DEMO_HOME/actual/default-stagingns.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @defaultStagingNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/default-stagingns
kustomize build $DEMO_HOME/overlay-default-stagingns -o $DEMO_HOME/actual/default-stagingns/
```

### Case 5: Base NS not set , Overlay NS set to default

Let's build into a file

<!-- @buildAsFileStaging @noNsDefault @test -->
```sh
kustomize build $DEMO_HOME/overlay-nons-default -o $DEMO_HOME/actual/nons-default.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @noNsDefault @test -->
```sh
mkdir -p $DEMO_HOME/actual/nons-default
kustomize build $DEMO_HOME/overlay-nons-default -o $DEMO_HOME/actual/nons-default/
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/basens-basens
mkdir -p ${DEMO_HOME}/expected/basens-stagingns
mkdir -p ${DEMO_HOME}/expected/default-nons
mkdir -p ${DEMO_HOME}/expected/default-stagingns
mkdir -p ${DEMO_HOME}/expected/nons-default
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/basens-basens/apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: base
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/basens-basens/apps_v1beta2_deployment_dply2.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
  namespace: base
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/basens-basens.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: base
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
  namespace: base
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/basens-stagingns/base_apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: base
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/basens-stagingns/staging_apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  namespace: staging
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/basens-stagingns.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: base
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  namespace: staging
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default-nons/apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: default
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default-nons/apps_v1beta2_deployment_dply2.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default-nons.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: default
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default-stagingns/default_apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: default
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default-stagingns/staging_apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  namespace: staging
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default-stagingns.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
  namespace: default
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply1
  namespace: staging
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected12

<!-- @createExpected12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/nons-default/apps_v1beta2_deployment_dply1.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


### Verification Step Expected13

<!-- @createExpected13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/nons-default/apps_v1beta2_deployment_dply2.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
  namespace: default
spec:
  template:
    metadata:
      labels:
        from: overlay
EOF
```


### Verification Step Expected14

<!-- @createExpected14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/nons-default.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dply2
  namespace: default
spec:
  template:
    metadata:
      labels:
        from: overlay
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: dply1
  name: dply1
spec:
  selector:
    matchLabels:
      app: dply1
  template:
    metadata:
      labels:
        app: dply1
    spec:
      containers:
      - image: alpine
        name: dply1
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

