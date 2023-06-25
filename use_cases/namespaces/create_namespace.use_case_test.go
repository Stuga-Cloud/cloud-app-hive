package namespaces

import (
	"fmt"
	"testing"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

// MockNamespaceRepository is a mock implementation of the NamespaceRepository interface above
type MockNamespaceRepository struct {
	FindByIDFunc     func(id string) (*domain.Namespace, error)
	ExistsByNameFunc func(name string) (bool, error)
	FindByNameFunc   func(name string) (*domain.Namespace, error)
	FindFunc         func(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error)
	CreateFunc       func(createNamespace commands.CreateNamespace) (*domain.Namespace, error)
	DeleteFunc       func(id string, userId string) (*domain.Namespace, error)
	UpdateFunc       func(updateNamespace commands.UpdateNamespace) (*domain.Namespace, error)
}

func (m *MockNamespaceRepository) FindByID(id string) (*domain.Namespace, error) {
	return m.FindByIDFunc(id)
}

func (m *MockNamespaceRepository) FindByName(name string) (*domain.Namespace, error) {
	return m.FindByNameFunc(name)
}

func (m *MockNamespaceRepository) ExistsByName(name string) (bool, error) {
	return m.ExistsByNameFunc(name)
}

func (m *MockNamespaceRepository) Find(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error) {
	return m.FindFunc(findNamespaces)
}

func (m *MockNamespaceRepository) Create(createNamespace commands.CreateNamespace) (*domain.Namespace, error) {
	return m.CreateFunc(createNamespace)
}

func (m *MockNamespaceRepository) Delete(id string, userId string) (*domain.Namespace, error) {
	return m.DeleteFunc(id, userId)
}

func (m *MockNamespaceRepository) Update(updateNamespace commands.UpdateNamespace) (*domain.Namespace, error) {
	return m.UpdateFunc(updateNamespace)
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
		ExistsByNameFunc: func(name string) (bool, error) {
			return false, nil
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
		ExistsByNameFunc: func(name string) (bool, error) {
			return true, nil
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
	expectedErrorMessage := fmt.Sprintf("namespace with name %s already exists", createNamespace.Name)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMessage, err.Error())
	}
}
