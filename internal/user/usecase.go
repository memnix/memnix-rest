package user

import "github.com/edgedb/edgedb-go"

type UseCase struct {
	IEdgeRepository
}

func (u UseCase) GetName(id string) string {
	// Convert uuid to edgedb.UUID
	uuidEdge, err := edgedb.ParseUUID(id)
	if err != nil {
		return ""
	}

	return u.IEdgeRepository.GetName(uuidEdge)
}

func NewUseCase(edgeRepo IEdgeRepository) IUseCase {
	return &UseCase{IEdgeRepository: edgeRepo}
}
