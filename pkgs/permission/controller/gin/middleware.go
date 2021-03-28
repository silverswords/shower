/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package controller

import (
	"errors"
	"net/http"

	permission "github.com/abserari/shower/permission/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	errPermission = errors.New("admin permission is wrong")
)

//CheckPermission middleware that checks the permission
func (c *Controller) CheckPermission() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		reqURL := ctx.Request.URL.Path

		adminID, err := c.getIDFunc(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		// not check to admin return
		if adminID == 1000 {
			return
		}

		adRole, err := permission.AdminGetRoleMap(c.db, adminID)
		if err != nil {
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		urlRole, err := permission.URLPermissions(c.db, &reqURL)
		if err != nil {
			ctx.AbortWithError(http.StatusFailedDependency, err)
			return
		}

		for urlkey := range urlRole {
			for adkey := range adRole {
				if urlkey == adkey {
					return
				}
			}
		}

		ctx.AbortWithError(http.StatusForbidden, errPermission)
		ctx.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})

	}
}
