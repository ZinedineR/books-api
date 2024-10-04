package service

import (
	"books-api/internal/entity"
	"books-api/internal/repository"
	"books-api/pkg/signature"
	"context"
	"github.com/google/uuid"

	//"books-api/pkg/exception"
	"books-api/pkg/exception"
	"books-api/pkg/xvalidator"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	signaturer signature.Signaturer
	validate   *xvalidator.Validator
}

func NewUserService(
	db *gorm.DB, repo repository.UserRepository,
	signaturer signature.Signaturer,
	validate *xvalidator.Validator,
) UserService {
	return &UserServiceImpl{
		db:         db,
		userRepo:   repo,
		signaturer: signaturer,
		validate:   validate,
	}
}

func (s *UserServiceImpl) Register(
	ctx context.Context, model *entity.UserLogin,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	duplicateCheck, err := s.userRepo.FindByName(ctx, s.db, "username", model.Username)
	if err != nil {
		return exception.Internal("err", err)
	}
	if duplicateCheck != nil {
		return exception.PermissionDenied("username already exists")
	}
	password, err := s.signaturer.HashBscryptPassword(model.Password)
	if err != nil {
		return exception.Internal("can't create password", err)
	}
	body := &entity.User{
		Id:       uuid.NewString(),
		Username: model.Username,
		Password: password,
	}
	if err := s.userRepo.CreateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, model *entity.UserLogin) (
	*UserLoginResponse, *exception.Exception,
) {
	if errs := s.validate.Struct(model); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	result, err := s.userRepo.FindByName(ctx, s.db, "username", model.Username)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.NotFound("username not found")
	}
	if ok := s.signaturer.CheckBscryptPasswordHash(model.Password, result.Password); !ok {
		return nil, exception.PermissionDenied("username/password unmatched")
	}
	jwtToken, err := s.signaturer.GenerateJWT(result.Username)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return &UserLoginResponse{
		Username: result.Username,
		Token:    jwtToken,
	}, nil
}
