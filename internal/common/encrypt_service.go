package common


import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/cast5"
)

type IEncryptDecryptService interface {
	Tokenization(salt, token, txt string) (map[string]bool, error)
	EncryptAES(aesKey []byte, plaintext string) (string, error)
	DecryptAES(aesKey []byte, enc string) (string, error)
	NGrams(str string, n int) []string
	HashToken(salt, token string) string
	DecryptOpenPGPCFB(key, hexCipher string) (string, error)
}

type encryptDecryptService struct {
	ctx   context.Context
	debug bool
}

func NewEncryptDecryptService(
	ctx context.Context, debug bool,
) IEncryptDecryptService {
	return &encryptDecryptService{
		ctx:   ctx,
		debug: debug,
	}
}

func (s *encryptDecryptService) Tokenization(salt, token, txt string) (map[string]bool, error) {
	index := make(map[string]bool)
	parts := strings.Fields(txt)
	for _, p := range parts {
		for _, g := range s.NGrams(p, 2) {
			index[s.HashToken(salt, g)] = true
		}
	}

	return index, nil
}

func (s *encryptDecryptService) EncryptAES(aesKey []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *encryptDecryptService) DecryptAES(aesKey []byte, enc string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

/* =========================
   SEARCH INDEX
========================= */

// สร้าง n-gram (รองรับภาษาไทย)
func (s *encryptDecryptService) NGrams(str string, n int) []string {
	r := []rune(str)
	var out []string

	for i := 0; i <= len(r)-n; i++ {
		out = append(out, string(r[i:i+n]))
	}
	return out
}

func (s *encryptDecryptService) HashToken(salt, token string) string {
	sum := sha256.Sum256([]byte(salt + token))
	return hex.EncodeToString(sum[:])
}

var reWS = regexp.MustCompile(`\s+`)

func stripWS(s string) string { return reWS.ReplaceAllString(s, "") }

func (s *encryptDecryptService) DecryptOpenPGPCFB(key, hexCipher string) (string, error) {

	padded := make([]byte, 16)
	copy(padded, []byte(key))

	ct, err := hex.DecodeString(stripWS(hexCipher))
	if err != nil {
		return "", err
	}
	if len(ct) < cast5.BlockSize+2 {
		return "", errors.New("ciphertext สั้นเกิน (ต้อง ≥ 10 ไบต์)")
	}

	block, err := cast5.NewCipher(padded)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize() // 8

	plain := make([]byte, 0, len(ct))

	/* ---------- ขั้น 1 : ถอด 8 ไบต์แรก ---------- */
	iv0 := make([]byte, bs)                    // ศูนย์ 8 ไบต์
	cfb1 := cipher.NewCFBDecrypter(block, iv0) // CFB-64
	prefix8 := make([]byte, bs)
	cfb1.XORKeyStream(prefix8, ct[:bs])
	plain = append(plain, prefix8...)

	/* ---------- ขั้น 2 : ถอด 2 ไบต์ sync ---------- */
	FR := ct[:bs] // feedback = ciphertext[0-7]
	var fre [8]byte
	block.Encrypt(fre[:], FR) // ENC_K(FR)

	p8 := fre[0] ^ ct[bs]   // ไบต์ที่ 8
	p9 := fre[1] ^ ct[bs+1] // ไบต์ที่ 9
	plain = append(plain, p8, p9)

	// เตรียม IV ใหม่สำหรับส่วนที่เหลือ
	iv1 := append(FR[2:], ct[bs], ct[bs+1]) // len = 8

	/* ---------- ขั้น 3 : ถอด ciphertext ที่เหลือ ---------- */
	restCT := ct[bs+2:]
	if len(restCT) > 0 {
		cfb2 := cipher.NewCFBDecrypter(block, iv1)
		restPT := make([]byte, len(restCT))
		cfb2.XORKeyStream(restPT, restCT)
		plain = append(plain, restPT...)
	}

	// ---------- ตัด random-prefix 10 ไบต์ ----------
	if len(plain) < cast5.BlockSize+2 {
		log.Fatal("plaintext สั้นเกินคาด")
	}
	jsonBytes := plain[cast5.BlockSize+2:]

	// (ทางเลือก) slice จนเจอ '{' เผื่อฟอร์แมตไม่ได้ตามสเปก
	if idx := strings.IndexByte(string(jsonBytes), '{'); idx >= 0 {
		jsonBytes = jsonBytes[idx:]
	}

	return string(jsonBytes), nil
}

