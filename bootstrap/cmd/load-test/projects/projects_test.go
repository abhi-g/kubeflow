package projects

import (
	"strings"
	"testing"
)

func TestProjectCreation(t *testing.T) {
	org, err := getKubeflowOrg()
	if err != nil {
		t.Error(err)
	}
	orgID := strings.Split(org.Name, "/")[1]
	if err := CreateProject("test-project-lib", orgID); err != nil {
		t.Error(err)
	}
}
