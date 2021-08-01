package main

import (
	"os"
	"time"

	"ricklepickle/docker/postgres"
	"ricklepickle/internal/app"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

func main() {
	var status int = 0
	defer func() {
		os.Exit(status)
	}()
	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		zap.L().Error("create docker client: %w")
		return
	}
	postgresClient, err := postgres.NewClient(dockerClient)
	if err != nil {
		zap.L().Error("create docker client: %w")
		return
	}
	appClient, err := app.NewClient(dockerClient)
	if err != nil {
		zap.L().Error("create docker client: %w")
		return
	}
	suite := &Suite{
		App:      appClient,
		Postgres: postgresClient,
	}
	status = godog.TestSuite{
		Name:                 "example",
		TestSuiteInitializer: InitializeTestSuite(suite),
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

}

type Suite struct {
	App      *app.Client
	Postgres *postgres.Client
}

func InitializeTestSuite(app *Suite) func(testSuiteContext *godog.TestSuiteContext) {
	return func(testSuiteContext *godog.TestSuiteContext) {

		testSuiteContext.BeforeSuite(func() {
			app.App.NewContainer()
		})

		testSuiteContext.AfterSuite(func() {
			app.App.NewContainer()
		})

	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^application exits with status (\d+)$`, applicationExitsWithStatus)
	ctx.Step(`^Call echo "([^"]*)" returns:$`, callEchoReturns)
	ctx.Step(`^Quick exit signal delivered$`, quickExitSignalDelivered)
}

func applicationExitsWithStatus(arg1 int) error {
	return godog.ErrPending
}

func callEchoReturns(arg1 string, arg2 *messages.PickleStepArgument_PickleDocString) error {
	return godog.ErrPending
}

func quickExitSignalDelivered() error {
	return godog.ErrPending
}
