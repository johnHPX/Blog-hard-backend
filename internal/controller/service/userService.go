package service

import (
	"github.com/google/uuid"
	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type userServiceInterface interface {
	Store(name, telephone, nick, email, secret, kind string) error
	List(offset, limit, page int) ([]model.User, error)
	Count() (int, error)
	ListName(name string, offset, limit, page int) ([]model.User, error)
	CountName(name string) (int, error)
	Find(id string) (*model.User, error)
	Update(id, name, telefone, nick, email, kind string) error
	Remove(id string) error
}

type userServiceImpl struct{}

func (s *userServiceImpl) Store(name, telephone, nick, email, secret, kind string) error {

	val := validator.NewValidator()
	Name, err := val.CheckAnyData("nome", 255, name, true)
	if err != nil {
		return err
	}
	Telephone, err := val.CheckAnyData("telefone", 13, telephone, true)
	if err != nil {
		return err
	}
	Nick, err := val.CheckAnyData("nick", 255, nick, true)
	if err != nil {
		return err
	}
	Email, err := val.CheckAnyData("email", 255, email, true)
	if err != nil {
		return err
	}
	Password, err := val.CheckPassword(255, secret, "", "create")
	if err != nil {
		return err
	}
	Kind, err := val.CheckAnyData("kind", 10, kind, true)
	if err != nil {
		return err
	}

	uid := uuid.New()
	pid := uuid.New()

	e := &model.User{
		UserID: uid.String(),
		Person: model.Person{
			PersonID:  pid.String(),
			Name:      Name.(string),
			Telephone: Telephone.(string),
		},
		Nick:   Nick.(string),
		Email:  Email.(string),
		Secret: Password,
		Kind:   Kind.(string),
	}

	repUser := repository.NewUserRepository()
	repPerson := repository.NewPersonRepository()

	err = repUser.Store(e)
	if err != nil {
		return err
	}

	err = repPerson.Store(&e.Person, e.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) List(offset, limit, page int) ([]model.User, error) {
	repUser := repository.NewUserRepository()
	entities, err := repUser.List(offset, limit, page)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *userServiceImpl) Count() (int, error) {
	repUser := repository.NewUserRepository()
	count, err := repUser.Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *userServiceImpl) ListName(name string, offset, limit, page int) ([]model.User, error) {
	repUser := repository.NewUserRepository()
	entities, err := repUser.ListName(name, offset, limit, page)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *userServiceImpl) CountName(name string) (int, error) {
	repUser := repository.NewUserRepository()
	count, err := repUser.CountListName(name)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *userServiceImpl) Find(id string) (*model.User, error) {
	repUser := repository.NewUserRepository()
	user, err := repUser.Find(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) Update(id, name, telefone, nick, email, kind string) error {
	val := validator.NewValidator()
	NameVal, err := val.CheckAnyData("nome", 255, name, true)
	if err != nil {
		return err
	}
	TelefoneVal, err := val.CheckAnyData("telefone", 13, telefone, true)
	if err != nil {
		return err
	}
	NickVal, err := val.CheckAnyData("nick", 255, nick, true)
	if err != nil {
		return err
	}
	EmailVal, err := val.CheckAnyData("email", 255, email, true)
	if err != nil {
		return err
	}
	KindVal, err := val.CheckAnyData("kind", 10, kind, true)
	if err != nil {
		return err
	}

	// repositorys
	repUser := repository.NewUserRepository()
	repPerson := repository.NewPersonRepository()

	// find id person
	user, err := repUser.Find(id)
	if err != nil {
		return err
	}

	// Update person
	person := new(model.Person)
	person.PersonID = user.PersonID
	person.Name = NameVal.(string)
	person.Telephone = TelefoneVal.(string)
	err = repPerson.Update(person)
	if err != nil {
		return err
	}

	// Update User
	userEntity := new(model.User)
	userEntity.UserID = user.UserID
	userEntity.Nick = NickVal.(string)
	userEntity.Email = EmailVal.(string)
	userEntity.Kind = KindVal.(string)
	err = repUser.Update(userEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Remove(id string) error {
	// repositorys
	repUser := repository.NewUserRepository()
	repPerson := repository.NewPersonRepository()

	// find id person
	user, err := repUser.Find(id)
	if err != nil {
		return err
	}

	// remove person
	err = repPerson.Remove(user.PersonID)
	if err != nil {
		return err
	}

	// remove user
	err = repUser.Remove(user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func NewUserService() userServiceInterface {
	return &userServiceImpl{}
}
