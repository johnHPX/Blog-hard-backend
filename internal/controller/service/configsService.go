package service

import (
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/blog-hard-backend/internal/utils/messages"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type configsServiceInterface interface {
	Store(collors, links, menuAs []string, bannerURL string) error
	List(offset, limit, page int) ([]model.Configs, int, error)
	Find(configID int) (*model.Configs, error)
	Update(id int, collors, links, menuAs []string, bannerURL string) error
	Remove(id int) error
}

type configsServiceImpl struct {
	userID string
	kindID string
}

func (s *configsServiceImpl) Store(collors, links, menuAs []string, bannerURL string) error {

	if s.kindID != "adm" {
		return errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()
	collorsVal := make([]string, 0)
	for _, v := range collors {
		collor, err := val.CheckAnyData("cor do site", 16, v, true)
		if err != nil {
			return err
		}
		collorsVal = append(collorsVal, collor.(string))
	}
	linksVal := make([]string, 0)
	for _, v := range links {
		link, err := val.CheckAnyData("link do site", 255, v, true)
		if err != nil {
			return err
		}
		linksVal = append(linksVal, link.(string))
	}
	menuAsVal := make([]string, 0)
	for _, v := range collors {
		menuA, err := val.CheckAnyData("menuA do site", 20, v, true)
		if err != nil {
			return err
		}
		menuAsVal = append(menuAsVal, menuA.(string))
	}
	bannerURLVal, err := val.CheckAnyData("banner url", 255, bannerURL, true)
	if err != nil {
		return err
	}
	configsEntity := new(model.Configs)
	configsEntity.Collors = collorsVal
	configsEntity.Links = linksVal
	configsEntity.MenuAs = menuAsVal
	configsEntity.BannerURL = bannerURLVal.(string)
	repConfigs := repository.NewConfigsRepository()
	err = repConfigs.Store(configsEntity)
	if err != nil {
		return err
	}

	return err
}

func (s *configsServiceImpl) List(offset, limit, page int) ([]model.Configs, int, error) {
	if s.kindID != "adm" {
		return nil, 0, errors.New(messages.AdmMessage)
	}

	repConfigs := repository.NewConfigsRepository()
	configs, err := repConfigs.List(offset, limit, page)
	if err != nil {
		return nil, 0, err
	}
	count, err := repConfigs.Count()
	if err != nil {
		return nil, 0, err
	}

	return configs, count, nil

}

func (s *configsServiceImpl) Find(configID int) (*model.Configs, error) {
	if s.kindID != "adm" {
		return nil, errors.New(messages.AdmMessage)
	}
	repConfigs := repository.NewConfigsRepository()
	config, err := repConfigs.Find(configID)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *configsServiceImpl) Update(id int, collors, links, menuAs []string, bannerURL string) error {
	if s.kindID != "adm" {
		return errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()

	// idVal, err := val.CheckAnyData("id de config", 100, id, true)
	// if err != nil {
	// 	return err
	// }

	collorsVal := make([]string, 0)
	for _, v := range collors {
		collor, err := val.CheckAnyData("cor do site", 16, v, true)
		if err != nil {
			return err
		}
		collorsVal = append(collorsVal, collor.(string))
	}

	linksVal := make([]string, 0)
	for _, v := range links {
		link, err := val.CheckAnyData("link do site", 255, v, true)
		if err != nil {
			return err
		}
		linksVal = append(linksVal, link.(string))
	}

	menuAsVal := make([]string, 0)
	for _, v := range collors {
		menuA, err := val.CheckAnyData("menuA do site", 20, v, true)
		if err != nil {
			return err
		}
		menuAsVal = append(menuAsVal, menuA.(string))
	}

	bannerURLVal, err := val.CheckAnyData("banner url", 255, bannerURL, true)
	if err != nil {
		return err
	}

	configsEntity := new(model.Configs)
	configsEntity.ConfigID = uint(id)
	configsEntity.Collors = collorsVal
	configsEntity.Links = linksVal
	configsEntity.MenuAs = menuAsVal
	configsEntity.BannerURL = bannerURLVal.(string)

	repConfigs := repository.NewConfigsRepository()
	err = repConfigs.Update(configsEntity)
	if err != nil {
		return err
	}

	return err
}

func (s *configsServiceImpl) Remove(id int) error {
	if s.kindID != "adm" {
		return errors.New(messages.AdmMessage)
	}
	repConfigs := repository.NewConfigsRepository()
	err := repConfigs.Remove(id)
	if err != nil {
		return err
	}

	return nil
}

func NewConfigsService(userID, kind string) configsServiceInterface {
	return &configsServiceImpl{
		userID: userID,
		kindID: kind,
	}
}
