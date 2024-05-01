package services

import (
	Orchestrator "0xKowalski1/container-orchestrator/api"
	OrchestratorModels "0xKowalski1/container-orchestrator/models"

	"0xKowalski1/server-hosting-web/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameserverService struct {
	DB                  *gorm.DB
	OrchestratorWrapper *Orchestrator.WrapperClient
}

func NewGameserverService(db *gorm.DB) *GameserverService {
	return &GameserverService{
		DB:                  db,
		OrchestratorWrapper: Orchestrator.NewApiWrapper("development", "localhost"), // Get me from env
	}
}

func (service *GameserverService) CreateGameserver(newGameserver models.Gameserver) (*models.Gameserver, error) {
	result := service.DB.Create(&newGameserver)
	if result.Error != nil {
		log.Printf("Failed to create gameserver: %v", result.Error)
		return nil, result.Error
	}

	return &newGameserver, nil
}

func (service *GameserverService) GetGameservers() ([]models.Gameserver, error) {
	var gameservers []models.Gameserver

	result := service.DB.Preload("Game").Find(&gameservers)
	if result.Error != nil {
		return nil, result.Error
	}

	return gameservers, nil
}

func (service *GameserverService) GetGameserverByID(gameserverID string) (models.Gameserver, error) {
	var gameserver models.Gameserver

	id, err := uuid.Parse(gameserverID)
	if err != nil {
		return gameserver, err
	}

	result := service.DB.First(&gameserver, "id = ?", id)
	if result.Error != nil {
		return gameserver, result.Error
	}

	return gameserver, nil
}

// Orchestrator

func (service *GameserverService) DeployGameserver(gameserver models.Gameserver) error {
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
		// Do something
		log.Printf("Error deploying gameserver: %v", err)
		return err
	}

	return nil
}
