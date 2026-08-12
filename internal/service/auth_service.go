package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"project-manager-go/internal/auth"
	"project-manager-go/internal/domain"
)

func (s *Service) Authenticate(token string) (domain.AuthUser, error) {
	return s.tokens.VerifyAccessToken(token)
}

func (s *Service) CurrentUser(ctx context.Context, authUser domain.AuthUser) (domain.User, error) {
	user, err := s.store.GetUserByID(ctx, authUser.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("getting current user: %w", err)
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *Service) Register(ctx context.Context, input CreateUserInput) (domain.User, error) {
	if err := validateUserRole(input.Role); err != nil {
		return domain.User{}, err
	}

	_, err := s.store.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return domain.User{}, domain.NewError(domain.ErrDuplicate, "User already exists")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, fmt.Errorf("checking existing user: %w", err)
	}

	return s.CreateUser(ctx, input)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	user, err := s.store.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, domain.NewError(domain.ErrInvalidCredential, "Invalid credentials")
		}
		return AuthResult{}, fmt.Errorf("getting user by email: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return AuthResult{}, domain.NewError(domain.ErrInvalidCredential, "Invalid credentials")
	}

	accessToken, err := s.tokens.GenerateAccessToken(user)
	if err != nil {
		return AuthResult{}, err
	}
	refreshToken, err := s.tokens.GenerateRefreshToken(user)
	if err != nil {
		return AuthResult{}, err
	}

	err = s.store.SaveRefreshToken(
		ctx,
		user.ID,
		auth.HashRefreshToken(refreshToken),
		s.now().Add(s.tokens.RefreshTTL),
	)
	if err != nil {
		return AuthResult{}, fmt.Errorf("saving refresh token: %w", err)
	}

	user.PasswordHash = ""
	return AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (AuthResult, error) {
	if refreshToken == "" {
		return AuthResult{}, domain.NewError(domain.ErrInvalidRefresh, "Refresh token not found")
	}

	authUser, err := s.tokens.VerifyRefreshToken(refreshToken)
	if err != nil {
		return AuthResult{}, domain.NewError(domain.ErrInvalidRefresh, "Invalid refresh token")
	}

	tokenHash := auth.HashRefreshToken(refreshToken)
	ok, err := s.store.FindValidRefreshToken(ctx, authUser.ID, tokenHash, s.now())
	if err != nil {
		return AuthResult{}, fmt.Errorf("finding refresh token: %w", err)
	}
	if !ok {
		return AuthResult{}, domain.NewError(domain.ErrInvalidRefresh, "Invalid refresh token")
	}

	user, err := s.store.GetUserByID(ctx, authUser.ID)
	if err != nil {
		return AuthResult{}, domain.NewError(domain.ErrInvalidRefresh, "Invalid refresh token")
	}

	if err := s.store.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return AuthResult{}, fmt.Errorf("revoking refresh token: %w", err)
	}

	newAccessToken, err := s.tokens.GenerateAccessToken(user)
	if err != nil {
		return AuthResult{}, err
	}
	newRefreshToken, err := s.tokens.GenerateRefreshToken(user)
	if err != nil {
		return AuthResult{}, err
	}

	err = s.store.SaveRefreshToken(
		ctx,
		user.ID,
		auth.HashRefreshToken(newRefreshToken),
		s.now().Add(s.tokens.RefreshTTL),
	)
	if err != nil {
		return AuthResult{}, fmt.Errorf("saving rotated refresh token: %w", err)
	}

	return AuthResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, user domain.AuthUser, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	tokenHash := auth.HashRefreshToken(refreshToken)
	ok, err := s.store.RefreshTokenBelongsToUser(ctx, user.ID, tokenHash)
	if err != nil {
		return fmt.Errorf("checking refresh token owner: %w", err)
	}
	if !ok {
		return domain.NewError(domain.ErrForbidden, "Refresh token not found or does not belong to user")
	}

	return s.store.RevokeRefreshToken(ctx, tokenHash)
}
