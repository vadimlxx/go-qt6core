package qt6core

import "testing"

func TestVersion(t *testing.T) {
	maj, min, pat := Version()
	t.Logf("Qt %d.%d.%d", maj, min, pat)

	if maj != 6 {
		t.Fatalf("expected Qt 6.x, got %d.%d.%d", maj, min, pat)
	}
	if min < 11 {
		t.Fatalf("expected Qt >= 6.11, got %d.%d.%d", maj, min, pat)
	}
}
