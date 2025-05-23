package models

import (
	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
)

type ServerType string

const (
	ServerTypeMyOwn          ServerType = "My Own"
	ServerTypeGaming         ServerType = "Gaming"
	ServerTypeFriend         ServerType = "Friend"
	ServerTypeStudyGroup     ServerType = "Study Group"
	ServerTypeSchoolClub     ServerType = "School Club"
	ServerTypeLocalCommunity ServerType = "Local Community"
	ServerTypeArtists        ServerType = "Artists & Creators"
)

const (
	PermViewChannel = "view_channel"
	PermEditChannel = "edit_channel"
	PermEditServer  = "edit_server"
	PermManageRoles = "manage_roles"
	PermInviteUsers = "invite_users"
	PermKickMembers = "kick_members"
	PermOwner       = "owner"
)

type Server struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string     `json:"name"`
	IconURL     string     `json:"icon_url,omitempty"`
	BannerURL   string     `json:"banner_url,omitempty"`
	Description string     `json:"description,omitempty"`
	Traits      []string   `json:"traits" gorm:"type:text[]"`
	Type        ServerType `json:"type"`
	IsPrivate   bool       `json:"is_private"`
	OwnerID     uuid.UUID  `json:"owner_id"` // FK to User

	// Relationships
	// Members []ServerMember `json:"members" gorm:"foreignKey:ServerID"`
	// Roles   []ServerRole   `json:"roles" gorm:"foreignKey:ServerID"`
	// Invites []ServerInvite `json:"invites" gorm:"foreignKey:ServerID"`
	// Rules   []ServerRule   `json:"rules" gorm:"foreignKey:ServerID"`

	models.AuditFields `json:",inline"`
}
