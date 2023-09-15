package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ruhancs/bank-go/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	
	username := util.RandomOwner()
	duration := time.Minute
	
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	
	token,err := maker.CreateToken(username,duration)
	require.NoError(t, err)
	require.NotEmpty(t,token)
	
	payload,err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t,payload)

	require.NotZero(t,payload.ID)
	require.Equal(t,username,payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssueAt,time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt,time.Second)

}

func TestExpiredJWTTokenCase(t *testing.T) {
	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token,err := maker.CreateToken(util.RandomOwner(),-time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t,token)
	
	payload,err := maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t, err, ErrorExpiredToken.Error())	
	require.Nil(t,payload)
}

func TestInvalidJWTTokenCase(t *testing.T) {
	payload,err := NewPaylload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token,err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload,err = maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t,payload)
}