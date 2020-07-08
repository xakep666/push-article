package token

import (
	"context"
	"sync"
)

type MemoryStorage struct {
	mu          sync.RWMutex
	userTokens  map[uint64][]Token
	tokenOwners map[string]uint64
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		userTokens:  map[uint64][]Token{},
		tokenOwners: map[string]uint64{},
	}
}

func (ms *MemoryStorage) SaveToken(ctx context.Context, userID uint64, token Token) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	owner, ok := ms.tokenOwners[token.Token]
	// if old user comes with some token it's ok
	if owner == userID {
		return nil
	}
	// if new user come with existing token we
	// should change it's owner to prevent push target mismatch
	if ok {
		ms.deleteTokenFromUser(token.Token, owner)
	}

	ut := ms.userTokens[userID]
	ut = append(ut, token)
	ms.userTokens[userID] = ut

	ms.tokenOwners[token.Token] = userID

	return nil
}

func (ms *MemoryStorage) deleteTokenFromUser(token string, userID uint64) {
	ut := ms.userTokens[userID]
	for i, t := range ut {
		if t.Token == token {
			ut[i], ut[len(ut)-1] = ut[len(ut)-1], Token{}
			ut = ut[:len(ut)-1]
			break
		}
	}
	ms.userTokens[userID] = ut
}

func (ms *MemoryStorage) UserTokens(ctx context.Context, userID uint64) ([]Token, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	tokens := ms.userTokens[userID]
	ret := make([]Token, len(tokens))
	copy(ret, tokens)

	return ret, nil
}

func (ms *MemoryStorage) DeleteTokens(ctx context.Context, tokens []string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, token := range tokens {
		user, ok := ms.tokenOwners[token]
		if !ok {
			return nil
		}

		ms.deleteTokenFromUser(token, user)
	}
	return nil
}
