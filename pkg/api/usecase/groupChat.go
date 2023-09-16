package usecase

import "chat/pkg/api/repository"

type GroupChatUsecase struct {
	GroupChatRepo repository.GroupChatRepoMethods
}

func NewGroupChatUsecase(repo repository.GroupChatRepoMethods) GroupChatUsecase {
	return GroupChatUsecase{
		GroupChatRepo: repo,
	}
}

type GroupChatUsecaseMethods interface {
	GroupChatStart()
}

func (u GroupChatUsecase) GroupChatStart() {
	u.GroupChatRepo.CreateGroupChat()
}
