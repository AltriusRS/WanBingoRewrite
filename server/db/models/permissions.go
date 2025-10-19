package models

// Permission represents a user permission as a bit flag
type Permission uint64

// Permission bit flags - using powers of 2 for bitwise operations
const (
	// Basic permissions
	PermCanChat Permission = 1 << iota
	PermCanSendMessages
	PermCanSendWhispers
	PermCanDeleteOwnMessages

	// Moderation permissions
	PermCanDeleteMessages
	PermCanBanUsers
	PermCanKickUsers
	PermCanMuteUsers
	PermCanUnmuteUsers

	// Host/Admin permissions
	PermCanHost
	PermCanModerate
	PermCanManageChat
	PermCanManageChatPermissions

	// Content management permissions
	PermCanSuggestTiles
	PermCanReviewTiles
	PermCanApproveTiles
	PermCanManageTiles

	// Advanced permissions
	PermCanPromotePlayers
	PermCanModifyShowData
	PermCanManageTimers
)

// PermissionNames maps permission bits to human-readable names
var PermissionNames = map[Permission]string{
	PermCanChat:                  "can_chat",
	PermCanSendMessages:          "can_send_messages",
	PermCanSendWhispers:          "can_send_whispers",
	PermCanDeleteOwnMessages:     "can_delete_own_messages",
	PermCanDeleteMessages:        "can_delete_messages",
	PermCanBanUsers:              "can_ban_users",
	PermCanKickUsers:             "can_kick_users",
	PermCanMuteUsers:             "can_mute_users",
	PermCanUnmuteUsers:           "can_unmute_users",
	PermCanHost:                  "can_host",
	PermCanModerate:              "can_moderate",
	PermCanManageChat:            "can_manage_chat",
	PermCanManageChatPermissions: "can_manage_chat_permissions",
	PermCanSuggestTiles:          "can_suggest_tiles",
	PermCanReviewTiles:           "can_review_tiles",
	PermCanApproveTiles:          "can_approve_tiles",
	PermCanManageTiles:           "can_manage_tiles",
	PermCanPromotePlayers:        "can_promote_players",
	PermCanModifyShowData:        "can_modify_show_data",
	PermCanManageTimers:          "can_manage_timers",
}

// PermissionDescriptions provides human-readable descriptions
var PermissionDescriptions = map[Permission]string{
	PermCanChat:                  "Can participate in chat",
	PermCanSendMessages:          "Can send chat messages",
	PermCanSendWhispers:          "Can send private messages",
	PermCanDeleteOwnMessages:     "Can delete own messages",
	PermCanDeleteMessages:        "Can delete any messages",
	PermCanBanUsers:              "Can ban users from chat",
	PermCanKickUsers:             "Can kick users from chat",
	PermCanMuteUsers:             "Can mute users",
	PermCanUnmuteUsers:           "Can unmute users",
	PermCanHost:                  "Can access host controls and streams",
	PermCanModerate:              "Can moderate chat and users",
	PermCanManageChat:            "Can manage chat settings",
	PermCanManageChatPermissions: "Can manage user chat permissions",
	PermCanSuggestTiles:          "Can suggest new bingo tiles",
	PermCanReviewTiles:           "Can review tile suggestions",
	PermCanApproveTiles:          "Can approve tiles for use",
	PermCanManageTiles:           "Can manage bingo tiles",
	PermCanPromotePlayers:        "Can promote other players",
	PermCanModifyShowData:        "Can modify show information",
	PermCanManageTimers:          "Can manage show timers",
}

// DefaultPermissions returns the default permission set for new users
func DefaultPermissions() Permission {
	return PermCanChat | PermCanSendMessages | PermCanSendWhispers | PermCanDeleteOwnMessages | PermCanSuggestTiles
}

// AdminPermissions returns full permissions for admin users
func AdminPermissions() Permission {
	// Give all permissions except the ones that might not be implemented yet
	return PermCanChat |
		PermCanSendMessages |
		PermCanSendWhispers |
		PermCanDeleteOwnMessages |
		PermCanDeleteMessages |
		PermCanBanUsers |
		PermCanKickUsers |
		PermCanMuteUsers |
		PermCanUnmuteUsers |
		PermCanHost |
		PermCanModerate |
		PermCanManageChat |
		PermCanManageChatPermissions |
		PermCanSuggestTiles |
		PermCanReviewTiles |
		PermCanApproveTiles |
		PermCanManageTiles
}

// HasPermission checks if a permission set contains a specific permission
func (p Permission) HasPermission(perm Permission) bool {
	return p&perm != 0
}

// AddPermission adds a permission to the set
func (p *Permission) AddPermission(perm Permission) {
	*p |= perm
}

// RemovePermission removes a permission from the set
func (p *Permission) RemovePermission(perm Permission) {
	*p &^= perm
}

// SetPermission sets a permission to a specific state
func (p *Permission) SetPermission(perm Permission, enabled bool) {
	if enabled {
		p.AddPermission(perm)
	} else {
		p.RemovePermission(perm)
	}
}

// GetPermissions returns a map of permission names to boolean values
func (p Permission) GetPermissions() map[string]bool {
	permissions := make(map[string]bool)
	for perm, name := range PermissionNames {
		permissions[name] = p.HasPermission(perm)
	}
	return permissions
}

// SetPermissionsFromMap sets permissions from a map of permission names to boolean values
func (p *Permission) SetPermissionsFromMap(permMap map[string]bool) {
	for name, enabled := range permMap {
		for perm, permName := range PermissionNames {
			if permName == name {
				p.SetPermission(perm, enabled)
			}
		}
	}
}

// GetAllPermissions returns all available permissions with their current state
func (p Permission) GetAllPermissions() map[string]interface{} {
	allPerms := make(map[string]interface{})
	for perm, name := range PermissionNames {
		allPerms[name] = map[string]interface{}{
			"enabled":     p.HasPermission(perm),
			"description": PermissionDescriptions[perm],
		}
	}
	return allPerms
}
