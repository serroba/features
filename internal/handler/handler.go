package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/serroba/features/internal/flags"
)

type Handler struct {
	service *flags.Service
}

func New(service *flags.Service) *Handler {
	return &Handler{service: service}
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
