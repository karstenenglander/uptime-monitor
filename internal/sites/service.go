package sites

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	repo "uptime-monitor/internal/adapters/postgresql/sqlc"
	tasks "uptime-monitor/internal/tasks"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListSites(ctx context.Context) ([]repo.Site, error)
	EnqueuePollSites(ctx context.Context) ([]*cloudtaskspb.Task, error)
	PollSite(ctx context.Context, params pollParams) (*http.Response, error)
	AddSite(ctx context.Context, params createAddParams) (int64, error)
	RemoveSite(ctx context.Context, params createIdParams) (string, error)
	FindSitesByID(ctx context.Context, params createIdParams) (repo.Site, error)
	updateSitePolled(ctx context.Context, id int64) (int64, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListSites(ctx context.Context) ([]repo.Site, error) {
	sites, err := s.repo.ListSites(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sites, nil
}

func (s *svc) updateSitePolled(ctx context.Context, id int64) (int64, error) {
	currentTime := pgtype.Timestamptz{Time: time.Now(), InfinityModifier: 0, Valid: true}
	params := repo.UpdateSitePolledParams{ID: id, PolledAt: currentTime}
	id, err := s.repo.UpdateSitePolled(ctx, params)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return id, nil
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

	endpointURL, endpointURLExists := os.LookupEnv("ENDPOINT_URL")
	if !endpointURLExists {
		endpointURL = ""
	}

	var enqueued []*cloudtaskspb.Task
	for _, v := range sites {
		// The message is the URL of the site to be polled
		var message = v.Url
		task, err := tasks.CreateHTTPTask(projectID, locationID, queueID, endpointURL, message)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		enqueued = append(enqueued, task)
	}

	return enqueued, nil
}

func (s *svc) PollSite(ctx context.Context, params pollParams) (*http.Response, error) {
	url := params.Url
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// TO-DO: Add checking of return code to determine action
	updateId := params.Id
	_, err = s.updateSitePolled(ctx, updateId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func (s *svc) FindSitesByID(ctx context.Context, params createIdParams) (repo.Site, error) {
	id := params.Id
	return s.repo.FindSitesByID(ctx, id)
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
