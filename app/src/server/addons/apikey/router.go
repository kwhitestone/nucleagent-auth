package apikey

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

var svc = &Service{}

// ctxKey 是 user_id 在 request context 中的键类型（避免字符串键冲突）。
type ctxKey int

const userIDKey ctxKey = 1

// BridgeMiddleware 把 JWT 中间件写入 gin.Context 的 user_id 复制到 request context，
// 供 huma handler 通过 ctx 读取（huma handler 拿到的是 request.Context()）。
func BridgeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("user_id")
		if exists {
			uid, _ := v.(uint)
			ctx := context.WithValue(c.Request.Context(), userIDKey, uid)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

// CreateInput 创建 API Key 请求体。
type CreateInput struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" doc:"Key用途说明"`
	}
}

// CreateOutput 创建 API Key 响应体（含明文，仅此一次）。
type CreateOutput struct {
	Body struct {
		Code    int         `json:"code" example:"0"`
		Message string      `json:"message" example:"success"`
		Data    *APIKeyView `json:"data"`
	}
}

// ListOutput API Key 列表响应。
type ListOutput struct {
	Body struct {
		Code    int          `json:"code" example:"0"`
		Message string       `json:"message" example:"success"`
		Data    []APIKeyView `json:"data"`
	}
}

// DeleteOutput 删除响应。
type DeleteOutput struct {
	Body struct {
		Code    int    `json:"code" example:"0"`
		Message string `json:"message" example:"success"`
	}
}

// RegisterRoutes 注册 API Key 路由到 Huma。
//
// 路由前缀由 plugin RoutePrefix 决定（/api/v1/addons/auth），
// JWT 中间件已强制认证并写入 user_id 到 gin context，
// BridgeMiddleware 再把它桥接到 request context。
func RegisterRoutes(api huma.API) {
	// 创建 API Key
	huma.Register(api, huma.Operation{
		OperationID: "createAPIKey",
		Method:      http.MethodPost,
		Path:        "/api/v1/addons/auth/api-keys",
		Summary:     "创建 API Key",
		Description: "为当前用户创建一个 API Key，明文仅返回一次",
		Tags:        []string{"APIKey"},
		Security: []map[string][]string{
			{"AuthTokenAuth": {}},
		},
	}, func(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
		userID := userIDFromCtx(ctx)
		if userID == 0 {
			return nil, huma.NewError(http.StatusUnauthorized, "未认证")
		}
		view, err := svc.Create(userID, input.Body.Name)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, err.Error())
		}
		resp := &CreateOutput{}
		resp.Body.Code = 0
		resp.Body.Message = "创建成功"
		resp.Body.Data = view
		return resp, nil
	})

	// 列出 API Key
	huma.Register(api, huma.Operation{
		OperationID: "listAPIKeys",
		Method:      http.MethodGet,
		Path:        "/api/v1/addons/auth/api-keys",
		Summary:     "列出 API Key",
		Description: "列出当前用户的所有 API Key（不含明文）",
		Tags:        []string{"APIKey"},
		Security: []map[string][]string{
			{"AuthTokenAuth": {}},
		},
	}, func(ctx context.Context, input *struct{}) (*ListOutput, error) {
		userID := userIDFromCtx(ctx)
		if userID == 0 {
			return nil, huma.NewError(http.StatusUnauthorized, "未认证")
		}
		keys, err := svc.List(userID)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, err.Error())
		}
		resp := &ListOutput{}
		resp.Body.Code = 0
		resp.Body.Message = "success"
		resp.Body.Data = keys
		return resp, nil
	})

	// 删除 API Key
	huma.Register(api, huma.Operation{
		OperationID: "deleteAPIKey",
		Method:      http.MethodDelete,
		Path:        "/api/v1/addons/auth/api-keys/{id}",
		Summary:     "删除 API Key",
		Description: "删除当前用户的指定 API Key",
		Tags:        []string{"APIKey"},
		Security: []map[string][]string{
			{"AuthTokenAuth": {}},
		},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"API Key ID"`
	}) (*DeleteOutput, error) {
		userID := userIDFromCtx(ctx)
		if userID == 0 {
			return nil, huma.NewError(http.StatusUnauthorized, "未认证")
		}
		id, err := strconv.ParseUint(input.ID, 10, 64)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "无效的 ID")
		}
		if err := svc.Delete(userID, uint(id)); err != nil {
			return nil, huma.NewError(http.StatusNotFound, err.Error())
		}
		resp := &DeleteOutput{}
		resp.Body.Code = 0
		resp.Body.Message = "删除成功"
		return resp, nil
	})
}

// userIDFromCtx 从 request context 提取 user_id（由 BridgeMiddleware 写入）。
func userIDFromCtx(ctx context.Context) uint {
	v, _ := ctx.Value(userIDKey).(uint)
	return v
}
