package services

import (
	Orchestrator "0xKowalski1/container-orchestrator/api"
	OrchestratorModels "0xKowalski1/container-orchestrator/models"
	"fmt"

	"0xKowalski1/server-hosting-web/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameserverService struct {
	DB                  *gorm.DB
	OrchestratorWrapper *Orchestrator.WrapperClient
}

func NewGameserverService(db *gorm.DB, orchestratorWrapper *Orchestrator.WrapperClient) *GameserverService {
	return &GameserverService{
		DB:                  db,
		OrchestratorWrapper: orchestratorWrapper,
	}
}

func (service *GameserverService) CreateGameserver(newGameserver models.Gameserver, stripeService *StripeService) (*models.Gameserver, error) {
	err := stripeService.CanAllocateResources(newGameserver.UserID, newGameserver.MemoryLimit, newGameserver.StorageLimit, service)
	if err != nil {
		return nil, fmt.Errorf("Failed to create gameserver as requested resources would exceed subscription limits: %v", err)
	}

	result := service.DB.Create(&newGameserver)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to create gameserver with options: %+v due to error: %v", newGameserver, result.Error)
	}

	return &newGameserver, nil
}

func (service *GameserverService) GetGameservers(userID string) ([]models.Gameserver, error) {
	var gameservers []models.Gameserver

	result := service.DB.Preload("Game").Find(&gameservers, "user_id = ?", userID)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to query database for gameservers: %v", result.Error)
	}

	return gameservers, nil
}

func (service *GameserverService) GetUsedResources(gameservers []models.Gameserver) (int, int) {
	var usedMemory, usedStorage int
	for _, gameserver := range gameservers {
		usedMemory += gameserver.MemoryLimit
		usedStorage += gameserver.StorageLimit
	}

	return usedMemory, usedStorage
}

func (service *GameserverService) GetGameserverByID(gameserverID string) (*models.Gameserver, error) {
	var gameserver *models.Gameserver

	id, err := uuid.Parse(gameserverID)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse gameserver ID: %s due to error: %v", gameserverID, err)
	}

	result := service.DB.First(&gameserver, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to query database for gameserver at ID: %s due to error: %v", gameserverID, result.Error)
	}

	return gameserver, nil
}

// Orchestrator

func (service *GameserverService) DeployGameserver(gameserver *models.Gameserver) error {
	newContainerRequest := OrchestratorModels.CreateContainerRequest{
		ID:           gameserver.ID.String(),
		Image:        gameserver.Game.ContainerImage,
		Env:          []string{"EULA=TRUE", "MEMORY=4"},
		StopTimeout:  5,
		MemoryLimit:  gameserver.MemoryLimit,
		CpuLimit:     2,
		StorageLimit: gameserver.StorageLimit,
		Ports: []OrchestratorModels.Port{
			{
				HostPort:      30001,
				ContainerPort: 25565,
				Protocol:      "tcp",
			},
		},
	}

	_, err := service.OrchestratorWrapper.CreateContainer(newContainerRequest)
	if err != nil {
		return fmt.Errorf("Error deploying gameserver: %v", err)
	}

	return nil
}
