package routes

const (
	AuthCallbackURL  = "/auth/{" + MuxIDParam + "}/callback"
	AuthLoginURL     = "/auth/{" + MuxIDParam + "}/login"
	CreateCommentURL = "/groups/{" + MuxGroupParam + ":[0-9]+}/notes/{" + MuxIDParam + ":[0-9]+}/comments/create"
	CreateGoalURL    = "/groups/{" + MuxGroupParam + ":[0-9]+}/meetings/{" + MuxIDParam + ":[0-9]+}/goals/create"
	CreateGroupURL   = "/groups/create"
	CreateMeetingURL = "/groups/{" + MuxIDParam + ":[0-9]+}/meetings/create"
	CreateNoteURL    = "/groups/{" + MuxIDParam + ":[0-9]+}/notes/create"
	GetGroupsURL     = "/groups"
	GetGroupURL      = "/groups/{" + MuxIDParam + ":[0-9]+}"
	GetLoginURL      = "/login"
	GetLogoffURL     = "/logoff"
	GetMeetingURL    = "/groups/{" + MuxGroupParam + ":[0-9]+}/meetings/{" + MuxIDParam + ":[0-9]+}"
	GetNoteURL       = "/groups/{" + MuxGroupParam + ":[0-9]+}/notes/{" + MuxIDParam + ":[0-9]+}"
	GetNotesURL      = "/groups/{" + MuxIDParam + ":[0-9]+}/notes"
	GetRootURL       = "/"
	GetSSEURL        = "/sse"
	GetWebsocketURL  = "/ws"
	NewGroupURL      = "/groups/new"
	NewNoteURL       = "/groups/{" + MuxIDParam + ":[0-9]+}/notes/new"
	UpdateMeetingURL = GetMeetingURL + "/update"

	// Router variable helper strings
	MuxIDParam    = "id"
	MuxGroupParam = "groupID"
)
