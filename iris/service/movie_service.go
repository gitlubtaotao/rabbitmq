package service

import "rabbitmq/iris/repositories"

type MovieService interface {
	ShowMovieName() string
}

type MovieServiceManager struct {
	repo repositories.MovieRepository
}

func NewMovieServiceManager(repo repositories.MovieRepository) MovieService {
	return  &MovieServiceManager{repo: repo}
}

func (m *MovieServiceManager) ShowMovieName() string {
	name :=  m.repo.GetMovieName()
	return name
}
