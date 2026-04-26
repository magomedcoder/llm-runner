package huggingface

import "testing"

func TestEffectiveParallelDownload(t *testing.T) {
	if g := effectiveParallelDownload(0, 1); g != 1 {
		t.Fatalf("single file: %d", g)
	}

	if g := effectiveParallelDownload(0, 10); g != 4 {
		t.Fatalf("auto cap 4: %d", g)
	}

	if g := effectiveParallelDownload(6, 10); g != 6 {
		t.Fatalf("explicit 6: %d", g)
	}

	if g := effectiveParallelDownload(100, 3); g != 3 {
		t.Fatalf("cap by nFiles: %d", g)
	}

	if g := effectiveParallelDownload(1, 10); g != 1 {
		t.Fatalf("sequential: %d", g)
	}

	if g := effectiveParallelDownload(8, 100); g != 8 {
		t.Fatalf("max 8: %d", g)
	}
}
