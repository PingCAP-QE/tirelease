// user_application is used to parse the reqeust and call the service layer to handle the request.
package controller

import (
	"tirelease/internal/service"

	"github.com/pkg/errors"
)

func HandleFindUser(request UserRequest) (UserResp, error) {

	// Find user by git code
	if request.ByGitCode() {
		err := request.Validate()
		if err != nil {
			return UserResp{}, err
		}

		user, err := service.FindUserByCode(*request.GitClientID, *request.GitClientSecret, *request.GitCode)

		if err != nil {
			return UserResp{}, err
		}

		return UserResp{*user}, nil
	}

	return UserResp{}, errors.New("unsupported request")

}
