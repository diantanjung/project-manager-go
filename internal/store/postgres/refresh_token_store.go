package postgres

import (
	"context"
	"time"
)

func (s *Store) SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (s *Store) FindValidRefreshToken(ctx context.Context, userID int, tokenHash string, now time.Time) (bool, error) {
	var id int
	err := s.db.GetContext(ctx, &id, `
		SELECT id FROM refresh_tokens
		WHERE user_id = $1 AND token = $2 AND is_revoked = false AND expires_at >= $3
	`, userID, tokenHash, now)
	return exists(err)
}

func (s *Store) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET is_revoked = true WHERE token = $1`, tokenHash)
	return err
}

func (s *Store) RefreshTokenBelongsToUser(ctx context.Context, userID int, tokenHash string) (bool, error) {
	var id int
	err := s.db.GetContext(ctx, &id, `
		SELECT id FROM refresh_tokens WHERE user_id = $1 AND token = $2
	`, userID, tokenHash)
	return exists(err)
}
