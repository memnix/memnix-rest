package user

import "github.com/edgedb/edgedb-go"

type UseCase struct {
	IEdgeRepository
}

func (u UseCase) GetName(id edgedb.UUID) string {
	return u.IEdgeRepository.GetName(id)
}

func NewUseCase(edgeRepo IEdgeRepository) IUseCase {
	return &UseCase{IEdgeRepository: edgeRepo}
}
