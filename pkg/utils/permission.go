package utils

//type PermissionService interface {
//	CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
//}
//
//type Utils struct {
//	PermissionService PermissionService
//}
//
//func NewUtils(permissionService PermissionService) *Utils {
//	return &Utils{
//		PermissionService: permissionService,
//	}
//}
//
//func (u *Utils) HasPermission(c *gin.Context, permission string) (bool, error) {
//	userID, err := GetUserIDFromContext(c)
//	if err != nil {
//		return false, err
//	}
//
//	hasPermission, err := u.PermissionService.CheckPermission(context.Background(), userID, permission)
//	if err != nil {
//		return false, err
//	}
//
//	return hasPermission, nil
//}
