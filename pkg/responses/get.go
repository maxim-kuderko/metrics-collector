package responses

import "context"

type Get struct {
	context.Context `json:"-"`
	BaseResponse    `json:"-"`
}
