package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/lucasreed/smol/pkg/data/models"
)

type Store struct {
	DBName     string
	User       string
	Password   string
	Host       string
	Port       string
	Connection *gorm.DB
}

// NewStore represents a new instance of a mysql storage location
func NewStore(host, port, user, password, database string) *Store {
	return &Store{
		DBName:   database,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (s *Store) GetURL(shortCode string) (models.URL, error) {
	url := models.URL{}
	err := s.Connection.Where(&models.URL{
		ShortCode: shortCode,
	}).First(&url).Error
	return url, err
}

func (s *Store) GetShortCode(destination string) (string, error) {
	url := models.URL{}
	err := s.Connection.Where(&models.URL{
		Destination: destination,
	}).First(&url).Error
	return url.ShortCode, err
}

func (s *Store) Health() bool {
	if err := s.Connection.DB().Ping(); err != nil {
		return false
	}
	return true
}

func (s *Store) Open() error {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", s.User, s.Password, s.Host, s.Port, s.DBName))
	if err != nil {
		return err
	}
	s.Connection = db
	if !s.Health() {
		return fmt.Errorf("[mysql] connection not established")
	}
	db.AutoMigrate(&models.URL{})
	return nil
}

func (s *Store) Close() error {
	return s.Connection.Close()
}

func (s *Store) SetURL(shortCode, url string) error {
	urlRow := models.URL{
		Destination: url,
		ShortCode:   shortCode,
	}
	return s.Connection.Create(&urlRow).Error
}

func (s *Store) Delete(shortCode string) error {
	url := models.URL{}
	if err := s.Connection.Where(&models.URL{
		ShortCode: shortCode,
	}).First(&url).Error; err != nil {
		return err
	}
	return s.Connection.Delete(&url).Error
}
