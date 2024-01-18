package enrich

import (
	"fio/internal/person"
	"github.com/sirupsen/logrus"
	"sync"
)

// AgeService возвращает возраст на основе имени
type AgeService interface {
	GetAge(name string) (int, error)
}

// NationalityService возвращает национальность на основе имени
type NationalityService interface {
	GetNationality(name string) (string, error)
}

// GenderService возвращает гендер на основе имени
type GenderService interface {
	GetGender(name string) (string, error)
}

// Service занимается обогащением сущности, используя соответствующие интерфейсы api
type Service struct {
	Logger             *logrus.Entry
	NationalityService NationalityService
	AgeService         AgeService
	GenderService      GenderService
}

func (m *Service) setAgeFromService(wg *sync.WaitGroup, person *person.Person) {
	defer wg.Done()

	age, err := m.AgeService.GetAge(person.Name)

	if err != nil {
		m.Logger.Warnf("get age err: %v", err)
		return
	}

	m.Logger.Debugf("get age for %s: %d", person.Name, age)

	person.Age = age
}

func (m *Service) setGenderFromService(wg *sync.WaitGroup, person *person.Person) {
	defer wg.Done()

	gender, err := m.GenderService.GetGender(person.Name)

	if err != nil {
		m.Logger.Warnf("get gender err: %v", err)
		return
	}

	m.Logger.Debugf("get gender for %s: %s", person.Name, gender)

	person.Gender = gender
}

func (m *Service) setNationalityFromService(wg *sync.WaitGroup, person *person.Person) {
	defer wg.Done()

	nationality, err := m.NationalityService.GetNationality(person.Name)

	if err != nil {
		m.Logger.Warnf("get nationality err: %v", err)
		return
	}

	m.Logger.Debugf("get nationality for %s: %s", person.Name, nationality)

	person.Nationality = nationality
}

// EnrichPerson запрашивает дополнительные данные из сторонних api
// Для ускорения обогащения - запускаем в разных горутинах методы записи в переменную
// и ждем выполнения всех горутин
func (m *Service) EnrichPerson(person *person.Person) {

	wg := sync.WaitGroup{}

	wg.Add(3)

	go m.setAgeFromService(&wg, person)
	go m.setGenderFromService(&wg, person)
	go m.setNationalityFromService(&wg, person)

	wg.Wait()

	// todo: сделать обработку ошибок из горутин, для принятия решения о передаче сущности в репозиторий

}
