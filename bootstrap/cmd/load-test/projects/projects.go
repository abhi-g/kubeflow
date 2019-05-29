// Package projects provides helpers to create GCP projects
package projects

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	crm "google.golang.org/api/cloudresourcemanager/v1"
)

var projectPrefix = "kf-load-test-project"
var ctx = context.Background()

func getKubeflowOrg() (*crm.Organization, error) {
	crmService, err := crm.NewService(ctx)
	var org *crm.Organization
	if err != nil {
		return nil, err
	}

	orgSearch := crm.SearchOrganizationsRequest{
		Filter: "domain:kubeflow.org",
	}
	req := crmService.Organizations.Search(&orgSearch)
	if err := req.Pages(ctx, func(page *crm.SearchOrganizationsResponse) error {
		org = page.Organizations[0]
		return nil
	}); err != nil {
		return nil, err
	}
	return org, nil
}

// CreateProject creates a GCP project for the load test with the given name
// under the resId resource which can be an organization or a folder.
func CreateProject(name string, resID *crm.ResourceId) error {
	log.Printf("Creating GCP Project: %v", name)
	p := crm.Project{Name: name, ProjectId: name, Parent: resID}

	crmService, err := crm.NewService(ctx)
	if err != nil {
		return err
	}

	op, err := crmService.Projects.Create(&p).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "alreadyExists") {
			log.Print("Project already exists")
			return nil
		}
		return err
	}

	status := crm.ProjectCreationStatus{}

	for i := 0; !status.Gettable && i < 10; i++ {
		time.Sleep(5 * time.Second)
		op, err = crmService.Operations.Get(op.Name).Context(ctx).Do()
		if err != nil {
			return err
		}

		json.Unmarshal(op.Metadata, &status)
		log.Print(status)
		log.Printf("Gettable: %v, Ready: %v",
			status.Gettable,
			status.Ready)
	}

	if !status.Gettable {
		return errors.New("Failed to create project")
	}
	return nil
}

// DeleteProject will delete the project with the given name
// as the project ID.
func DeleteProject(name string) error {
	// TODO: Do we need this? After delete, the project needs to be restored, and
	// cannot be created again within 30days?
	log.Printf("Deleting GCP Project: %v", name)
	return nil
}

// CreateAllProjects creates all num GCP projects for the load test.
func CreateAllProjects(num int) error {
	log.Print("Creating ", num, " GCP projects for the load test")
	// Following is the folder ID of "gcp-deploy" folder under kubeflow.org
	folderID := &crm.ResourceId{Type: "folder", Id: "838562927550"}
	for i := 0; i < num; i++ {
		projName := projectPrefix + strconv.Itoa(i)
		if err := CreateProject(projName, folderID); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}
