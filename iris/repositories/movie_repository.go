package repositories

import "rabbitmq/iris/datamodels"

type MovieRepository interface {
	GetMovieName() string
}

type MovieManager struct {
}

func NewMovieManager() MovieRepository {
	return &MovieManager{}
}

func (m *MovieManager) GetMovieName() string {
	movie := datamodels.Movie{Name: "immo慕课网"}
	return movie.Name
}
