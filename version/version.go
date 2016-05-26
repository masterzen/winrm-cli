package version

import (
	"bytes"
	"fmt"
)

// GitSHA contains the Git SHA commit - this is filled by the compiler
var GitSHA string

// Version is the "marketing" version
const Version = "0.8.0"

func GetFullVersion() string {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "%s", Version)
	if len(GitSHA) >= 8 {
		fmt.Fprintf(&versionString, " (%s)", GitSHA[:8])
	}
	return versionString.String()
}
