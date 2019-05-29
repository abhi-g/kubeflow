package projects

import (
	"testing"

	crm "google.golang.org/api/cloudresourcemanager/v1"
)

func TestProjectCreation(t *testing.T) {
	// Following is the folder ID of "gcp-deploy" folder under kubeflow.org
	folderID := &crm.ResourceId{Type: "folder", Id: "838562927550"}
	if err := CreateProject("test-project-lib-folder", folderID); err != nil {
		t.Error(err)
	}
}
