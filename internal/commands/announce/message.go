package announce

import (
	"fmt"
	"strings"
)

const downloadNote = "It's already on GitHub releases. It will take some time for F-Droid update."

const (
	releasedFmt        = "**%s** has been released!\n\n"
	whatsChangedHeader = "**What's Changed?**"
	whereToDownloadFmt = "**Where to download?**\n- %s\n"
	roleMentionFmt     = "\n|| <@&%s> ||"
)

// combines the app name and version
func releaseTitle(app, version string) string {
	v := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(version), "v"))
	if v == "" {
		return app
	}
	return app + " v" + v
}

func buildAnnouncement(app, version, changes, roleID string) string {
	var b strings.Builder

	fmt.Fprintf(&b, releasedFmt, releaseTitle(app, version))

	if feats := splitChanges(changes); len(feats) > 0 {
		b.WriteString(whatsChangedHeader + "\n\n")
		for _, f := range feats {
			fmt.Fprintf(&b, "- %s\n", f)
		}
		b.WriteString("\n")
	}

	fmt.Fprintf(&b, whereToDownloadFmt, downloadNote)

	if roleID != "" {
		fmt.Fprintf(&b, roleMentionFmt, roleID)
	}

	return b.String()
}

// splitChanges turns " - a - b" into ["a", "b"], dropping empty entries.
func splitChanges(s string) []string {
	var out []string
	for _, part := range strings.Split(s, " - ") {
		part = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(part), "-"))
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
