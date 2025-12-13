package usecase

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/config"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/database"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/logger"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/repository"
	"github.com/fazriegi/fintrack-be/internal/pkg"

	"github.com/sirupsen/logrus"
)

type AuthUsecase interface {
	Register(props *entity.RegisterRequest) (resp pkg.Response)
	Login(props *entity.LoginRequest) (resp pkg.Response)
}

type authUsecase struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
	log      *logrus.Logger
	jwt      *pkg.JWT
}

func NewAuthUsecase(userRepo repository.UserRepository, authRepo repository.AuthRepository, jwt *pkg.JWT) AuthUsecase {
	log := logger.Get()

	return &authUsecase{
		userRepo,
		authRepo,
		log,
		jwt,
	}
}

func (u *authUsecase) Register(props *entity.RegisterRequest) (resp pkg.Response) {
	var (
		err            error
		user           entity.User
		hashedPassword string
		db             = database.Get()
	)

	user, err = u.userRepo.GetByUsername(props.Username, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.log.Errorf("userRepo.GetByUsername: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	if user.Username != "" {
		return pkg.NewResponse(http.StatusBadRequest, "username already exists", nil, nil)
	}

	if hashedPassword, err = pkg.HashPassword(props.Password); err != nil {
		u.log.Errorf("pkg.HashPassword: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error start transaction: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}
	defer tx.Rollback()

	user = entity.User{
		Name:     props.Name,
		Email:    props.Email,
		Password: hashedPassword,
		Username: props.Username,
	}

	_, err = u.userRepo.Insert(&user, tx)
	if err != nil {
		u.log.Errorf("userRepo.Insert: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	createdUser := entity.UserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}

	return pkg.NewResponse(http.StatusCreated, "success", createdUser, nil)
}

func (u *authUsecase) Login(props *entity.LoginRequest) (resp pkg.Response) {
	db := database.Get()

	existingUser, err := u.userRepo.GetByUsername(props.Username, db)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return pkg.NewResponse(http.StatusUnauthorized, "invalid username or password", nil, nil)
	} else if err != nil {
		u.log.Errorf("userRepo.GetByUsername: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	if !pkg.CheckPasswordHash(props.Password, existingUser.Password) {
		return pkg.NewResponse(http.StatusUnauthorized, "invalid username or password", nil, nil)
	}

	accessToken, err := u.jwt.GenerateAccessToken(existingUser.ID, existingUser.Email, existingUser.Username)
	if err != nil {
		u.log.Errorf("u.jwt.GenerateAccessToken: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	refreshToken, err := u.jwt.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		u.log.Errorf("u.jwt.GenerateRefreshToken: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error start transaction: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}
	defer tx.Rollback()

	refreshTokenData := map[string]any{
		"user_id":    existingUser.ID,
		"token":      pkg.Hash(refreshToken),
		"expired_at": time.Now().Add(time.Duration(config.GetUint("jwt.refreshToken.expDay")) * 24 * time.Hour),
	}

	_, err = u.authRepo.InsertRefreshToken(refreshTokenData, tx)
	if err != nil {
		u.log.Errorf("authRepo.InsertRefreshToken: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		return pkg.NewResponse(http.StatusInternalServerError, pkg.ErrServer.Error(), nil, nil)
	}

	data := map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": entity.UserResponse{
			Name:     existingUser.Name,
			Username: existingUser.Username,
			Email:    existingUser.Email,
		},
	}

	return pkg.NewResponse(http.StatusOK, "success", data, nil)
}
