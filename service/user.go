package service

import (
	"context"
	"myapp/graph/model"
	"myapp/middleware"
	"myapp/tools"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func (s *Service) UserRegister(ctx context.Context, input model.NewUser) (string, error) {
	if input.Email == "" || input.Password == "" {
		panic(gqlerror.Errorf("invalid email/password input"))
	}

	exist, _ := s.UserCheckExistByEmail(ctx, input.Email)
	if exist {
		panic(gqlerror.Errorf("email already used"))
	}

	s.UserCreate(ctx, input)

	return "Success", nil
}

func (s *Service) UserLogin(ctx context.Context, input model.LoginInput) (*model.LoginResponse, error) {
	if input.Email == "" || input.Password == "" {
		panic(gqlerror.Errorf("invalid email/password input"))
	}

	user, _ := s.UserGetByEmail(ctx, input.Email)
	valid, err := tools.CompareHash(user.Password, input.Password)
	if err != nil {
		panic(err)
	}

	if !valid {
		panic(gqlerror.Errorf("user not found"))
	}

	token := tools.TokenCreate(user.ID)

	return &model.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) UserCreate(ctx context.Context, input model.NewUser) (*model.User, error) {
	password, err := tools.HashAndSalt(input.Password)
	if err != nil {
		panic(err)
	}

	user := model.User{
		Email:     input.Email,
		Password:  password,
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.DB.Model(&user).Omit("updated_at").Create(&user).Error; err != nil {
		panic(err)
	}

	return &user, nil
}

func (s *Service) UserGetMe(ctx context.Context) (*model.User, error) {
	var (
		getUser = middleware.AuthContext(ctx)
	)

	return s.UserGetByID(ctx, getUser.ID)
}

func (s *Service) UserGetByID(ctx context.Context, id int) (*model.User, error) {
	var (
		user model.User
	)

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("id = ?", id).First(&user).Error; err == gorm.ErrRecordNotFound {
		panic(gqlerror.Errorf("user not found"))
	} else if err != nil {
		panic(err)
	}

	return &user, nil
}

func (s *Service) UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	var (
		user model.User
	)

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("email = ?", email).First(&user).Error; err == gorm.ErrRecordNotFound {
		panic(gqlerror.Errorf("user not found"))
	} else if err != nil {
		panic(err)
	}

	return &user, nil
}

func (s *Service) UserCheckExistByID(ctx context.Context, id int) (bool, error) {
	var (
		user  model.User
		exist bool = false
		count int64
	)

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("id = ?", id).Count(&count).Error; err != nil {
		panic(err)
	}

	if int(count) > 0 {
		exist = true
	}

	return exist, nil
}

func (s *Service) UserCheckExistByEmail(ctx context.Context, email string) (bool, error) {
	var (
		user  model.User
		exist bool = false
		count int64
	)

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("email = ?", email).Count(&count).Error; err != nil {
		panic(err)
	}

	if int(count) > 0 {
		exist = true
	}

	return exist, nil
}
