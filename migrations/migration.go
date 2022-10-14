package migrations

import "github.com/xy-planning-network/trails/postgres"

// List is used to run migration against the archive application database.
var List = []postgres.Migration{
	{Executor: CreateGroupsTable, Key: "20201109_create_groups"},
	{Executor: CreateUsersTable, Key: "20201109_create_users"},
	{Executor: CreateNotesTable, Key: "20201110_create_notes"},
	{Executor: CreateCommentsTable, Key: "20201113_create_comments"},
	{Executor: AddIsAdminToUsers, Key: "20211022_add_is_admin_to_users"},
	{Executor: CreateMeetingsTable, Key: "20211022_create_meetings"},
	{Executor: CreateGoalsTable, Key: "20211022_create_goals"},
	{Executor: AddStyleToGoals, Key: "20221007_add_style_to_goals"},
}
