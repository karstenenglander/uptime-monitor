package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"uptime-monitor/internal/hash"

	types "uptime-monitor/internal/types"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type TaskCreator interface {
	CreateHTTPTaskWithToken(projectID, locationID, queueID, url string, message *types.PollParams) (*taskspb.Task, error)
}

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func CreateHTTPTaskWithToken(projectID, locationID, queueID, url string, message *types.PollParams) (*taskspb.Task, error) {

	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build Task name
	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	taskId := message.Id
	currentTime := time.Now()
	hour := currentTime.Hour()
	minute := currentTime.Minute()
	taskWithTime := strconv.FormatInt(taskId, 10) + "_" + strconv.Itoa(hour) + "_" + strconv.Itoa(minute)
	taskNameHash := hash.GetFNVHash(taskWithTime)
	taskName := strconv.FormatUint(taskNameHash, 10)
	fullTaskName := queuePath + "/tasks/" + taskName

	serviceAccountEmail, serviceAccountEmailExists := os.LookupEnv("RUNTIME_SA")
	if !serviceAccountEmailExists {
		serviceAccountEmail = ""
	}

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			Name: fullTaskName,
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: serviceAccountEmail,
						},
					},
					Body:       messageJson,
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
				},
			},
		},
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %w", err)
	}

	return createdTask, nil
}
