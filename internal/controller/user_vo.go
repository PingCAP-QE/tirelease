package controller

import (
	"errors"
	"tirelease/internal/service"
)

type UserRequest struct {
	GitUserRequest
}

type UserResp struct {
	service.User
}

type GitUserRequest struct {
	GitCode         *string `json:"git_code,omitempty" form:"git_code"`
	GitClientID     *string `json:"git_client_id,omitempty" form:"git_client_id"`
	GitClientSecret *string `json:"git_client_secret,omitempty" form:"git_client_secret"`
}

func (r *GitUserRequest) ByGitCode() bool {
	return r.GitCode != nil
}

func (r *GitUserRequest) Validate() error {
	if r.GitCode == nil {
		return errors.New("git_code is required")
	}
	if r.GitClientID == nil {
		return errors.New("git_client_id is required")
	}
	if r.GitClientSecret == nil {
		return errors.New("git_client_secret is required")
	}
	return nil
}
