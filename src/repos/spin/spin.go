package spin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/hugocortes/tools/src/domain"
	"github.com/hugocortes/tools/src/utils/rest"
)

const (
	applicationPipelineConfigPath = "/applications/%s/pipelineConfigs"
	applicationPipelinePath       = "/applications/%s/pipelines"
	pipelinesPath                 = "/pipelines"
	pipelinePath                  = "/pipelines/%s"
)

type repo struct {
	rest        *rest.Client
	accessToken string
}

func New(ctx context.Context, gateURL string, accessToken string) (domain.SpinRepo, error) {
	spin := &repo{
		accessToken: accessToken,
	}

	if spin.accessToken == "" {
		return nil, errors.New("access token not provided")
	}

	cookieClient, err := createHTTPClient()
	if err != nil {
		return nil, err
	}
	restClient, err := rest.NewClientFromHTTP(cookieClient, gateURL)
	if err != nil {
		return nil, err
	}

	spin.rest = restClient
	return spin, nil
}

func createHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Jar: jar,
	}, nil
}

func (r *repo) addHeaders(req *http.Request) {
	req.Header.Set("authorization", r.accessToken)
	req.Header.Set("Content-Type", "application/json")
}

func (r *repo) CreatePipeline(ctx context.Context, pipelineConfig *domain.PipelineConfig) (*domain.PipelineConfig, error) {
	path := pipelinesPath

	req, err := r.rest.NewRequest("POST", path, nil, pipelineConfig)
	if err != nil {
		return nil, err
	}
	r.addHeaders(req)

	newConfig := &domain.PipelineConfig{}
	_, err = r.rest.Do(ctx, req, newConfig)
	return newConfig, err
}

func (r *repo) GetApplicationPipelineConfigs(ctx context.Context, application string) ([]*domain.PipelineConfig, error) {
	path := fmt.Sprintf(applicationPipelineConfigPath, application)

	req, err := r.rest.NewRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}
	r.addHeaders(req)

	pipelineConfigs := &[]*domain.PipelineConfig{}
	_, err = r.rest.Do(ctx, req, pipelineConfigs)
	return *pipelineConfigs, err
}

func (r *repo) GetApplicationPipelineExecutions(ctx context.Context, application string) ([]*domain.PipelineExecution, error) {
	query := map[string]string{
		"expanded": "true",
		"limit":    "50",
	}
	path := fmt.Sprintf(applicationPipelinePath, application)

	req, err := r.rest.NewRequest("GET", path, query, nil)
	if err != nil {
		return nil, err
	}
	r.addHeaders(req)

	executions := &[]*domain.PipelineExecution{}
	_, err = r.rest.Do(ctx, req, executions)
	return *executions, err
}

func (r *repo) DeletePipelineExecution(ctx context.Context, ID string) error {
	path := fmt.Sprintf(pipelinePath, ID)

	req, err := r.rest.NewRequest("DELETE", path, nil, nil)
	if err != nil {
		return err
	}
	r.addHeaders(req)

	_, err = r.rest.Do(ctx, req, nil)
	return err
}
