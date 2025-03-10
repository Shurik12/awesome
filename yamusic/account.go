package yamusic

import (
	"context"
	"net/http"
	"time"
)

type (
	// AccountService is a service to deal with accounts.
	AccountService struct {
		client *Client
	}
	// AccountStatusResp describes account get status method response
	AccountStatusResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         struct {
			Account struct {
				UID              int       `json:"uid"`
				Region           int       `json:"region"`
				Login            string    `json:"login"`
				FullName         string    `json:"fullName"`
				SecondName       string    `json:"secondName"`
				FirstName        string    `json:"firstName"`
				DisplayName      string    `json:"displayName"`
				Birthday         string    `json:"birthday"`
				ServiceAvailable bool      `json:"serviceAvailable"`
				HostedUser       bool      `json:"hostedUser"`
				RegisteredAt     time.Time `json:"registeredAt"`
				Now              time.Time `json:"now"`
				PassportPhones   []struct {
					Phone string `json:"phone"`
				} `json:"passport-phones"`
			} `json:"account"`
			Permissions struct {
				Until   time.Time `json:"until"`
				Values  []string  `json:"values"`
				Default []string  `json:"default"`
			} `json:"permissions"`
			Subscription struct {
				CanStartTrial bool `json:"canStartTrial"`
				Mcdonalds     bool `json:"mcdonalds"`
			} `json:"subscription"`
		} `json:"result"`
	}
	AccountSettingsResp struct {
		InvocationInfo InvocationInfo `json:"invocationInfo"`
		Error          Error          `json:"error"`
		Result         struct {
			UID                       int       `json:"uid"`
			LastFmScrobblingEnabled   bool      `json:"lastFmScrobblingEnabled"`
			FacebookScrobblingEnabled bool      `json:"facebookScrobblingEnabled"`
			ShuffleEnabled            bool      `json:"shuffleEnabled"`
			AddNewTrackOnPlaylistTop  bool      `json:"addNewTrackOnPlaylistTop"`
			VolumePercents            int       `json:"volumePercents"`
			UserMusicVisibility       bool      `json:"userMusicVisibility"`
			UserSocialVisibility      bool      `json:"userSocialVisibility"`
			AdsDisabled               bool      `json:"adsDisabled"`
			Modified                  time.Time `json:"modified"`
			RbtDisabled               bool      `json:"rbtDisabled"`
			Theme                     string    `json:"theme"`
			PromosDisabled            bool      `json:"promosDisabled"`
			AutoPlayRadio             bool      `json:"autoPlayRadio"`
			SyncQueueEnabled          bool      `json:"syncQueueEnabled"`
			ExplicitForbidden         bool      `json:"explicitForbidden"`
			ChildModEnabled           bool      `json:"childModEnabled"`
			WizardIsPassed            bool      `json:"wizardIsPassed"`
			UserCollectionHue         int       `json:"userCollectionHue"`
		} `json:"result"`
	}
)

// GetStatus returns account's status
func (s *AccountService) GetStatus(
	ctx context.Context,
) (*AccountStatusResp, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "account/status", nil)
	if err != nil {
		return nil, nil, err
	}

	accountStatus := new(AccountStatusResp)
	resp, err := s.client.Do(ctx, req, accountStatus)
	return accountStatus, resp, err
}

// GetStatus returns account's status
func (s *AccountService) GetUser(
	ctx context.Context,
) (*AccountSettingsResp, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "account/settings", nil)
	if err != nil {
		return nil, nil, err
	}

	accountStatus := new(AccountSettingsResp)
	resp, err := s.client.Do(ctx, req, accountStatus)
	return accountStatus, resp, err
}
