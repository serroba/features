package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/serroba/features/internal/flags"
)

//go:generate mockgen -destination=mock_service_test.go -package=handler_test . FlagService

type FlagService interface {
	Create(ctx context.Context, flag flags.Flag) (flags.Flag, error)
	Evaluate(ctx context.Context, key flags.FlagKey, evalCtx flags.EvalContext) (flags.EvalResult, error)
}

type Handler struct {
	service FlagService
}

func New(service FlagService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateFlag(ctx context.Context, req *CreateFlagRequest) (*CreateFlagResponse, error) {
	flag, err := h.service.Create(ctx, ToFlag(req.Body))
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

	result, err := h.service.Evaluate(ctx, flags.FlagKey(req.Key), evalCtx)
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
