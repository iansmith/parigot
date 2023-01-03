package go_

import (
	"fmt"
	"testing"
)

func TestFilename(t *testing.T) {
	badFilename := []string{
		"",             // nonsensical
		".",            // we require full pathnames (no CWD action)
		"..",           // no chance
		"/..",          // probably an attempt to subvert fs
		"/.",           // probably an attempt to subvert fs
		"/foo",         // must start with /app
		"/foo/bar/baz", // must start with /app
	}
	controlChar := []string{
		"\x00",
		"\t",
		"\n",
		"\r",
		"\f",
		"\v",
	}
	controlCharTest := []string{
		"%s",                     // sanity
		"/app/foo%s",             //part of filename
		"/app/foo%sbar",          //part of filename
		"/app/%s/foo",            // is directory
		"/app/foo/fleazil%sfrak", // in middle of text of filename
		"/foo/%s",                //fails because not /app
		"/app/%s/foo",            // as dir
		"/app/bar%s/foo",         // part of dir
		"/foo/bleah%sgack/foo",   // part of dir
	}

	okFilename := []string{
		"/app",
		"/app/",                              // bad trailing / cleaned up by lexical processing
		"/app/foo/",                          // bad trailing / cleaned up by lexical processing
		"/foo/../app/baz",                    // ok, because result of .. being processed lexically is /app
		"/foo/./bar/../..///app///////gack/", // ok, because result of all this cruft being processed lexically is /app/gack
		"/app/baz/fleazil",                   // simple case
	}

	for _, bad := range badFilename {
		_, err := ValidatePathForParigot(bad, "open")
		if err == nil {
			t.Errorf("path '%s' succeeded but should have failed", bad)
		}
	}
	for _, good := range okFilename {
		_, err := ValidatePathForParigot(good, "open")
		if err == nil {
			t.Errorf("path '%s' failed", good)
		}
	}
	for _, testCase := range controlCharTest {
		for _, ch := range controlChar {
			bad := fmt.Sprintf(testCase, ch)
			_, err := ValidatePathForParigot(bad, "open")
			if err == nil {
				t.Errorf("path '%s' succeeded but should have failed", bad)
			}
		}
	}
}
