package strings

import baseStrings "strings"

// removes "Data" from the message name and returns it all lower case
// SelectorRowData --> selectorrow
func NormalizeWidgetDataMessageName(s string) string {
	lower := baseStrings.ToLower(s)
	return baseStrings.TrimSuffix(lower, "data")
}

// removes "_" from the widget type and returns it all lower case
// SELECTOR_ROW --> selectorrow
func NormalizeWidgetType(s string) string {
	lower := baseStrings.ToLower(s)
	return baseStrings.ReplaceAll(lower, "_", "")
}
