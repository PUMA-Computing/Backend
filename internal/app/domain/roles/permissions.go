package permissions

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//var Permissions = map[string]Permission{
//	PermissionManageAll:                  {1, "Manage All", "Allow user to manage all"},
//	PermissionManageProfile:              {2, "Manage Profile", "Allow user to manage their profile"},
//	PermissionEditProfile:                {3, "Edit Profile", "Allow user to edit their profile"},
//	PermissionChangePassword:             {4, "Change Password", "Allow user to change their password"},
//	PermissionDeleteAccount:              {5, "Delete Account", "Allow user to delete their account"},
//	PermissionRegisterEvent:              {6, "Register Event", "Allow user to register event"},
//	PermissionManageUser:                 {7, "Manage User", "Allow user to manage user"},
//	PermissionCreateUser:                 {8, "Create User", "Allow user to create user"},
//	PermissionEditUser:                   {9, "Edit User", "Allow user to edit user"},
//	PermissionDeleteUser:                 {10, "Delete User", "Allow user to delete user"},
//	PermissionManageEvent:                {11, "Manage Event", "Allow user to manage event"},
//	PermissionCreateEvent:                {12, "Create Event", "Allow user to create event"},
//	PermissionEditEvent:                  {13, "Edit Event", "Allow user to edit event"},
//	PermissionDeleteEvent:                {14, "Delete Event", "Allow user to delete event"},
//	PermissionGetUsersRegisteredEvent:    {15, "Get Users Registered Event", "Allow user to get users registered event"},
//	PermissionDeleteUsersRegisteredEvent: {16, "Delete Users Registered Event", "Allow user to delete users registered event"},
//	PermissionManageNews:                 {17, "Manage News", "Allow user to manage news"},
//	PermissionCreateNews:                 {18, "Create News", "Allow user to create news"},
//	PermissionEditNews:                   {19, "Edit News", "Allow user to edit news"},
//	PermissionDeleteNews:                 {20, "Delete News", "Allow user to delete news"},
//	PermissionManageRoleAndPermissions:   {21, "Manage Role and Permissions", "Allow user to manage role and permissions"},
//	PermissionCreateRole:                 {22, "Create Role", "Allow user to create role"},
//	PermissionEditRole:                   {23, "Edit Role", "Allow user to edit role"},
//	PermissionDeleteRole:                 {24, "Delete Role", "Allow user to delete role"},
//	PermissionAssignRolePermission:       {25, "Assign Role Permission", "Allow user to assign role permission"},
//	PermissionAssignUserRole:             {26, "Assign User Role", "Allow user to assign user role"},
//}
