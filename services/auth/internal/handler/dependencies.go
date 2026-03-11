package handler

import loginUsecase "authService/internal/usecase/login"

type Dependencies struct {
	LoginUseCase loginUsecase.Usecase
}
