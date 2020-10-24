// Package slack - user
package slack

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	userApiResponse struct {
		// Ok indicates success or failure
		Ok bool `json:"ok"`

		// Error short machine-readable error code
		Error string `json:"error"`

		// User contains user information
		User User `json:"user"`
	}

	// User entity
	User struct {
		// ID identifier for this workspace user. It is unique to the workspace containing the user. Use this field
		// together with team_id as a unique key when storing related data or when specifying the user in API requests.
		ID string `json:"id"`

		// TeamID
		TeamID string `json:"team_id"`

		// Name deprecated field
		Name string `json:"name"`

		// Deleted this user has been deactivated when the value of this field is true. Otherwise the value is false,
		// or the field may not appear at all.
		Deleted bool `json:"deleted"`

		// Color used in some clients to display a special username color.
		Color string `json:"color"`

		// RealName the real name that the user specified in their workspace profile.
		RealName string `json:"real_name"`

		// TimeZone a human-readable string for the geographic timezone-related region this user has specified in their
		// account.
		TimeZone string `json:"tz"`

		// TimeZoneLabel вescribes the commonly used name of the tz timezone.
		TimeZoneLabel string `json:"tz_label"`

		// TimeZoneOffset indicates the number of seconds to offset UTC time by for this user's tz.
		TimeZoneOffset int `json:"tz_offset"`

		// Profile an object containing the default fields of a user's workspace profile
		Profile struct {
			// Title
			Title string `json:"title"`

			// Phone user's phone
			Phone string `json:"phone"`

			// Skype user's skype
			Skype string `json:"skype"`

			// RealName the real name that the user specified in their workspace profile.
			RealName string `json:"real_name"`

			// RealNameNormalized the real_name field, but with any non-Latin characters filtered out.
			RealNameNormalized string `json:"real_name_normalized"`

			// DisplayName indicates the display name that the user has chosen to identify themselves by in their
			// workspace profile.
			DisplayName string `json:"display_name"`

			// DisplayNameNormalized the display_name field, but with any non-Latin characters filtered out.
			DisplayNameNormalized string `json:"display_name_normalized"`

			// StatusText text of user's status
			StatusText string `json:"status_text"`

			// StatusEmoji emoji of user's status
			StatusEmoji string `json:"status_emoji"`

			// StatusExpiration expiration of user's status
			StatusExpiration int `json:"status_expiration"`

			// AvatarHash
			AvatarHash string `json:"avatar_hash"`

			// FirstName user's first name
			FirstName string `json:"first_name"`

			// LastName user's last name
			LastName string `json:"last_name"`

			// Email user's email
			Email string `json:"email"`

			// ImageOriginal these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			ImageOriginal string `json:"image_original"`

			// Image24 these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			Image24 string `json:"image_24"`

			// Image32 these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			Image32 string `json:"image_32"`

			// Image48 these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			Image48 string `json:"image_48"`

			// Image72 these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			Image72 string `json:"image_72"`

			// Image192 these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			Image192 string `json:"image_192"`

			// Image512 these various fields will contain https URLs that point to square ratio, web-viewable
			// images (GIFs, JPEGs, or PNGs) that represent different sizes of a user's profile picture.
			Image512 string `json:"image_512"`

			// Team
			Team string `json:"team"`
		} `json:"profile"`

		// IsAdmin indicates whether the user is an Admin of the current workspace.
		IsAdmin bool `json:"is_admin"`

		// IsOwner indicates whether the user is an Owner of the current workspace.
		IsOwner bool `json:"is_owner"`

		// IsPrimaryOwner indicates whether the user is the Primary Owner of the current workspace.
		IsPrimaryOwner bool `json:"is_primary_owner"`

		// IsRestricted indicates whether or not the user is a guest user. Use in combination with the
		// is_ultra_restricted field to check if the user is a single-channel guest user.
		IsRestricted bool `json:"is_restricted"`

		// IsUltraRestricted indicates whether or not the user is a single-channel guest.
		IsUltraRestricted bool `json:"is_ultra_restricted"`

		// IsBot indicates whether the user is actually a bot user. Bleep bloop.
		IsBot bool `json:"is_bot"`

		// IsStranger if true, this user belongs to a different workspace than the one associated with your app's token,
		// and isn't in any shared channels visible to your app. If false (or this field is not present), the user is
		// either from the same workspace as associated with your app's token, or they are from a different workspace,
		// but are in a shared channel that your app has access to.
		IsStranger bool `json:"is_stranger"`

		// Updated a unix timestamp indicating when the user object was last updated.
		Updated int `json:"updated"`

		// IsAppUser indicates whether the user is an authorized user of the calling app.
		IsAppUser bool `json:"is_app_user"`

		// IsInvitedUser indicates whether the user signed up for the workspace directly (false or the field is absent),
		// or they joined via an invite (true).
		IsInvitedUser bool `json:"is_invited_user"`

		// Has2FA describes whether two-factor authentication is enabled for this user.
		Has2FA bool `json:"has_2fa"`

		// Locale contains a IETF language code that represents this user's chosen display language for Slack clients.
		Locale string `json:"locale"`
	}
)

func getUserByEmail(ctx context.Context, c *client, email string) (User, error) {
	respBody, err := c.get(ctx, fmt.Sprintf("/api/users.lookupByEmail?email=%s", email))
	if err != nil {
		return User{}, err
	}

	var resp userApiResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return User{}, fmt.Errorf("can't unmarshal response: %w", err)
	}

	if resp.Ok != true {
		return User{}, fmt.Errorf("slack respond with error: %s", resp.Error)
	}

	return resp.User, nil
}
