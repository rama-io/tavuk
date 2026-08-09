package announce

import (
	"strings"
	"testing"
)

func TestBuildAnnouncementFullMessage(t *testing.T) {
	got := buildAnnouncement("Tui", "10", "- Fix bug where a query would never match again - Add media files filters - Checkbox for .nomedia folders", "1517897036098310245")

	want := `**Tui v10** has been released!

**What's Changed?**

- Fix bug where a query would never match again
- Add media files filters
- Checkbox for .nomedia folders

**Where to download?**
- It's already on GitHub releases. It will take some time for F-Droid update.

|| <@&1517897036098310245> ||`

	if got != want {
		t.Fatalf("message mismatch:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestBuildAnnouncementMentionsRoleOnce(t *testing.T) {
	got := buildAnnouncement("App", "1", "- one", "role1")
	if n := strings.Count(got, "<@&"); n != 1 {
		t.Fatalf("expected exactly one role mention, got %d:\n%q", n, got)
	}
}

func TestBuildAnnouncementOmitsRole(t *testing.T) {
	got := buildAnnouncement("App", "2", "- one", "")
	if got != "**App v2** has been released!\n\n**What's Changed?**\n\n- one\n\n**Where to download?**\n- It's already on GitHub releases. It will take some time for F-Droid update.\n" {
		t.Fatalf("unexpected message without role:\n%q", got)
	}
}

func TestBuildAnnouncementOmitsChangesSectionWhenEmpty(t *testing.T) {
	got := buildAnnouncement("App", "3", "   -   ", "1")
	if got != "**App v3** has been released!\n\n**Where to download?**\n- It's already on GitHub releases. It will take some time for F-Droid update.\n\n|| <@&1> ||" {
		t.Fatalf("unexpected message without changes:\n%q", got)
	}
}

func TestSplitChanges(t *testing.T) {
	got := splitChanges("- feat 1 - feat 2 - feat 3")
	want := []string{"feat 1", "feat 2", "feat 3"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestReleaseTitle(t *testing.T) {
	cases := []struct {
		app, version, want string
	}{
		{"Txori", "10", "Txori v10"},
		{"Txori", "v10", "Txori v10"},
		{"Txori", "1.5.3", "Txori v1.5.3"},
		{"Txori", "  ", "Txori"},
	}
	for _, c := range cases {
		if got := releaseTitle(c.app, c.version); got != c.want {
			t.Fatalf("releaseTitle(%q, %q) = %q, want %q", c.app, c.version, got, c.want)
		}
	}
}
