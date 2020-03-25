// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/api/resid"
)

func TestPathSlice(t *testing.T) {
	type path struct {
		input  string
		parsed []string
	}
	paths := []path{
		{
			input:  "spec/metadata/annotations",
			parsed: []string{"spec", "metadata", "annotations"},
		},
		{
			input:  `metadata/annotations/nginx.ingress.kubernetes.io\/auth-secret`,
			parsed: []string{"metadata", "annotations", "nginx.ingress.kubernetes.io/auth-secret"},
		},
	}
	for _, p := range paths {
		fs := FieldSpec{Path: p.input}
		actual := fs.PathSlice()
		if !reflect.DeepEqual(actual, p.parsed) {
			t.Fatalf("expected %v, but got %v", p.parsed, actual)
		}
	}
}

var mergeTests = []struct {
	name     string
	original FsSlice
	incoming FsSlice
	err      error
	result   FsSlice
}{
	{
		"normal",
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "pear"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "home",
					Gvk:                resid.Gvk{Group: "beans"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		nil,
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "pear"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "home",
					Gvk:                resid.Gvk{Group: "beans"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
		},
	},
	{
		"ignore copy",
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "pear"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		nil,
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "pear"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
		},
	},
	{
		"error on conflict",
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "pear"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "whatever",
					Gvk:                resid.Gvk{Group: "apple"},
					CreateIfNotPresent: true,
				},
				Behavior: "add",
			},
		},
		fmt.Errorf("hey"),
		FsSlice{},
	},
	{
		"remove",
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field1",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field2",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path: "spec/field1",
					Gvk:  resid.Gvk{Kind: "MyCRD"},
				},
				Behavior: "remove",
			},
		},
		nil,
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field2",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
		},
	},
	{
		"remove2",
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "metadata/labels",
					CreateIfNotPresent: true,
				},
			},
			{
				FieldSpec: FieldSpec{
					Path:               "spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels",
					Gvk:                resid.Gvk{Kind: "Deployment", Group: "apps"},
					CreateIfNotPresent: false,
				},
			},
			{
				FieldSpec: FieldSpec{
					Path:               "spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels",
					Gvk:                resid.Gvk{Kind: "Deployment", Group: "apps"},
					CreateIfNotPresent: false,
				},
			},
		},
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels",
					Gvk:                resid.Gvk{Kind: "Deployment", Group: "apps"},
					CreateIfNotPresent: false,
				},
				Behavior: "remove",
			},
		},
		nil,
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "metadata/labels",
					CreateIfNotPresent: true,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels",
					Gvk:                resid.Gvk{Kind: "Deployment", Group: "apps"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
		},
	},
	{
		"replace",
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field1",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field2",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: false,
				},
				Behavior: "add",
			},
		},
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field2",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: true,
				},
				Behavior: "replace",
			},
		},
		nil,
		FsSlice{
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field1",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: false,
				},
				Behavior: "",
			},
			{
				FieldSpec: FieldSpec{
					Path:               "spec/field2",
					Gvk:                resid.Gvk{Kind: "MyCRD"},
					CreateIfNotPresent: true,
				},
				Behavior: "",
			},
		},
	},
}

func TestFsSlice_MergeAll(t *testing.T) {
	for _, item := range mergeTests {
		result := FsSlice{}
		var err error

		// Normalize and merge original FsSlice
		result, err = result.MergeAll(item.original)
		if err != nil {
			t.Fatalf("test %s: unexpected err %v", item.name, err)
		}

		// Normalize and merge incoming FsSlice
		result, err = result.MergeAll(item.incoming)
		if item.err == nil {
			if err != nil {
				t.Fatalf("test %s: unexpected err %v", item.name, err)
			}
			if !reflect.DeepEqual(item.result, result) {
				t.Fatalf("test %s: expected: %v\n but got: %v\n",
					item.name, item.result, result)
			}
		} else {
			if err == nil {
				t.Fatalf("test %s: expected err: %v", item.name, err)
			}
			if !strings.Contains(err.Error(), "conflicting fieldspecs") {
				t.Fatalf("test %s: unexpected err: %v", item.name, err)
			}
		}
	}
}
