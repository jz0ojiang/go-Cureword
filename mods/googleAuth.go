package mods

// CODE FROM: https://my.oschina.net/qiongtaoli/blog/3049118

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (gauth *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (gauth *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (gauth *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (gauth *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (gauth *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (gauth *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (gauth *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := gauth.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := gauth.toUint32(hashParts)
	return number % 1000000
}

func (gauth *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, gauth.un())
	return strings.ToUpper(gauth.base32encode(gauth.hmacSha1(buf.Bytes(), nil)))
}

func (gauth *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := gauth.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := gauth.oneTimePassword(secretKey, gauth.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

func (gauth *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

func (gauth *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := gauth.GetQrcode(user, secret)
	return fmt.Sprintf("http://www.google.cn/chart?chs=200x200&chld=M%%7C0&cht=qr&chl=%s", qrcode)
}

func (gauth *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := gauth.GetCode(secret)
	fmt.Println(_code, code, err)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}
