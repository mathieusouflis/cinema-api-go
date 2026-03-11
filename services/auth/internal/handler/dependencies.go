package handler

import (
	loginUsecase "authService/internal/usecase/login"
	registerUsecase "authService/internal/usecase/register"
)

type Dependencies struct {
	LoginUseCase    loginUsecase.Usecase
	RegisterUseCase registerUsecase.Usecase
}
