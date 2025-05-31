package ambulance_virtual_patient_list

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Robert-Prikryl/ambulance-virtual-patient-api/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type VirtualPatientListSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[VirtualPatient]
}

func TestVirtualPatientListSuite(t *testing.T) {
	suite.Run(t, new(VirtualPatientListSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) ListDocuments(ctx context.Context) ([]DocType, error) {
	args := this.Called(ctx)
	return args.Get(0).([]DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *VirtualPatientListSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[VirtualPatient]{}

	// Compile time Assert that the mock is of type db_service.DbService[VirtualPatient]
	var _ db_service.DbService[VirtualPatient] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&VirtualPatient{
				Id:         "test-patient",
				Name:       "Test Patient",
				Difficulty: 2,
				Symptoms:   []string{"fever", "cough"},
				Anamnesis:  "Test anamnesis",
			},
			nil,
		)
}

func (suite *VirtualPatientListSuite) Test_CreateVirtualPatient_DbServiceCreateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("CreateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"name": "New Patient",
		"difficulty": 3,
		"symptoms": ["headache", "dizziness"],
		"anamnesis": "New patient anamnesis"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Request = httptest.NewRequest("POST", "/virtual-patient", strings.NewReader(json))

	// ACT
	NewPatientApi().CreateVirtualPatient(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "CreateDocument", mock.Anything, mock.Anything, mock.Anything)
	suite.Equal(http.StatusCreated, recorder.Code)
}

func (suite *VirtualPatientListSuite) Test_UpdateVirtualPatient_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"name": "Updated Patient",
		"difficulty": 4,
		"symptoms": ["fever", "cough", "fatigue"],
		"anamnesis": "Updated anamnesis"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "virtualPatientId", Value: "test-patient"},
	}
	ctx.Request = httptest.NewRequest("PUT", "/virtual-patient/test-patient", strings.NewReader(json))

	// ACT
	NewPatientApi().UpdateVirtualPatient(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-patient", mock.Anything)
	suite.Equal(http.StatusOK, recorder.Code)
}

func (suite *VirtualPatientListSuite) Test_DeleteVirtualPatient_DbServiceDeleteCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("DeleteDocument", mock.Anything, mock.Anything).
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "virtualPatientId", Value: "test-patient"},
	}
	ctx.Request = httptest.NewRequest("DELETE", "/virtual-patient/test-patient", nil)

	// ACT
	NewPatientApi().DeleteVirtualPatient(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "DeleteDocument", mock.Anything, "test-patient")
	suite.Equal(http.StatusNoContent, recorder.Code)
}

func (suite *VirtualPatientListSuite) Test_GetVirtualPatientList_DbServiceListCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("ListDocuments", mock.Anything).
		Return([]VirtualPatient{
			{
				Id:         "test-patient-1",
				Name:       "Test Patient 1",
				Difficulty: 2,
				Symptoms:   []string{"fever", "cough"},
				Anamnesis:  "Test anamnesis 1",
			},
			{
				Id:         "test-patient-2",
				Name:       "Test Patient 2",
				Difficulty: 3,
				Symptoms:   []string{"headache", "dizziness"},
				Anamnesis:  "Test anamnesis 2",
			},
		}, nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Request = httptest.NewRequest("GET", "/virtual-patient", nil)

	// ACT
	NewPatientApi().GetVirtualPatientList(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "ListDocuments", mock.Anything)
	suite.Equal(http.StatusOK, recorder.Code)
}
