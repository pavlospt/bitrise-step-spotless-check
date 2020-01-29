package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/bitrise-step-android-unit-test/cache"
	"github.com/bitrise-tools/go-android/gradle"
	"github.com/bitrise-tools/go-steputils/stepconf"
	shellquote "github.com/kballard/go-shellquote"
)

// Config ...
type Config struct {
	ProjectLocation   string `env:"project_location,dir"`
	Module            string `env:"module"`
	Arguments         string `env:"arguments"`
	CacheLevel        string `env:"cache_level,opt[none,only_deps,all]"`
	DeployDir         string `env:"BITRISE_DEPLOY_DIR,dir"`
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

func runSpotlessCheckTask(config Config) error {
	gradleProject, err := gradle.NewProject(config.ProjectLocation)
	if err != nil {
		return fmt.Errorf("Failed to open project, error: %s", err)
	}

	spotlessCheckTask := gradleProject.GetTask("spotlessCheck")

	args, err := shellquote.Split(config.Arguments)
	if err != nil {
		return fmt.Errorf("Failed to parse arguments, error: %s", err)
	}

	log.Infof("Run Spotless check")

	emptyVariants := gradle.Variants{}
	emptyVariants[config.Module] = append(emptyVariants[config.Module], "")

	spotlessCheckCommand := spotlessCheckTask.GetCommand(emptyVariants, args...)
	fmt.Println()
	log.Donef("Printable command args: %s" + spotlessCheckCommand.PrintableCommandArgs())
	fmt.Println()

	taskError := spotlessCheckCommand.Run()
	if taskError != nil {
		log.Errorf("Spotless check task failed, error: %v", taskError)
	}

	return taskError
}

func main() {
	var config Config

	if err := stepconf.Parse(&config); err != nil {
		failf("Couldn't create step config: %v\n", err)
	}

	stepconf.Print(config)
	fmt.Println()

	if err := runSpotlessCheckTask(config); err != nil {
		failf("%s", err)
	}

	fmt.Println()
	log.Infof("Collecting cache:")
	if warning := cache.Collect(config.ProjectLocation, cache.Level(config.CacheLevel)); warning != nil {
		log.Warnf("%s", warning)
	}
	log.Donef("  Done")
}
