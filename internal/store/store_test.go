package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenMissingFileCreatesEmptyStore(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "apps.json"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, ok := s.Get("nope"); ok {
		t.Fatal("expected empty store")
	}
}

func TestSetPersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "apps.json")

	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	target := Target{GuildID: "g1", ChannelID: "c1", RoleID: "r1"}
	if err := s.Set("Tui", target); err != nil {
		t.Fatalf("Set: %v", err)
	}

	s2, err := Open(path)
	if err != nil {
		t.Fatalf("Reopen: %v", err)
	}
	got, ok := s2.Get("Tui")
	if !ok {
		t.Fatal("expected Tui to be present after reopen")
	}
	if got != target {
		t.Fatalf("got %+v, want %+v", got, target)
	}
}

func TestSetOverwriteUpdatesTarget(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "apps.json"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := s.Set("Tui", Target{ChannelID: "c1"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Set("Tui", Target{ChannelID: "c2", RoleID: "r2"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, _ := s.Get("Tui")
	if got.ChannelID != "c2" || got.RoleID != "r2" {
		t.Fatalf("got %+v, want ChannelID=c2 RoleID=r2", got)
	}
}

func TestNoTempFileLeftBehind(t *testing.T) {
	path := filepath.Join(t.TempDir(), "apps.json")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := s.Set("a", Target{ChannelID: "c"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp file should have been renamed away")
	}
}

func TestAppsSorted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "apps.json")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	for _, app := range []string{"Tui", "Txori", "Beta"} {
		if err := s.Set(app, Target{ChannelID: "c1"}); err != nil {
			t.Fatalf("Set(%q): %v", app, err)
		}
	}
	want := []string{"Beta", "Tui", "Txori"}
	got := s.Apps()
	if len(got) != len(want) {
		t.Fatalf("Apps() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Apps() = %v, want %v", got, want)
		}
	}
}
