package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"eigenflux_server/api/clients"
	auth "eigenflux_server/kitex_gen/eigenflux/auth"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		header := string(c.GetHeader("Authorization"))
		if header == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "missing or invalid authorization header",
			})
			c.Abort()
			return
		}
		accessToken := strings.TrimPrefix(header, "Bearer ")

		resp, err := clients.AuthClient.ValidateSession(ctx, &auth.ValidateSessionReq{
			AccessToken: accessToken,
		})
		if err != nil || resp.BaseResp.Code != 0 {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("agent_id", resp.AgentId)
		c.Next(ctx)
	}
}
