package service

type User struct {
	// Basic Info
	Name  string `json:"name"`
	Email string `json:"email"`
	// Git Info
	GitUser
}

func FindUserByCode(clientId, clientSecret, code string) (*User, error) {
	user, err := GetUserByGitCode(clientId, clientSecret, code)
	if err != nil {
		return nil, err
	}
	// TODO Replenish user info

	return &User{
		GitUser: *user,
	}, nil
}
