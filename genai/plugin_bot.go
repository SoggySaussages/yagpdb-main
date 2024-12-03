package genai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strconv"

	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"golang.org/x/crypto/scrypt"
)

func (p *Plugin) BotInit() {

}

func createKey(gs *dstate.GuildState) ([]byte, error) {
	salt := []byte(strconv.FormatInt(gs.ID+gs.OwnerID, 10))
	return scrypt.Key([]byte(common.GetBotToken()), salt, 1048576, 8, 1, 32)
}

func encryptAPIToken(gs *dstate.GuildState, token string) (string, error) {
	key, err := createKey(gs)
	if err != nil {
		return "", err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	cypheredToken := gcm.Seal(nonce, nonce, []byte(token), nil)

	return string(cypheredToken), nil
}

func decryptAPIToken(gs *dstate.GuildState, encryptedToken string) (string, error) {
	key, err := createKey(gs)
	if err != nil {
		return "", err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	encryptedTokenBytes := []byte(encryptedToken)
	nonce, encryptedTokenBytes := encryptedTokenBytes[:gcm.NonceSize()], encryptedTokenBytes[gcm.NonceSize():]

	decryptedToken, err := gcm.Open(nil, nonce, encryptedTokenBytes, nil)
	if err != nil {
		return "", ErrorAPIKeyInvalid
	}

	return string(decryptedToken), nil
}

func getAPIToken(gs *dstate.GuildState) (string, error) {
	config, err := GetConfig(gs.ID)
	if err != nil {
		logger.WithError(err).WithField("guild", gs.ID).Error("Failed retrieving openai config")
		return "", err
	}

	if !config.Enabled {
		return "", nil
	}

	return decryptAPIToken(gs, config.Key)
}
