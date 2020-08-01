package spin

import (
	"context"
	"fmt"

	"github.com/hugocortes/tools/api"
)

const (
	appPipelinePath = "/applications/%s/pipelines"
	pipelinePath    = "/pipelines/%s"

	Succeeded = "SUCCEEDED"
	Failed    = "FAILED"
	Running   = "RUNNING"
)

type Orca struct {
	client *api.Client
}

type Execution struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewOrca(orcaURL string) (*Orca, error) {
	orca := &Orca{}

	client, err := api.NewClient(orcaURL)
	if err != nil {
		return nil, err
	}

	orca.client = client
	return orca, nil
}

func (o *Orca) GetExecutions(appID string) ([]*Execution, error) {
	query := map[string]string{
		"expanded": "true",
		"limit":    "50",
	}
	path := fmt.Sprintf(appPipelinePath, appID)

	req, err := o.client.NewRequest("GET", path, query)
	if err != nil {
		return nil, err
	}

	executions := &[]*Execution{}
	_, err = o.client.Do(context.Background(), req, executions)
	if err != nil {
		return nil, err
	}

	return *executions, nil
}

func (o *Orca) DeletePipeline(ID string) error {
	path := fmt.Sprintf(pipelinePath, ID)

	req, err := o.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = o.client.Do(context.Background(), req, nil)
	return err
}
