package commands

import "testing"

func TestTagsFromFlags(t *testing.T) {
	tags, err := tagsFromFlags([]string{"signal=EntraID_AiTM_Session", " signal = X_local ", "q=a=b"})
	if err != nil {
		t.Fatal(err)
	}
	want := [][2]string{{"signal", "EntraID_AiTM_Session"}, {"signal", "X_local"}, {"q", "a=b"}}
	if len(tags) != len(want) {
		t.Fatalf("tags %d", len(tags))
	}
	for i, w := range want {
		if tags[i].Name != w[0] || tags[i].Value != w[1] {
			t.Errorf("tag %d = %s=%s, want %s=%s", i, tags[i].Name, tags[i].Value, w[0], w[1])
		}
	}
	if tags, err := tagsFromFlags(nil); err != nil || tags != nil {
		t.Errorf("no flags: %v %v", tags, err)
	}
	for _, bad := range []string{"signal", "=x", "signal=", " = "} {
		if _, err := tagsFromFlags([]string{bad}); err == nil {
			t.Errorf("%q accepted", bad)
		}
	}
}
