package version

import (
	"testing"
)

func TestGetMavenVersions(t *testing.T) {
    versions := GetMavenVersions()
    t.Log(versions)
}
