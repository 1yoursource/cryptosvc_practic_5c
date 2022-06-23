package cryptosvc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"projects/practic_5course_cesar/internal/storage"
	"projects/practic_5course_cesar/pkg/custom"
	"strings"
	"time"
)

type CryptoService struct {
	enCrypt custom.LikeCesarCrypt
	ukCrypt custom.LikeCesarCrypt
	cache   Storage
}

type Storage interface {
	Set(ctx context.Context, k storage.Key, v *storage.DecryptInfo)
	Get(ctx context.Context, k storage.Key) (*storage.DecryptInfo, error)
}

func New(en, uk custom.LikeCesarCrypt, storage Storage) *CryptoService {
	return &CryptoService{
		enCrypt: en,
		ukCrypt: uk,
		cache:   storage,
	}
}

type Data string

func (s Data) String() string {
	return string(s)
}

type Language string

const ukrainian Language = "uk"
const english Language = "en"

type Result string

func (s *CryptoService) Encrypt(ctx context.Context, lang Language, phrase Data) (Result, error) {
	handler, err := s.getHandlerByLang(lang)
	if err != nil {
		return "", fmt.Errorf("indetidy handler error: %w", err)
	}

	data := strings.ToUpper(phrase.String())

	var encResult string
	var decryptKey []int

	for _, runa := range data {
		shiftBy := custom.Shift(s.getRandomDigit())

		result, shift := handler.Crypt(ctx, custom.Data(runa), shiftBy)

		encResult += result.String()

		decryptKey = append(decryptKey, -shift.Int())
	}

	s.cache.Set(ctx, storage.Key(encResult), &storage.DecryptInfo{Key: decryptKey})

	return Result(encResult), nil
}

func (s *CryptoService) Decrypt(ctx context.Context, lang Language, phrase Data) (Result, error) {
	data := strings.ToUpper(phrase.String())

	decryptInfo, err := s.cache.Get(ctx, storage.Key(data))
	if err != nil {
		return "", fmt.Errorf("get cache for decrypting: %w", err)
	}

	handler, err := s.getHandlerByLang(lang)
	if err != nil {
		return "", fmt.Errorf("indetidy handler error: %w", err)
	}

	var encResult string
	var count int

	for _, runa := range data {
		result, _ := handler.Crypt(ctx, custom.Data(runa), custom.Shift(decryptInfo.Key[count]))

		encResult += result.String()

		count++
	}

	return Result(encResult), nil
}

var ErrUnknownLanguage = errors.New("unknown language")

func (s *CryptoService) getHandlerByLang(l Language) (custom.LikeCesarCrypt, error) {
	if l == english {
		return s.enCrypt, nil
	}

	if l == ukrainian {
		return s.ukCrypt, nil
	}

	return nil, ErrUnknownLanguage
}

func (s *CryptoService) getRandomDigit() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}
