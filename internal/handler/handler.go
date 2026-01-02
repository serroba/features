package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/serroba/features/internal/flags"
)

type Handler struct {
	service *flags.Service
}

func New(service *flags.Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) CreateFlag(ctx context.Context, req *CreateFlagRequest) (*CreateFlagResponse, error) {
	flag := ToFlag(req.Body)

	err := h.service.Create(ctx, flag)
	if err != nil {
		if errors.Is(err, flags.ErrFlagExists) {
			return nil, huma.Error409Conflict("flag already exists")
		}

		return nil, huma.Error500InternalServerError("failed to create flag")
	}

	return &CreateFlagResponse{
		Body: CreateFlagResponseBody{
			Key:       flag.Key,
			Version:   flag.Version,
			CreatedAt: flag.UpdatedAt,
		},
	}, nil
}

func (h *Handler) EvaluateFlag(ctx context.Context, req *EvaluateFlagRequest) (*EvaluateFlagResponse, error) {
	evalCtx := ToEvalContext(req.Body)

	result, err := h.service.Evaluate(ctx, req.Key, evalCtx)
	if err != nil {
		if errors.Is(err, flags.ErrFlagNotFound) {
			return nil, huma.Error404NotFound("flag not found")
		}

		return nil, huma.Error500InternalServerError("failed to evaluate flag")
	}

	return &EvaluateFlagResponse{
		Body: ToEvalResultBody(result),
	}, nil
}
