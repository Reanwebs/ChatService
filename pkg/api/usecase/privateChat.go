package usecase

import (
	"chat/pkg/api/repository"
)

type PrivateChatUsecase struct {
	PrivateChatRepo repository.PrivateChatRepoMethods
}

func NewPrivateChatUsecase(repo repository.PrivateChatRepoMethods) PrivateChatUsecase {
	return PrivateChatUsecase{
		PrivateChatRepo: repo,
	}
}

type PrivateChatUsecaseMethods interface {
	PrivateChatStart()
}

func (r PrivateChatUsecase) PrivateChatStart() {

}
