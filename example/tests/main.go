package main

import (
	"os"
	"time"

	"github.com/7d4b9/lever/example/tests/docker"
	"github.com/7d4b9/lever/example/tests/internal/app"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/docker/docker/client"
)

func main() {
	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	c, _ := client.NewClientWithOpts(client.FromEnv)
	dockerClient, _ := docker.NewClient(c, "")

	app := &App{
		app.NewClient(dockerClient),
	}

	status := godog.TestSuite{
		Name:                 "example",
		TestSuiteInitializer: InitializeTestSuite(app),
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

type App struct {
	App *app.Client
}

func InitializeTestSuite(app *App) func(testSuiteContext *godog.TestSuiteContext) {
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
