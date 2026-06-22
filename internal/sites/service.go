package sites

import (
	"context"
	"log"
	"os"
	repo "uptime-monitor/internal/adapters/postgresql/sqlc"
	tasks "uptime-monitor/internal/tasks"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type Service interface {
	ListSites(ctx context.Context) ([]repo.Site, error)
	EnqueuePollSites(ctx context.Context) ([]*cloudtaskspb.Task, error)
	// pollSite(ctx context.Context) (string, error)
	AddSite(ctx context.Context, params createAddParams) (int64, error)
	RemoveSite(ctx context.Context, params createIdParams) (string, error)
	FindSiteByID(ctx context.Context, params createIdParams) (repo.Site, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListSites(ctx context.Context) ([]repo.Site, error) {
	return s.repo.ListSites(ctx)
}

func (s *svc) EnqueuePollSites(ctx context.Context) ([]*cloudtaskspb.Task, error) {
	sites, err := s.repo.ListSites(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	projectID, projectIDExists := os.LookupEnv("PROJECT_ID")
	if !projectIDExists {
		projectID = ""
	}
	locationID, locationIDExists := os.LookupEnv("LOCATION_ID")
	if !locationIDExists {
		locationID = ""
	}
	queueID, queueIDExists := os.LookupEnv("QUEUE_ID")
	if !queueIDExists {
		queueID = ""
	}

	var enqueued []*cloudtaskspb.Task
	for _, v := range sites {
		var url = v.Url
		var message = v.Name
		task, err := tasks.CreateHTTPTask(projectID, locationID, queueID, url, message)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		enqueued = append(enqueued, task)
	}

	return enqueued, nil
}

// func (s *svc) pollSite(ctx context.Context) (string, error) {

// 	return
// }

func (s *svc) FindSiteByID(ctx context.Context, params createIdParams) (repo.Site, error) {
	id := params.Id
	return s.repo.FindSiteByID(ctx, id)
}

func (s *svc) AddSite(ctx context.Context, params createAddParams) (int64, error) {
	name := params.Name
	url := params.Url
	return s.repo.AddSite(ctx, repo.AddSiteParams{Name: name, Url: url})
}

func (s *svc) RemoveSite(ctx context.Context, params createIdParams) (string, error) {
	id := params.Id
	return s.repo.RemoveSiteByID(ctx, id)
}
