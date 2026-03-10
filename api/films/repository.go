package films

import "fmt"

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Stub data — replace with DB/external API calls
var store = []Film{
	{ID: "1", Title: "The Godfather", Year: 1972},
	{ID: "2", Title: "Inception", Year: 2010},
}

func (r *Repository) FindAll() ([]Film, error) {
	return store, nil
}

func (r *Repository) FindByID(id string) (*Film, error) {
	for _, f := range store {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("film %s not found", id)
}
