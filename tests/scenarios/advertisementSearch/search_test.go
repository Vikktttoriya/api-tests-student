package advertisementSearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/fixtures"
	advList "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	photoPath       = "testdata/photo.png"
	titlePartLength = 6
)

type TestSuite struct {
	suite.Suite
	loginData          models.LoginOkResponse
	advertisementTitle string
	searchBeginning    string
	searchMiddle       string
	searchEnd          string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.loginData = base.LoginTestUser(s.T())

	base.Precondition("Создаём объявление с уникальным названием из трёх частей: начало, середина, конец")
	s.searchBeginning = utils.RandomString(titlePartLength)
	s.searchMiddle = utils.RandomString(titlePartLength)
	s.searchEnd = utils.RandomString(titlePartLength)

	request := fixtures.AdvertisementRequest(photoPath)
	request.Title = strings.Join([]string{s.searchBeginning, s.searchMiddle, s.searchEnd}, " ")
	s.advertisementTitle = fixtures.CreateAdvertisement(s.T(), s.loginData.Token, request).Title

	base.Precondition("Ждём, пока объявление появится в поиске")
	list := advList.SearchAdvertisementsUntilFound(s.T(), s.loginData.Token, s.advertisementTitle)
	assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.advertisementTitle)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, s.searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.advertisementTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", s.searchBeginning},
		{"середина", s.searchMiddle},
		{"конец", s.searchEnd},
		{"нижний регистр", strings.ToLower(s.searchMiddle)},
		{"верхний регистр", strings.ToUpper(s.searchMiddle)},
		{"смешанный регистр", utils.MixCase(s.searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.advertisementTitle)
		})
	}
}
