package booking

import "sync"

type ConcurrentMemoryStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentMemoryStore() *ConcurrentMemoryStore {
	return &ConcurrentMemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentMemoryStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentMemoryStore) ListBookings(movieID string) []Booking {
	s.RLock() // just for read data
	defer s.RUnlock()
	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}
