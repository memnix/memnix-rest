package user

import (
	"github.com/memnix/memnix-rest/pkg/utils"
)

type UseCase struct {
	IRepository
}

func (u UseCase) GetName(id string) string {
	uintID, _ := utils.ConvertStrToUInt(id)

	return u.IRepository.GetName(uintID)
}

func NewUseCase(repo IRepository) IUseCase {
	return &UseCase{IRepository: repo}
}
