package model

import "testing"

var isInRangeTests = []struct {
	version  RangeVersion
	value    string
	expected bool
}{
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "1.0.1", true},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "2.0.0", true},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "2", true},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "2.0", true},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "1.1", true},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "2.0.1", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "0.1.0", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "a.a.a", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "a", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "0", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "1", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "1.0", false},
	{RangeVersion{Min: "1.0.0", Max: "2.0.0"}, "2.1", false},
	{RangeVersion{Min: "1.0.0"}, "2.0.0", true},
	{RangeVersion{Min: "1.0.0"}, "1.0.1", true},
	{RangeVersion{Min: "1.0.0"}, "0.1.1", false},
	{RangeVersion{Min: "1.1.2"}, "1.1.1", false},
	{RangeVersion{Max: "1.0.2"}, "1.0.1", true},
	{RangeVersion{Max: "1.0.2"}, "1.0.3", false},
}

func TestRangeVersion_IsInRange(t *testing.T) {
	for _, tt := range isInRangeTests {
		actual := tt.version.IsInRange(tt.value)
		if actual != tt.expected {
			t.Errorf("RangeVersion(%s, %s).IsInRange(%s): expected %t, actual %t", tt.version.Min, tt.version.Max, tt.value, tt.expected, actual)
		}
	}
}
