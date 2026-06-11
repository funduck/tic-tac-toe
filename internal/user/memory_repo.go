package user

import "sync"

type MemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{
		users: make(map[string]*User),
	}
}

func (r *MemoryUserRepo) FindByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MemoryUserRepo) FindBySessionID(sessionID string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.ID == sessionID { // This is just a placeholder. In a real implementation, you'd check the session ID.
			return user, nil
		}
	}
	return nil, nil
}

func (r *MemoryUserRepo) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}
