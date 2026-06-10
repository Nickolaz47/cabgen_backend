package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewTicketRepo(t *testing.T) {
	db := testutils.NewMockDB()
	ticketRepo := repositories.NewTicketRepo(db)

	assert.NotEmpty(t, ticketRepo)
}

func TestGetTickets(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	db := testutils.NewMockDB()
	ticketRepo := repositories.NewTicketRepo(db)

	ticket1 := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)
	ticket2 := testmodels.NewTicket(
		uuid.NewString(),
		"Maria",
		"maria@mail.com",
		"INCA",
		"Wrong username",
		"Cannot access my account.",
		&admin,
	)
	ticket2.Status = models.TicketStatusResolved
	db.Create(&ticket1)
	db.Create(&ticket2)

	t.Run("Success - Without Filter", func(t *testing.T) {
		tickets, err := ticketRepo.GetTickets(ctx, "")

		assert.NoError(t, err)
		assert.Len(t, tickets, 2)
		assert.Equal(t, ticket1.ID, tickets[0].ID)
		assert.Equal(t, ticket1.Status, tickets[0].Status)
		assert.Equal(t, ticket2.ID, tickets[1].ID)
		assert.Equal(t, ticket2.Status, tickets[1].Status)
		assert.NotNil(t, tickets[0].Admin)
		assert.Equal(t, admin.ID, *tickets[0].AdminID)
		assert.NotNil(t, tickets[1].Admin)
		assert.Equal(t, admin.ID, *tickets[1].AdminID)
	})

	t.Run("Success - With Filter", func(t *testing.T) {
		tickets, err := ticketRepo.GetTickets(ctx, models.TicketStatusOpen)

		assert.NoError(t, err)
		assert.Len(t, tickets, 1)
		assert.Equal(t, ticket1.ID, tickets[0].ID)
		assert.Equal(t, models.TicketStatusOpen, tickets[0].Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockTicketRepo := repositories.NewTicketRepo(mockDB)
		tickets, err := mockTicketRepo.GetTickets(ctx, "")

		assert.Empty(t, tickets)
		assert.Error(t, err)
	})
}

func TestGetTicketByID(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	db := testutils.NewMockDB()
	ticketRepo := repositories.NewTicketRepo(db)

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)
	db.Create(&ticket)

	t.Run("Success", func(t *testing.T) {
		result, err := ticketRepo.GetTicketByID(ctx, ticket.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ticket.ID, result.ID)
		assert.Equal(t, ticket.Name, result.Name)
		assert.Equal(t, ticket.Email, result.Email)
		assert.Equal(t, ticket.Status, result.Status)
		assert.NotNil(t, result.Admin)
		assert.Equal(t, admin.ID, result.Admin.ID)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		result, err := ticketRepo.GetTicketByID(ctx, uuid.New())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockTicketRepo := repositories.NewTicketRepo(mockDB)
		ticket, err := mockTicketRepo.GetTicketByID(ctx, ticket.ID)

		assert.Error(t, err)
		assert.Empty(t, ticket)
	})
}

func TestCreateTicket(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	ticketRepo := repositories.NewTicketRepo(db)

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		nil,
	)

	t.Run("Success", func(t *testing.T) {
		err := ticketRepo.CreateTicket(ctx, &ticket)
		assert.NoError(t, err)

		var result models.Ticket
		err = db.Where("id = ?", ticket.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, ticket.ID, result.ID)
		assert.Equal(t, ticket.Name, result.Name)
		assert.Equal(t, ticket.Status, result.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockTicketRepo := repositories.NewTicketRepo(mockDB)
		err = mockTicketRepo.CreateTicket(ctx, &models.Ticket{})

		assert.Error(t, err)
	})
}

func TestUpdateTicket(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	ticketRepo := repositories.NewTicketRepo(db)

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		nil,
	)
	db.Create(&ticket)

	t.Run("Success", func(t *testing.T) {
		ticketToUpdate := ticket
		ticketToUpdate.Status = models.TicketStatusResolved

		err := ticketRepo.UpdateTicket(ctx, &ticketToUpdate)
		assert.NoError(t, err)

		var result models.Ticket
		err = db.Where("id = ?", ticket.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, ticket.ID, result.ID)
		assert.Equal(t, models.TicketStatusResolved, result.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockTicketRepo := repositories.NewTicketRepo(mockDB)
		err = mockTicketRepo.UpdateTicket(ctx, &models.Ticket{})

		assert.Error(t, err)
	})
}

func TestDeleteTicket(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	db := testutils.NewMockDB()
	ticketRepo := repositories.NewTicketRepo(db)

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)
	db.Create(&ticket)

	t.Run("Success", func(t *testing.T) {
		err := ticketRepo.DeleteTicket(ctx, &ticket)
		assert.NoError(t, err)

		var result models.Ticket
		err = db.Where("id = ?", ticket.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockTicketRepo := repositories.NewTicketRepo(mockDB)
		err = mockTicketRepo.DeleteTicket(ctx, &models.Ticket{})

		assert.Error(t, err)
	})
}
