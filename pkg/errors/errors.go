package errors

import "errors"

// Erreurs domaine — ce que retournent les usecases
var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")     // ex: email déjà utilisé
	ErrForbidden  = errors.New("forbidden")    // ex: pas le droit de supprimer
	ErrBadRequest = errors.New("bad request")  // ex: champs manquants
	ErrUnauth     = errors.New("unauthorized") // ex: token invalide
)
