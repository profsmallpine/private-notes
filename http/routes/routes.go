package routes

const (
	AuthCallbackURL  = "/auth/{" + MuxIDParam + "}/callback"
	AuthLoginURL     = "/auth/{" + MuxIDParam + "}/login"
	CreateCommentURL = "/groups/{" + MuxGroupParam + ":[0-9]+}/notes/{" + MuxIDParam + ":[0-9]+}/comments/create"
	CreateGroupURL   = "/groups/create"
	CreateNoteURL    = "/groups/{" + MuxIDParam + ":[0-9]+}/notes/create"
	GetGroupsURL     = "/groups"
	GetGroupURL      = "/groups/{" + MuxIDParam + ":[0-9]+}"
	GetLoginURL      = "/login"
	GetLogoffURL     = "/logoff"
	GetNoteURL       = "/groups/{" + MuxGroupParam + ":[0-9]+}/notes/{" + MuxIDParam + ":[0-9]+}"
	GetNotesURL      = "/groups/{" + MuxIDParam + ":[0-9]+}/notes"
	GetRootURL       = "/"
	NewGroupURL      = "/groups/new"
	NewNoteURL       = "/groups/{" + MuxIDParam + ":[0-9]+}/notes/new"

	// Router variable helper strings
	MuxIDParam    = "id"
	MuxGroupParam = "groupID"
)
