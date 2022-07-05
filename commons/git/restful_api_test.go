package git

import "testing"

func TestGetAccessTokenByClient(t *testing.T) {
	t.Skip()

	accessToken, err := GetAccessTokenByClient(ClientID, ClientSecret, Code)
	if err != nil {
		t.Error(err)
	}

	t.Log(accessToken)
}

func TestGetUserByToken(t *testing.T) {
	// t.Skip()

	user, err := GetUserByToken(TestToken)
	if err != nil {
		t.Error(err)
	}

	t.Log(user)
}
