package database

import (
	"context"
	"fmt"
	"time"

	"github.com/wizenheimer/cascade/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DatabaseClient is an interface for a database client
type DatabaseClient interface {
	Ready() bool

	// Scenario related methods
	// CRUD Methods for Scenario
	CreateScenario(ctx context.Context, scenario *models.Scenario) (*models.Scenario, error)
	UpdateScenario(ctx context.Context, scenarioID string, updatedScenario *models.Scenario) (*models.Scenario, error)
	DeleteScenario(ctx context.Context, scenarioID string) error
	GetScenarioByID(ctx context.Context, scenarioID string) ([]models.Scenario, error)
	// Listing Methods for Scenarios
	ListScenarios(ctx context.Context) ([]models.Scenario, error)
	ListScenariosByTeamID(ctx context.Context, teamID string) ([]models.Scenario, error)
	ListScenarioVersion(ctx context.Context, scenarioID string) ([]models.Scenario, error)

	// Session related methods
	CreateSession(ctx context.Context, scenarioID string, version int) (*models.Session, error)
	StartSession(ctx context.Context, sessionID string) (*models.Session, error)
	GracefullyEndSession(ctx context.Context, sessionID string) (*models.Session, error)
	TerminateSession(ctx context.Context, sessionID string) (*models.Session, error)
	// Listing Method for Sessions
	ListSessionByScenarioID(ctx context.Context, scenarioID string, version int) ([]models.Session, error)

	// Metrics related methods
	GetSessionMetrics(ctx context.Context, scenarioID string) ([]models.SessionMetrics, error)

	// UserManagement related methods
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, email string, updatedUser *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, email string) (*models.User, error)
	DeactivateUser(ctx context.Context, email string) (*models.User, error)

	// TeamManagement related methods
	CreateTeam(ctx context.Context, name string, description string, creator *models.User) (*models.Team, error)
	GetTeamByID(ctx context.Context, teamID string) (*models.Team, error)
	AddUserToTeam(ctx context.Context, user *models.User, team *models.Team) (*models.User, error)
	RemoveUserFromTeam(ctx context.Context, user *models.User, team *models.Team) (*models.User, error)
	ListUsersByTeam(ctx context.Context, teamID string) ([]models.User, error)
	UpdateTeambyTeamID(ctx context.Context, teamID string, updatedTeam *models.Team) (*models.Team, error)
	DeleteTeam(ctx context.Context, teamID string) (*models.Team, error)
	DeactivateTeam(ctx context.Context, teamID string) (*models.Team, error)
}

// Client is a database client
type Client struct {
	DB *gorm.DB
}

// NewDatabaseClient creates a new database client
func NewDatabaseClient(host string, user string, password string, dbname string, port int32, sslmode string) (DatabaseClient, error) {
	// Prepare the data source name (DSN) for the PostgreSQL connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	// Open a new connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "cascade.",
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	})
	if err != nil {
		return nil, err
	}

	// Return a new database client
	client := Client{
		DB: db,
	}
	return client, nil
}

// Ready checks if the database connection is ready
func (c Client) Ready() bool {
	var ready string
	tx := c.DB.Raw("SELECT 1 as ready").Scan(&ready)
	if tx.Error != nil {
		return false
	}
	if ready == "1" {
		return true
	}
	return false
}
