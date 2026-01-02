package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-flag",
		Method:      http.MethodPost,
		Path:        "/flags",
		Summary:     "Create a new feature flag",
		Tags:        []string{"Flags"},
	}, h.CreateFlag)

	huma.Register(api, huma.Operation{
		OperationID: "evaluate-flag",
		Method:      http.MethodPost,
		Path:        "/flags/{key}/evaluate",
		Summary:     "Evaluate a feature flag",
		Tags:        []string{"Flags"},
	}, h.EvaluateFlag)
}
