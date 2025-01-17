package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func resetClients(ctx *cli.Context) error {
	gethRunning := isRunning(gethDependencyName)
	prysmRunning := isRunning(prysmDependencyName)
	validatorRunning := isRunning(validatorDependencyName)

	if gethRunning || prysmRunning || validatorRunning {
		message := "Please stop the following clients before resetting: "
		if gethRunning {
			message += "geth "
		}
		if prysmRunning {
			message += "prysm "
		}
		if validatorRunning {
			message += "validator "
		}

		log.Warn(message)

		return nil
	}

	message := fmt.Sprintf("WARNING: THIS ACTION WILL REMOVE DATA DIRECTORIES FROM ALL OF RUNNING CLIENTS.\n"+
		"Are you sure you want to continue?\nDirectories that will be deleted:\n"+
		"- %s\n- %s\n- %s\n[Y/n]: ", ctx.String(gethDatadirFlag), ctx.String(prysmDatadirFlag), ctx.String(validatorDatadirFlag))

	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") {
		log.Info("Aborting...")

		return nil
	}

	err := resetGeth(ctx)
	if err != nil {
		return err
	}

	err = resetPrysm(ctx)
	if err != nil {
		return err
	}

	err = resetValidator(ctx)

	return err
}

func resetGeth(ctx *cli.Context) error {
	dataDirPath := ctx.String(gethDatadirFlag)
	if dataDirPath == "" {
		return errFlagMissing
	}

	return os.RemoveAll(dataDirPath)
}

func resetPrysm(ctx *cli.Context) error {
	dataDirPath := ctx.String(prysmDatadirFlag)
	if dataDirPath == "" {
		return errFlagMissing
	}

	return os.RemoveAll(dataDirPath)
}

func resetValidator(ctx *cli.Context) error {
	dataDirPath := ctx.String(validatorDatadirFlag)
	if dataDirPath == "" {
		return errFlagMissing
	}

	return os.RemoveAll(dataDirPath)
}
