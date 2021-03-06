package user_test

import (
	. "cf/commands/user"
	"cf/configuration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestCreateUserFailsWithUsage", func() {
		defaultArgs, defaultReqs, defaultUserRepo := getCreateUserDefaults()

		ui := callCreateUser([]string{}, defaultReqs, defaultUserRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callCreateUser(defaultArgs, defaultReqs, defaultUserRepo)
		Expect(ui.FailedWithUsage).To(BeFalse())
	})

	It("TestCreateUserRequirements", func() {
		defaultArgs, defaultReqs, defaultUserRepo := getCreateUserDefaults()

		callCreateUser(defaultArgs, defaultReqs, defaultUserRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())

		notLoggedInReq := &testreq.FakeReqFactory{LoginSuccess: false}
		callCreateUser(defaultArgs, notLoggedInReq, defaultUserRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	It("TestCreateUser", func() {
		defaultArgs, defaultReqs, defaultUserRepo := getCreateUserDefaults()
		ui := callCreateUser(defaultArgs, defaultReqs, defaultUserRepo)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Creating user", "my-user", "current-user"},
			{"OK"},
			{"TIP"},
		})
		Expect(defaultUserRepo.CreateUserUsername).To(Equal("my-user"))
	})

	It("TestCreateUserWhenItAlreadyExists", func() {
		defaultArgs, defaultReqs, userAlreadyExistsRepo := getCreateUserDefaults()
		userAlreadyExistsRepo.CreateUserExists = true

		ui := callCreateUser(defaultArgs, defaultReqs, userAlreadyExistsRepo)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Creating user"},
			{"FAILED"},
			{"my-user"},
			{"already exists"},
		})
	})
})

func getCreateUserDefaults() (defaultArgs []string, defaultReqs *testreq.FakeReqFactory, defaultUserRepo *testapi.FakeUserRepository) {
	defaultArgs = []string{"my-user", "my-password"}
	defaultReqs = &testreq.FakeReqFactory{LoginSuccess: true}
	defaultUserRepo = &testapi.FakeUserRepository{}
	return
}

func callCreateUser(args []string, reqFactory *testreq.FakeReqFactory, userRepo *testapi.FakeUserRepository) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("create-user", args)
	configRepo := testconfig.NewRepositoryWithDefaults()
	accessToken, err := testconfig.EncodeAccessToken(configuration.TokenInfo{
		Username: "current-user",
	})
	Expect(err).NotTo(HaveOccurred())
	configRepo.SetAccessToken(accessToken)

	cmd := NewCreateUser(ui, configRepo, userRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
