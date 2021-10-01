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
	appPipelinePath = "/applications/%s/pipelines"
	pipelinePath    = "/pipelines/%s"
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

func (r *repo) addAuthn(req *http.Request) {
	req.Header.Set("authorization", r.accessToken)
}

func (r *repo) GetApplicationPipelineExecutions(ctx context.Context, appID string) ([]*domain.PipelineExecution, error) {
	query := map[string]string{
		"expanded": "true",
		"limit":    "50",
	}
	path := fmt.Sprintf(appPipelinePath, appID)

	req, err := r.rest.NewRequest("GET", path, query)
	if err != nil {
		return nil, err
	}
	r.addAuthn(req)

	executions := &[]*domain.PipelineExecution{}
	_, err = r.rest.Do(context.Background(), req, executions)
	if err != nil {
		return nil, err
	}

	return *executions, nil
}

func (r *repo) DeletePipelineExecution(ctx context.Context, ID string) error {
	path := fmt.Sprintf(pipelinePath, ID)

	req, err := r.rest.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	r.addAuthn(req)

	_, err = r.rest.Do(context.Background(), req, nil)
	return err
}
