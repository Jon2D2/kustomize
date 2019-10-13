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

# Proper behavior is only obtained if
# nameSuffix is not set in here
# but in the overlay
nameSuffix: -suffix

configMapGenerator:
- name: bar
  files: [file2]
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# Proper behavior is only obtained if
# nameSuffix is not set in the base
# layer but set here in the overlay
# nameSuffix: -suffix

resources:
- ../base
- podspec.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/podspec.yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod
spec:
  volumes:
    - name: config-bar
      configMap:
        # From a layering standpoint, overlay expect
        # base to have added the suffix. If you try
        # to add the suffix in the overlay, don't forget
        # to remove the suffix here.
        name: bar-suffix
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/file2
file2 content
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_bar-suffix-ttmh4ct745.yaml
apiVersion: v1
data:
  file2: |
    file2 content
kind: ConfigMap
metadata:
  name: bar-suffix-ttmh4ct745
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_pod_pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod
spec:
  volumes:
  - configMap:
      name: bar-suffix
    name: config-bar
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

