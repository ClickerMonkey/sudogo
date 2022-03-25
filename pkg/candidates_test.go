package sudogo

import (
	"fmt"
	"testing"
)

func TestBits(t *testing.T) {
	b := Candidates{}

	if b.Count != 0 {
		t.Error("default count is not zero")
	}
	if b.Value != 0 {
		t.Error("default value is not zero")
	}

	b.Fill(4)

	if b.Value != 30 {
		t.Error("fill has wrong value")
	}
	if b.Count != 4 {
		t.Error("fill has wrong count")
	}

	if b.Set(1, true) {
		t.Error("bit 1 is already set and I was allowed to set it")
	}
	if b.Set(2, true) {
		t.Error("bit 2 is already set and I was allowed to set it")
	}
	if b.Set(3, true) {
		t.Error("bit 3 is already set and I was allowed to set it")
	}
	if b.Set(4, true) {
		t.Error("bit 4 is already set and I was allowed to set it")
	}
	if b.Count != 4 {
		t.Error("set of already set bits affected count")
	}

	if !b.Set(1, false) {
		t.Error("bit 1 could not be set")
	}
	if b.Count != 3 {
		t.Error("count wrong after set 1")
	}
	if fmt.Sprint(b.ToSlice()) != "[2 3 4]" {
		t.Error("setting bit 0 resulted in the wrong slice")
	}

	if !b.Set(2, false) {
		t.Error("bit 2 could not be set")
	}
	if b.Count != 2 {
		t.Error("count wrong after set 2")
	}
	if fmt.Sprint(b.ToSlice()) != "[3 4]" {
		t.Error("setting bit 1 resulted in the wrong slice")
	}

	if !b.Set(3, false) {
		t.Error("bit 3 could not be set")
	}
	if b.Count != 1 {
		t.Error("count wrong after set 3")
	}
	if fmt.Sprint(b.ToSlice()) != "[4]" {
		t.Error("setting bit 3 resulted in the wrong slice")
	}

	if !b.Set(4, false) {
		t.Error("bit 4 could not be set")
	}
	if b.Count != 0 {
		t.Error("count wrong after set 4")
	}
	if fmt.Sprint(b.ToSlice()) != "[]" {
		t.Error("setting bit 4 resulted in the wrong slice")
	}
}
