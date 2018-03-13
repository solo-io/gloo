package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	. "github.com/solo-io/gloo-plugins/grpc/grpc-test-service/bookstore"
)

func NewServer() BookstoreServer {
	return &bookstoreServer{}
}

type bookstoreServer struct {
	shelves []*Shelf
}

// Returns a list of all shelves in the bookstore.
func (s *bookstoreServer) ListShelves(context.Context, *empty.Empty) (*ListShelvesResponse, error) {
	return &ListShelvesResponse{Shelves: s.shelves}, nil
}

// Creates a new shelf in the bookstore.
func (s *bookstoreServer) CreateShelf(ctx context.Context, r *CreateShelfRequest) (*Shelf, error) {
	s.shelves = append(s.shelves, r.Shelf)
	return r.Shelf, nil
}

// Creates multiple shelves with one streaming call
func (s *bookstoreServer) BulkCreateShelf(bc Bookstore_BulkCreateShelfServer) error {
	for {
		select {
		case <-bc.Context().Done():
			return nil
		default:
		}
		req, err := bc.Recv()
		if err != nil {
			return errors.Wrap(err, "receiving next create req")
		}
		if _, err := s.CreateShelf(bc.Context(), req); err != nil {
			return errors.Wrap(err, "creating shelf")
		}
	}
}

// Returns a specific bookstore shelf.
func (s *bookstoreServer) GetShelf(ctx context.Context, r *GetShelfRequest) (*Shelf, error) {
	for _, shelf := range s.shelves {
		if shelf.Id == r.Shelf {
			return shelf, nil
		}
	}
	return nil, errors.Errorf("shelf %v not found", r.Shelf)
}

// Deletes a shelf, including all books that are stored on the shelf.
func (s *bookstoreServer) DeleteShelf(context.Context, *DeleteShelfRequest) (*empty.Empty, error) {

}

// Returns a list of books on a shelf.
func (s *bookstoreServer) ListBooks(*ListBooksRequest, Bookstore_ListBooksServer) error {}

// Creates a new book.
func (s *bookstoreServer) CreateBook(context.Context, *CreateBookRequest) (*Book, error) {}

// Returns a specific book.
func (s *bookstoreServer) GetBook(context.Context, *GetBookRequest) (*Book, error) {}

// Deletes a book from a shelf.
func (s *bookstoreServer) DeleteBook(context.Context, *DeleteBookRequest) (*empty.Empty, error) {
}
func (s *bookstoreServer) UpdateBook(context.Context, *UpdateBookRequest) (*Book, error) {}
func (s *bookstoreServer) BookstoreOptions(context.Context, *GetShelfRequest) (*empty.Empty, error) {
}

// Returns a specific author.
func (s *bookstoreServer) GetAuthor(context.Context, *GetAuthorRequest) (*Author, error) {}
