package service

import (
	"fmt"
	"myapp/config"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func GetService() *Service {
	s := Service{
		DB: config.GetDB(),
	}

	return &s
}

func GetTransaction() *Service {
	fmt.Println("begin...")
	s := Service{
		DB: config.GetDB().Begin(),
	}

	return &s
}

func (s *Service) Commit() {
	if err := s.DB.Commit().Error; err != nil {
		panic(err)
	}

	fmt.Println("commit...")
}

func (s *Service) Rollback(err ...interface{}) {
	s.DB.Rollback()
	fmt.Println("rollback...")

	if len(err) > 0 && err[0] != nil {
		panic(err[0])
	}
}
