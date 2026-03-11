package handler

import usecase "authService/internal/usecase/login"

type Dependencies struct {
	LoginUseCase *usecase.Usecase
}
