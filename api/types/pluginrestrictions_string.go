// Code generated by "stringer -type=PluginRestrictions"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PluginRestrictionsUnknown-0]
	_ = x[PluginRestrictionsBuiltinsOnly-1]
	_ = x[PluginRestrictionsNone-2]
}

const _PluginRestrictionsName = "PluginRestrictionsUnknownPluginRestrictionsBuiltinsOnlyPluginRestrictionsNone"

var _PluginRestrictionsIndex = [...]uint8{0, 25, 55, 77}

func (i PluginRestrictions) String() string {
	if i < 0 || i >= PluginRestrictions(len(_PluginRestrictionsIndex)-1) {
		return "PluginRestrictions(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PluginRestrictionsName[_PluginRestrictionsIndex[i]:_PluginRestrictionsIndex[i+1]]
}
