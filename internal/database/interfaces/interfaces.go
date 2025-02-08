package interfaces

type Repository[T any] interface {
	Get() ([]T, error)
	GetById(int) (T, error)
	Insert(T) (int, error)
	Update(T) error
	Delete(int) (T, error)
}
