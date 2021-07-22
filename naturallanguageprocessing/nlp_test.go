package naturallanguageprocessing

import "testing"

func TestIntentMap(t *testing.T) {
	tests := map[string][]string{
		"empty the recent users list":      {"reset", "recent"},
		"put my name on the list":          {"add", ""},
		"put <@ABCDEFGHIJK> on the list":   {"add", "ABCDEFGHIJK"},
		"please include <@ABCDEFGHIJK>":    {"add", "ABCDEFGHIJK"},
		"take me off the list":             {"remove", ""},
		"take <@ABCDEFGHIJK> off the list": {"remove", "ABCDEFGHIJK"},
		"forget <@ABCDEFGHIJK>":            {"remove", "ABCDEFGHIJK"},
		"gimme the list of names":          {"give", ""},
		"display all the names":            {"give", ""},
		"show me who's on the list":        {"give", ""},
	}

	for submit, expect := range tests {
		expected := NLPResult{expect[0], expect[1]}
		actual := IntentMap(submit)

		if actual.Action != expected.Action {
			t.Errorf("Actions differ:\n\tExpected:%s\n\tGot:%s", expected.Action, actual.Action)
		}

		if actual.Target != expected.Target {
			t.Errorf("Targets differ:\n\tExpected:%s\n\tGot:%s", expected.Target, actual.Target)
		}
	}
}
