package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ndavidson19/quanta-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:            util.RandomUsername(),
		HashedPassword:      hashedPassword,
		FullName:            util.RandomUsername(),
		Email:               util.RandomEmail(),
		PhoneNumber:         sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		PasswordChangedAt:   time.Now().UTC(),
		CreatedAt:           time.Now().UTC(),
		LastLoginAt:         sql.NullTime{Time: time.Now().UTC(), Valid: true},
		LoginAttempts:       sql.NullInt32{Int32: 0, Valid: true},
		LockedUntil:         sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ResetToken:          sql.NullString{String: util.RandomString(6), Valid: true},
		ResetTokenExpiresAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.PasswordChangedAt, user.PasswordChangedAt)
	require.Equal(t, arg.CreatedAt, user.CreatedAt)
	require.Equal(t, arg.LastLoginAt, user.LastLoginAt)
	require.Equal(t, arg.LoginAttempts, user.LoginAttempts)
	require.Equal(t, arg.LockedUntil, user.LockedUntil)
	require.Equal(t, arg.ResetToken, user.ResetToken)
	require.Equal(t, arg.ResetTokenExpiresAt, user.ResetTokenExpiresAt)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.Equal(t, user1.LastLoginAt, user2.LastLoginAt)
	require.Equal(t, user1.LoginAttempts, user2.LoginAttempts)
	require.Equal(t, user1.LockedUntil, user2.LockedUntil)
	require.Equal(t, user1.ResetToken, user2.ResetToken)
	require.Equal(t, user1.ResetTokenExpiresAt, user2.ResetTokenExpiresAt)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserParams{
		Username:            user1.Username,
		HashedPassword:      util.RandomString(6),
		FullName:            util.RandomUsername(),
		Email:               util.RandomEmail(),
		PhoneNumber:         sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		PasswordChangedAt:   time.Now().UTC(),
		LastLoginAt:         sql.NullTime{Time: time.Now().UTC(), Valid: true},
		LoginAttempts:       sql.NullInt32{Int32: 0, Valid: true},
		LockedUntil:         sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ResetToken:          sql.NullString{String: util.RandomString(6), Valid: true},
		ResetTokenExpiresAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}

	// Execute the update and get the updated user
	err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	// Assert all updated fields
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.Equal(t, arg.FullName, user2.FullName)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.PhoneNumber.String, user2.PhoneNumber.String)
	require.WithinDuration(t, arg.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, arg.LastLoginAt.Time, user2.LastLoginAt.Time, time.Second)
	require.Equal(t, arg.LoginAttempts.Int32, user2.LoginAttempts.Int32)
	require.WithinDuration(t, arg.LockedUntil.Time, user2.LockedUntil.Time, time.Second)
	require.Equal(t, arg.ResetToken.String, user2.ResetToken.String)
	require.WithinDuration(t, arg.ResetTokenExpiresAt.Time, user2.ResetTokenExpiresAt.Time, time.Second)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
