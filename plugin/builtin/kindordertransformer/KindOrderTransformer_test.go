// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestKubectlApplyOrderTransformer(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		PrepBuiltin("KindOrderTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: KindOrderTransformer
metadata:
  name: kubectl-apply
builtinordername: kubectlapply
`, `
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Namespace
metadata:
  name: apple
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: Namespace
metadata:
  name: apple
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
`)
}

func TestKubectlApplyOrderTransformerCustomOrder(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		PrepBuiltin("KindOrderTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: KindOrderTransformer
metadata:
  name: kubectl-apply
kindorder:
- Namespace
- LimitRange
- CRD1
- CRD2
- Secret
- ConfigMap
- Role
- Service
- Deployment
- Ingress
- ValidatingWebhookConfiguration
`, `
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Namespace
metadata:
  name: apple
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: Namespace
metadata:
  name: apple
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
`)
}

func TestKubectlDeleteOrderTransformer(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		PrepBuiltin("KindOrderTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: KindOrderTransformer
metadata:
  name: kubectl-delete
builtinordername: kubectldelete
`, `
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Namespace
metadata:
  name: apple
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Namespace
metadata:
  name: apple
`)
}

func TestKubectlDeleteOrderTransformerCustomOrder(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		PrepBuiltin("KindOrderTransformer")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: KindOrderTransformer
metadata:
  name: kubectl-delete
kindorder:
- ValidatingWebhookConfiguration
- Ingress
- Deployment
- Service
- Role
- ConfigMap
- Secret
- CRD2
- CRD1
- LimitRange
- Namespace
`, `
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Namespace
metadata:
  name: apple
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: ValidatingWebhookConfiguration
metadata:
  name: pomegranate
---
apiVersion: v1
kind: Ingress
metadata:
  name: durian
---
apiVersion: v1
kind: Deployment
metadata:
  name: pear
---
apiVersion: v1
kind: Service
metadata:
  name: papaya
---
apiVersion: v1
kind: Role
metadata:
  name: banana
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: apricot
---
apiVersion: v1
kind: Secret
metadata:
  name: quince
---
apiVersion: v1
kind: CRD2
metadata:
  name: durian
---
apiVersion: v1
kind: CRD1
metadata:
  name: durian
---
apiVersion: v1
kind: LimitRange
metadata:
  name: peach
---
apiVersion: v1
kind: Namespace
metadata:
  name: apple
`)
}
