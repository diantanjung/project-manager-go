package service

import "testing"

func TestNormalizePage(t *testing.T) {
	got := NormalizePage(PageFilter{Page: -1, Limit: 999})
	if got.Page != 1 || got.Limit != 100 {
		t.Fatalf("got %#v, want page 1 limit 100", got)
	}
}
