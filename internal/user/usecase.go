package user

import (
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/pkg/utils"
)

type UseCase struct {
	IRepository
}

func (u UseCase) GetName(id string) string {
	uintID, _ := utils.ConvertStrToUInt(id)

	return u.IRepository.GetName(uintID)
}

func (u UseCase) GetByID(id uint) (domain.User, error) {
	return u.IRepository.GetByID(id)
}

func NewUseCase(repo IRepository) IUseCase {
	return &UseCase{IRepository: repo}
}
