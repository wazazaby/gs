package cache

import "testing"

func TestLRUOrder(t *testing.T) {
	var l LRU[int, string]
	l.Insert(1, "a")
	l.Insert(2, "b")
	l.Insert(3, "c")

	if got, ok := l.GetLRU(); !ok || got != "a" {
		t.Fatalf("GetLRU initial = %q, %v; want %q, true", got, ok, "a")
	}
	if got, ok := l.GetMRU(); !ok || got != "c" {
		t.Fatalf("GetMRU initial = %q, %v; want %q, true", got, ok, "c")
	}

	l.MakeMRU(1)
	if got, ok := l.GetLRU(); !ok || got != "b" {
		t.Fatalf("GetLRU after MakeMRU(1) = %q, %v; want %q, true", got, ok, "b")
	}
	if got, ok := l.GetMRU(); !ok || got != "a" {
		t.Fatalf("GetMRU after MakeMRU(1) = %q, %v; want %q, true", got, ok, "a")
	}

	l.MakeLRU(3)
	if got, ok := l.GetLRU(); !ok || got != "c" {
		t.Fatalf("GetLRU after MakeLRU(3) = %q, %v; want %q, true", got, ok, "c")
	}
	if got, ok := l.GetMRU(); !ok || got != "a" {
		t.Fatalf("GetMRU after MakeLRU(3) = %q, %v; want %q, true", got, ok, "a")
	}
}

func TestLRUUpsert(t *testing.T) {
	var l LRU[int, string]
	l.Upsert(1, "a")
	if got, ok := l.GetLRU(); !ok || got != "a" {
		t.Fatalf("GetLRU after Upsert = %q, %v; want %q, true", got, ok, "a")
	}

	l.Upsert(1, "b")
	if got, ok := l.GetMRU(); !ok || got != "b" {
		t.Fatalf("GetMRU after Upsert overwrite = %q, %v; want %q, true", got, ok, "b")
	}
}

func TestLRUInsert(t *testing.T) {
	var l LRU[int, string]
	l.Insert(1, "a")
	l.Insert(1, "b")

	if got, ok := l.GetLRU(); !ok || got != "a" {
		t.Fatalf("GetLRU after Insert duplicate = %q, %v; want %q, true", got, ok, "a")
	}
	if got, ok := l.GetMRU(); !ok || got != "a" {
		t.Fatalf("GetMRU after Insert duplicate = %q, %v; want %q, true", got, ok, "a")
	}
}

func TestLRUUpdate(t *testing.T) {
	var l LRU[int, string]
	l.Update(1, "a")
	if _, ok := l.GetMRU(); ok {
		t.Fatalf("GetMRU after Update missing ok = true; want false")
	}

	l.Upsert(1, "a")
	l.Upsert(2, "b")
	l.Update(1, "c")
	if got, ok := l.GetMRU(); !ok || got != "c" {
		t.Fatalf("GetMRU after Update existing = %q, %v; want %q, true", got, ok, "c")
	}
}

func TestLRURemove(t *testing.T) {
	var l LRU[int, string]
	l.Insert(1, "a")
	l.Insert(2, "b")

	l.Remove(1)
	if got, ok := l.GetLRU(); !ok || got != "b" {
		t.Fatalf("GetLRU after Remove(1) = %q, %v; want %q, true", got, ok, "b")
	}
	if got, ok := l.GetMRU(); !ok || got != "b" {
		t.Fatalf("GetMRU after Remove(1) = %q, %v; want %q, true", got, ok, "b")
	}

	l.Remove(2)
	if _, ok := l.GetLRU(); ok {
		t.Fatalf("GetLRU after Remove(2) ok = true; want false")
	}
	if _, ok := l.GetMRU(); ok {
		t.Fatalf("GetMRU after Remove(2) ok = true; want false")
	}
}

func TestLRUClear(t *testing.T) {
	var l LRU[int, string]
	l.Insert(1, "a")
	l.Insert(2, "b")
	l.Insert(3, "c")

	l.Clear()
	if _, ok := l.GetLRU(); ok {
		t.Fatalf("GetLRU after Clear ok = true; want false")
	}
	if _, ok := l.GetMRU(); ok {
		t.Fatalf("GetMRU after Clear ok = true; want false")
	}

	l.Insert(4, "d")
	if got, ok := l.GetLRU(); !ok || got != "d" {
		t.Fatalf("GetLRU after Clear/Add = %q, %v; want %q, true", got, ok, "d")
	}
	if got, ok := l.GetMRU(); !ok || got != "d" {
		t.Fatalf("GetMRU after Clear/Add = %q, %v; want %q, true", got, ok, "d")
	}
}
