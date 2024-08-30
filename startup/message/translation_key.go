package message

var (
	ErrPlayerNotFound       = translationKey{"player.not_found", "player"}            // This means the target player was not found
	ErrTeamNotFound         = translationKey{"team.not_found", "team"}                // This means the target team was not found
	ErrTeamAlreadyExists    = translationKey{"team.already_exists", "team"}           // This means a team with the same name already exists
	ErrPlayerAlreadyInTeam  = translationKey{"team.player_already_in_team", "player"} // This means the target player is already in a team
	ErrSelfAlreadyInTeam    = translationKey{"team.self_already_in_team"}             // This means the sender is already in a team
	ErrPlayerNotInTeam      = translationKey{"team.player_not_in_team", "player"}     // This means the target player is not in a team
	ErrPlayerNotTeamMember  = translationKey{"team.player_not_team_member", "player"} // This means the target player is not a member of the team
	ErrPlayerAlreadyMember  = translationKey{"team.player_already_member", "player"}  // This means the target player is already a member of the team
	ErrPlayerAlreadyInvited = translationKey{"team.player_already_invited", "player"} // This means the target player is already invited to the team
	ErrPlayerHighestRole    = translationKey{"team.player_highest_role"}              // This means the target player has the highest role in the team
	ErrSelfNotInTeam        = translationKey{"team.self_not_in_team"}                 // This means the sender is not in a team
	ErrSelfNotLeader        = translationKey{"team.self_not_leader"}                  // This means the sender is not the leader of the team
	ErrSelfNotOfficer       = translationKey{"team.self_not_officer"}                 // This means the sender is not an officer of the team
	ErrSelfNotInvited       = translationKey{"team.self_not_invited", "team"}         // This means the sender is not invited to the team
	ErrCannotUseOnSelf      = translationKey{"team.cannot_use_on_self"}               // This means the sender cannot use the command on themselves

	SuccessTeamCreated     = translationKey{"team.success_broadcast_team_created", "player", "team"} // This means a team was successfully created
	SuccessSelfTeamCreated = translationKey{"team.success_self_team_created", "team"}                // This means the sender successfully created a team

	SuccessTeamInviteSent          = translationKey{"team.success_team_invite_sent", "player"}                     // This means the sender successfully sent an invitation to the target player
	SuccessBroadcastTeamInviteSent = translationKey{"team.success_broadcast_team_invite_sent", "sender", "player"} // This means the sender successfully sent an invitation to the target player
	SuccessTeamInviteReceived      = translationKey{"team.success_team_invite_received", "sender", "team"}         // This means the target player successfully received an invitation

	SuccessSelfTeamDisband = translationKey{"team.success_self_team_disband", "team"}      // This means the sender successfully disbanded their team
	SuccessTeamDisband     = translationKey{"team.success_team_disband", "player", "team"} // This means a team was successfully disbanded

	SuccessTeamMemberLeft = translationKey{"team.success_team_member_left", "player"} // This means a player successfully left the team
	SuccessSelfLeftTeam   = translationKey{"team.success_self_left_team", "team"}     // This means the sender successfully left the team

	SuccessSelfTeamMemberKicked = translationKey{"team.success_self_team_member_kicked", "player"} // This means the sender successfully kicked the target player from the team
	SuccessTeamKick             = translationKey{"team.success_team_kick", "player", "sender"}     // This means the target player was successfully kicked from the team
	SuccessSelfTeamKicked       = translationKey{"team.success_self_team_kicked", "team"}          // This means the target player was successfully kicked from the team
)

type translationKey []string

func (t translationKey) Build(args ...string) string {
	return ""
}
