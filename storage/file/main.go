package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type CacheFile struct {
	Data    []byte              `json:"data"`
	Headers map[string][]string `json:"headers"`
}

type CacheStore struct{}

func (s CacheStore) getFilename(key string) string {
	filename := fmt.Sprintf("blitz_%s.dat", key)
	return path.Join(os.TempDir(), filename)
}

func (s *CacheStore) Exists(key string) bool {
	filename := s.getFilename(key)
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (s *CacheStore) Save(key string, data []byte, headers map[string][]string) error {
	filename := s.getFilename(key)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	cacheFile := CacheFile{
		Data:    data,
		Headers: headers,
	}
	content, err := json.Marshal(&cacheFile)
	if err != nil {
		return fmt.Errorf("could not marshal data to json: %w", err)
	}

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (s *CacheStore) Get(key string) ([]byte, map[string][]string, error) {
	filename := s.getFilename(key)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading cache file: %w", err)
	}
	var cacheFile CacheFile
	if err := json.Unmarshal(content, &cacheFile); err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling cache data: %w", err)
	}
	return cacheFile.Data, cacheFile.Headers, nil
}

func (s *CacheStore) DeleteById(key string) error {
	filename := s.getFilename(key)
	err := os.Remove(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
