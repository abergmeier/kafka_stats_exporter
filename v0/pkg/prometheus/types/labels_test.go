package types

import "testing"

func TestLabelNames(t *testing.T) {
	lns := LabelNames{}
	err := lns.AddStrings("doo_dab", "u2")
	if err != nil {
		t.Fatal("Unexpected error from AddString:", err)
	}
	err = lns.AddStrings("doo_dab", "1_2")
	if err != nil {
		return
	}
	err = lns.AddStrings("foo.bar", "too_doo")
	if err != nil {
		return
	}

	t.Fatal("Expected returned error from AddStrings")
}
