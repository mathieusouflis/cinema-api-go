package handler

import (
	loginUsecase "authService/internal/usecase/login"
	logoutUsecase "authService/internal/usecase/logout"
	oauthCallbackUsecase "authService/internal/usecase/oauth"
	refreshUsecase "authService/internal/usecase/refresh"
	registerUsecase "authService/internal/usecase/register"
)

type Dependencies struct {
	LoginUseCase         loginUsecase.Usecase
	RegisterUseCase      registerUsecase.Usecase
	RefreshUseCase       refreshUsecase.Usecase
	LogoutUseCase        logoutUsecase.Usecase
	OauthCallbackUseCase oauthCallbackUsecase.Usecase
}
