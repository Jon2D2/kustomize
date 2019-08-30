// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package builtinconfig

import (
	"log"
	"sort"

	"github.com/pkg/errors"

	"sigs.k8s.io/kustomize/api/ifc"
	"sigs.k8s.io/kustomize/api/konfig/builtinpluginconsts"
	"sigs.k8s.io/kustomize/api/types"
)

// TransformerConfig holds the data needed to perform transformations.
type TransformerConfig struct {
	NamePrefix        types.FsSlice `json:"namePrefix,omitempty" yaml:"namePrefix,omitempty"`
	NameSuffix        types.FsSlice `json:"nameSuffix,omitempty" yaml:"nameSuffix,omitempty"`
	NameSpace         types.FsSlice `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	CommonLabels      types.FsSlice `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
	CommonAnnotations types.FsSlice `json:"commonAnnotations,omitempty" yaml:"commonAnnotations,omitempty"`
	NameReference     nbrSlice      `json:"nameReference,omitempty" yaml:"nameReference,omitempty"`
	VarReference      types.FsSlice `json:"varReference,omitempty" yaml:"varReference,omitempty"`
	Images            types.FsSlice `json:"images,omitempty" yaml:"images,omitempty"`
	Replicas          types.FsSlice `json:"replicas,omitempty" yaml:"replicas,omitempty"`
}

// MakeEmptyConfig returns an empty TransformerConfig object
func MakeEmptyConfig() *TransformerConfig {
	return &TransformerConfig{}
}

// MakeDefaultConfig returns a default TransformerConfig.
func MakeDefaultConfig() *TransformerConfig {
	c, err := makeTransformerConfigFromBytes(
		builtinpluginconsts.GetDefaultFieldSpecs())
	if err != nil {
		log.Fatalf("Unable to make default transformconfig: %v", err)
	}
	return c
}

// MakeTransformerConfig returns a merger of custom config,
// if any, with default config.
func MakeTransformerConfig(
	ldr ifc.Loader, paths []string) (*TransformerConfig, error) {

	result := &TransformerConfig{}
	var err error

	// merge and normalize Default Config
	defaultcfg := MakeDefaultConfig()
	result, err = result.Merge(defaultcfg)
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return defaultcfg, nil
	}

	// merge and normalize user provided configurations
	result, err = loadDefaultConfig(result, ldr, paths)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// sortFields provides determinism in logging, tests, etc.
func (t *TransformerConfig) sortFields() {
	sort.Sort(t.NamePrefix)
	sort.Sort(t.NameSuffix)
	sort.Sort(t.NameSpace)
	sort.Sort(t.CommonLabels)
	sort.Sort(t.CommonAnnotations)
	sort.Sort(t.NameReference)
	sort.Sort(t.VarReference)
	sort.Sort(t.Images)
	sort.Sort(t.Replicas)
}

// AddPrefixFieldSpec adds a types.FieldSpec to NamePrefix
func (t *TransformerConfig) AddPrefixFieldSpec(fs types.FieldSpec) (err error) {
	t.NamePrefix, err = t.NamePrefix.MergeOne(types.FieldSpecConfig{FieldSpec: fs, Behavior: "add"})
	return err
}

// AddLabelFieldSpec adds a types.FieldSpec to CommonLabels
func (t *TransformerConfig) AddLabelFieldSpec(fs types.FieldSpec) (err error) {
	t.CommonLabels, err = t.CommonLabels.MergeOne(types.FieldSpecConfig{FieldSpec: fs, Behavior: "add"})
	return err
}

// AddAnnotationFieldSpec adds a types.FieldSpec to CommonAnnotations
func (t *TransformerConfig) AddAnnotationFieldSpec(fs types.FieldSpec) (err error) {
	t.CommonAnnotations, err = t.CommonAnnotations.MergeOne(types.FieldSpecConfig{FieldSpec: fs, Behavior: "add"})
	return err
}

// AddNamereferenceFieldSpec adds a NameBackReferences to NameReference
func (t *TransformerConfig) AddNamereferenceFieldSpec(
	nbrs NameBackReferences) (err error) {
	t.NameReference, err = t.NameReference.mergeOne(nbrs)
	return err
}

// Merge merges two TransformerConfigs objects into
// a new TransformerConfig object
func (t *TransformerConfig) Merge(input *TransformerConfig) (
	merged *TransformerConfig, err error) {
	if input == nil {
		return t, nil
	}
	merged = &TransformerConfig{}
	merged.NamePrefix, err = t.NamePrefix.MergeAll(input.NamePrefix)
	if err != nil {
		return nil, errors.Wrap(err, "NamePrefix")
	}
	merged.NameSpace, err = t.NameSpace.MergeAll(input.NameSpace)
	if err != nil {
		return nil, errors.Wrap(err, "NameSpace")
	}
	merged.CommonAnnotations, err = t.CommonAnnotations.MergeAll(
		input.CommonAnnotations)
	if err != nil {
		return nil, errors.Wrap(err, "CommonAnnotations")
	}
	merged.CommonLabels, err = t.CommonLabels.MergeAll(input.CommonLabels)
	if err != nil {
		return nil, errors.Wrap(err, "CommonLabels")
	}
	merged.VarReference, err = t.VarReference.MergeAll(input.VarReference)
	if err != nil {
		return nil, errors.Wrap(err, "VarReference")
	}
	merged.NameReference, err = t.NameReference.mergeAll(input.NameReference)
	if err != nil {
		return nil, errors.Wrap(err, "NameReference")
	}
	merged.Images, err = t.Images.MergeAll(input.Images)
	if err != nil {
		return nil, errors.Wrap(err, "Images")
	}
	merged.Replicas, err = t.Replicas.MergeAll(input.Replicas)
	if err != nil {
		return nil, errors.Wrap(err, "Replicas`")
	}
	merged.sortFields()
	return merged, nil
}

func (t *TransformerConfig) NamePrefixFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.NamePrefix)
}
func (t *TransformerConfig) NameSuffixFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.NameSuffix)
}
func (t *TransformerConfig) NameSpaceFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.NameSpace)
}
func (t *TransformerConfig) CommonLabelsFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.CommonLabels)
}
func (t *TransformerConfig) CommonAnnotationsFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.CommonAnnotations)
}
func (t *TransformerConfig) VarReferenceFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.VarReference)
}
func (t *TransformerConfig) ImagesFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.Images)
}
func (t *TransformerConfig) ReplicasFieldSpecs() types.FieldSpecs {
	return types.NewFieldSpecs(t.Replicas)
}
