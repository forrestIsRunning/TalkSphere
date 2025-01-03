package encrypt

import (
	"crypto/md5"
	"encoding/hex"

	"TalkSphere/setting"
)

func EncryptPassword(opassword string) string {
	h := md5.New()
	h.Write([]byte(setting.Conf.EncryptConfig.SecretKey))
	return hex.EncodeToString(h.Sum([]byte(opassword)))
}
