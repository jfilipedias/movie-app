package model

import model "github.com/jfilipedias/movie-app/metadata/pkg"

type MovieDetails struct {
	Rating   *float64       `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
