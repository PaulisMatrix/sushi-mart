package testserver

import (
	"encoding/base64"
	"os"
	"sushi-mart/common"

	"github.com/sirupsen/logrus"
)

const testingCustID = 1

func (s *GinServer) GetJWTToken(config *common.Config, logger *logrus.Logger) string {
	//generate 1hr long token and return
	token, tokenErr := common.GenerateNewToken(int32(testingCustID), config)

	if tokenErr != nil {
		logger.WithError(tokenErr).Error("failed to generate the jwt token")
		os.Exit(1)
	}
	return token
}

func (s *GinServer) GetAdminBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
