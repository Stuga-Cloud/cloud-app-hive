package namespaces

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"fmt"
	"testing"
)

// MockNamespaceRepository is a mock implementation of the NamespaceRepository interface above
type MockNamespaceRepository struct {
	FindByIDFunc   func(id string) (*domain.Namespace, error)
	FindByNameFunc func(name string) (*domain.Namespace, error)
	FindFunc       func(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error)
	CreateFunc     func(createNamespace commands.CreateNamespace) (*domain.Namespace, error)
	DeleteFunc     func(id string) (*domain.Namespace, error)
}

func (m *MockNamespaceRepository) FindByID(id string) (*domain.Namespace, error) {
	return m.FindByIDFunc(id)
}

func (m *MockNamespaceRepository) FindByName(name string) (*domain.Namespace, error) {
	return m.FindByNameFunc(name)
}

func (m *MockNamespaceRepository) Find(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error) {
	return m.FindFunc(findNamespaces)
}

func (m *MockNamespaceRepository) Create(createNamespace commands.CreateNamespace) (*domain.Namespace, error) {
	return m.CreateFunc(createNamespace)
}

func (m *MockNamespaceRepository) Delete(id string) (*domain.Namespace, error) {
	return m.DeleteFunc(id)
}

func TestExecute_CreateNamespace_Success(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepository := &MockNamespaceRepository{
		FindByNameFunc: func(name string) (*domain.Namespace, error) {
			return nil, nil // Namespace doesn't exist
		},
		CreateFunc: func(createNamespace commands.CreateNamespace) (*domain.Namespace, error) {
			return &domain.Namespace{ID: "123", Name: createNamespace.Name}, nil
		},
	}

	// Create the use case instance with the mock repository
	useCase := CreateNamespaceUseCase{
		NamespaceRepository: mockRepository,
	}

	// Prepare the test case
	createNamespace := commands.CreateNamespace{Name: "myNamespace"}

	// Execute the use case
	createdNamespace, err := useCase.Execute(createNamespace)

	// Assert the results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if createdNamespace == nil {
		t.Error("Expected a non-nil created namespace, but got nil")
	}
	if createdNamespace.Name != createNamespace.Name {
		t.Errorf("Expected created namespace name to be %s, but got %s", createNamespace.Name, createdNamespace.Name)
	}
}

func TestExecute_CreateNamespace_AlreadyExists(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepository := &MockNamespaceRepository{
		FindByNameFunc: func(name string) (*domain.Namespace, error) {
			return &domain.Namespace{ID: "123", Name: name}, nil // Namespace already exists
		},
	}

	// Create the use case instance with the mock repository
	useCase := CreateNamespaceUseCase{
		NamespaceRepository: mockRepository,
	}

	// Prepare the test case
	createNamespace := commands.CreateNamespace{Name: "myNamespace"}

	// Execute the use case
	createdNamespace, err := useCase.Execute(createNamespace)

	// Assert the results
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if createdNamespace != nil {
		t.Error("Expected createdNamespace to be nil, but got non-nil")
	}
	expectedErrorMessage := fmt.Sprintf("namespace %s already exists", createNamespace.Name)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMessage, err.Error())
	}
}
