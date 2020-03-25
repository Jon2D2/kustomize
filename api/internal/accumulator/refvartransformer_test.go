// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package accumulator

import (
	"reflect"
	"testing"

	"sigs.k8s.io/kustomize/api/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	resmaptest_test "sigs.k8s.io/kustomize/api/testutils/resmaptest"
	"sigs.k8s.io/kustomize/api/types"
)

func TestRefVarTransformer(t *testing.T) {
	type given struct {
		varMap map[string]interface{}
		fs     []types.FieldSpec
		res    resmap.ResMap
	}
	type expected struct {
		res    resmap.ResMap
		unused []string
	}
	testCases := []struct {
		description string
		given       given
		expected    expected
	}{
		{
			description: "var replacement in map[string]",
			given: given{
				varMap: map[string]interface{}{
					"FOO": "replacementForFoo",
					"BAR": "replacementForBar",
					"BAZ": int64(5),
					"BOO": true,
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/map"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/slice"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/interface"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/nil"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/num"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"map": map[string]interface{}{
								"item1": "$(FOO)",
								"item2": "bla",
								"item3": "$(BAZ)",
								"item4": "$(BAZ)+$(BAZ)",
								"item5": "$(BOO)",
								"item6": "if $(BOO)",
								"item7": 2019,
							},
							"slice": []interface{}{
								"$(FOO)",
								"bla",
								"$(BAZ)",
								"$(BAZ)+$(BAZ)",
								"$(BOO)",
								"if $(BOO)",
							},
							"interface": "$(FOO)",
							"nil":       nil,
							"num":       2019,
						}}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"map": map[string]interface{}{
								"item1": "replacementForFoo",
								"item2": "bla",
								"item3": int64(5),
								"item4": "5+5",
								"item5": true,
								"item6": "if true",
								"item7": 2019,
							},
							"slice": []interface{}{
								"replacementForFoo",
								"bla",
								int64(5),
								"5+5",
								true,
								"if true",
							},
							"interface": "replacementForFoo",
							"nil":       nil,
							"num":       2019,
						}}).ResMap(),
				unused: []string{"BAR"},
			},
		},
		{
			description: "inlining",
			given: given{
				varMap: map[string]interface{}{
					"FOO": map[string]interface{}{
						"foofield1": "foovalue1",
						"foofield2": "foovalue2",
					},
					"BAR": "replacementForBar",
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/item1"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"item1": "$(FOO)",
							"item2": "bla",
						},
					}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"item1": map[string]interface{}{
								"foofield1": "foovalue1",
								"foofield2": "foovalue2",
							},
							"item2": "bla",
						}}).ResMap(),
				unused: []string{"BAR"},
			},
		},
		{
			description: "parent-inlining",
			given: given{
				varMap: map[string]interface{}{
					"FOO": map[string]interface{}{
						"foofield1": "foovalue1",
						"foofield2": "foovalue2",
						"foofield3": "foovalue3",
						"foofield4": "foovalue4",
					},
					"BAR": "replacementForBar",
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"parent-inline": "$(FOO)",
							"foofield3":     "bla",
						},
					}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"foofield1": "foovalue1",
							"foofield2": "foovalue2",
							"foofield3": "bla",
							"foofield4": "foovalue4",
						}}).ResMap(),
				unused: []string{"BAR"},
			},
		},
		{
			description: "parent-inlining-and-variable-replacement",
			given: given{
				varMap: map[string]interface{}{
					"FOO": map[string]interface{}{
						"foofield1": "foovalue1",
						"foofield2": "foovalue2",
						"foofield3": "foovalue3",
						"foofield4": "foovalue4",
					},
					"BAR": "replacementForBar",
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"parent-inline": "$(FOO)",
							"foofield3":     "$(BAR)",
						},
					}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"foofield1": "foovalue1",
							"foofield2": "foovalue2",
							"foofield3": "replacementForBar",
							"foofield4": "foovalue4",
						}}).ResMap(),
				unused: []string{""},
			},
		},
		{
			description: "deeply-nested-parent-inlining",
			given: given{
				varMap: map[string]interface{}{
					"FOO": map[string]interface{}{
						"foofield1": "foovalue1",
						"foofield2": "foovalue2",
						"foofield3": "foovalue3",
						"foofield4": "foovalue4",
						"foofield5": map[string]interface{}{
							"nestedfield": "nestedvalue",
						},
					},
					"BAR": "replacementForBar",
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"parent-inline": "$(FOO)",
							"foofield3":     "$(BAR)",
						},
					}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"foofield1": "foovalue1",
							"foofield2": "foovalue2",
							"foofield3": "replacementForBar",
							"foofield4": "foovalue4",
							"foofield5": map[string]interface{}{
								"nestedfield": "nestedvalue",
							},
						}}).ResMap(),
				unused: []string{""},
			},
		},
		{
			description: "deeply-nested-parent-inlining-and-variable-replacement",
			given: given{
				varMap: map[string]interface{}{
					"FOO": map[string]interface{}{
						"foofield1": "foovalue1",
						"foofield2": "foovalue2",
						"foofield3": "foovalue3",
						"foofield4": "foovalue4",
						"foofield5": map[string]interface{}{
							"nestedfield1": "nestedvalue1",
							"nestedfield2": "nestedvalue2",
							"nestedfield3": "nestedvalue3",
						},
					},
					"BAR": "replacementForBar",
					"BAZ": map[string]interface{}{
						"nestedfield1": "nestedvalue1",
						"nestedfield2": "nestedvalue2",
						"nestedfield3": "nestedvalue3",
					},
					"QUX": "replacementForQux",
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/foofield5"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"parent-inline": "$(FOO)",
							"foofield3":     "$(BAR)",
							"foofield5": map[string]interface{}{
								"parent-inline": "$(BAZ)",
								"nestedfield2":  "$(QUX)",
							},
						},
					}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"foofield1": "foovalue1",
							"foofield2": "foovalue2",
							"foofield3": "replacementForBar",
							"foofield4": "foovalue4",
							"foofield5": map[string]interface{}{
								"nestedfield1": "nestedvalue1",
								"nestedfield2": "replacementForQux",
								"nestedfield3": "nestedvalue3",
							},
						}}).ResMap(),
				unused: []string{""},
			},
		},
		{
			description: "recursive-inlining",
			given: given{
				varMap: map[string]interface{}{
					"FOO": map[string]interface{}{
						"level2": "$(BAR)",
					},
					"BAR": map[string]interface{}{
						"level3": "$(QUX)",
					},
					"QUX": map[string]interface{}{
						"level4": "$(ZOO)",
					},
					"ZOO": "replacementForZoo",
				},
				fs: []types.FieldSpec{
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/level1"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/level1/level2"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/level1/level2/level3"},
					{Gvk: resid.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/level1/level2/level3/level4"},
				},
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"level1": "$(FOO)",
						},
					}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(
					t, resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"level1": map[string]interface{}{
								"level2": map[string]interface{}{
									"level3": map[string]interface{}{
										"level4": "replacementForZoo",
									},
								},
							},
						}}).ResMap(),
				unused: []string{""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// arrange
			tr := newRefVarTransformer(NewVarMap(tc.given.varMap), tc.given.fs)

			// act
			err := tr.Transform(tc.given.res)

			// assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			a, e := tc.given.res, tc.expected.res
			if !reflect.DeepEqual(a, e) {
				err = e.ErrorIfNotEqualLists(a)
				t.Fatalf("actual doesn't match expected: \nACTUAL:\n%v\nEXPECTED:\n%v\nERR: %v", a, e, err)
			}

		})
	}
}
