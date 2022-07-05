package service

import "tirelease/commons/git"

type GitUser struct {
	GitID        int64  `json:"git_id"`
	GitLogin     string `json:"git_login"`
	GitAvatarURL string `json:"git_avatar_url"`
	GitName      string `json:"git_name"`
}

func GetUserByGitCode(clientId, clientSecret, code string) (*GitUser, error) {
	accessToken, err := git.GetAccessTokenByClient(clientId, clientSecret, code)
	if err != nil {
		return nil, err
	}

	user, err := git.GetUserByToken(accessToken)
	if err != nil {
		return nil, err
	}

	return &GitUser{
		GitID:        user.GetID(),
		GitLogin:     user.GetLogin(),
		GitAvatarURL: user.GetAvatarURL(),
		GitName:      user.GetName(),
	}, nil
}
