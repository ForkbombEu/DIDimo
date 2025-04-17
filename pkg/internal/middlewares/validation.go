// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package middlewares

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

var validate = validator.New()

func ValidateInputMiddleware[T any]() func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var input T

		if err := json.NewDecoder(e.Request.Body).Decode(&input); err != nil {
			return apis.NewBadRequestError("Invalid JSON body", err)
		}

		if err := validate.Struct(input); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				details := make([]map[string]interface{}, 0, len(validationErrors))
				for _, ve := range validationErrors {
					details = append(details, map[string]interface{}{
						"field":   ve.Field(),
						"tag":     ve.Tag(),
						"param":   ve.Param(),
						"value":   ve.Value(),
						"message": ve.Error(),
					})
				}
				return apis.NewBadRequestError("Validation failed", map[string]interface{}{
					"errors": details,
				})
			}
			return apis.NewBadRequestError("Validation failed", err)
		}
		ctx := e.Request.Context()
		ctx = context.WithValue(ctx, "validatedInput", input)
		e.Request = e.Request.WithContext(ctx)
		e.Request = e.Request.WithContext(ctx)

		return e.Next()
	}
}
