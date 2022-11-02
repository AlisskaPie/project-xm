package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AlisskaPie/project-xm/pkg/domain"
	"github.com/AlisskaPie/project-xm/pkg/domain/mocks"
)

var testUUID = uuid.MustParse("10000000-0000-0000-0000-000000000000")

func TestEventSenderWrapper_CreateSuccess(t *testing.T) {
	m := &mocks.CompanyRepository{}
	m.On("Create", mock.Anything, mock.Anything).Return(nil)

	createCompany := domain.CreateCompany{
		ID:                testUUID,
		Name:              "1",
		Description:       "2",
		AmountOfEmployees: 3,
		Registered:        true,
		CompanyType:       domain.CooperativeType,
	}
	e := &mocks.CompanyEventSender{}
	e.On("Send", mock.Anything, domain.CompanyEvent{
		Action: domain.InsertEventActionType,
		ID:     createCompany.ID,
		State:  domain.Company(createCompany),
	}).Return(nil)
	w := NewEventSenderWrapper(m, e)

	err := w.Create(context.TODO(), createCompany)
	assert.NoError(t, err)

	m.AssertExpectations(t)
	e.AssertExpectations(t)
}

func TestEventSenderWrapper_DeleteSuccess(t *testing.T) {
	m := &mocks.CompanyRepository{}
	m.On("Delete", mock.Anything, testUUID).Return(nil)

	e := &mocks.CompanyEventSender{}
	e.On("Send", mock.Anything, domain.CompanyEvent{
		Action: domain.DeleteEventActionType,
		ID:     testUUID,
	}).Return(nil)
	w := NewEventSenderWrapper(m, e)

	err := w.Delete(context.TODO(), testUUID)
	assert.NoError(t, err)

	m.AssertExpectations(t)
	e.AssertExpectations(t)
}

func TestEventSenderWrapper_PatchSuccess(t *testing.T) {
	testPatchCompany := domain.PatchCompany{
		Name:              getPointer("1"),
		Description:       getPointer("2"),
		AmountOfEmployees: getPointer(uint32(3)),
		Registered:        getPointer(true),
		CompanyType:       getPointer(domain.CooperativeType),
	}
	testExpCompany := domain.Company{
		Name:              "1",
		Description:       "2",
		AmountOfEmployees: 3,
		Registered:        true,
		CompanyType:       domain.CooperativeType,
	}
	m := &mocks.CompanyRepository{}
	m.On("Patch", mock.Anything, testUUID, testPatchCompany).
		Return(domain.Company{
			Name:              "1",
			Description:       "2",
			AmountOfEmployees: 3,
			Registered:        true,
			CompanyType:       domain.CooperativeType,
		}, nil)

	e := &mocks.CompanyEventSender{}
	e.On("Send", mock.Anything, domain.CompanyEvent{
		Action: domain.UpdateEventActionType,
		ID:     testUUID,
		State:  testExpCompany,
	}).Return(nil)
	w := NewEventSenderWrapper(m, e)

	_, err := w.Patch(context.TODO(), testUUID, testPatchCompany)
	assert.NoError(t, err)

	m.AssertExpectations(t)
	e.AssertExpectations(t)
}
