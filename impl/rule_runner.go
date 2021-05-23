package impl

import (
	"fmt"
	"github.com/fdaines/arch-go/config"
	"github.com/fdaines/arch-go/impl/contents"
	"github.com/fdaines/arch-go/impl/cycles"
	"github.com/fdaines/arch-go/impl/dependencies"
	"github.com/fdaines/arch-go/impl/functions"
	"github.com/fdaines/arch-go/impl/model"
	baseModel "github.com/fdaines/arch-go/model"
	"github.com/fdaines/arch-go/utils"
	"github.com/fdaines/arch-go/utils/packages"
	"github.com/fdaines/arch-go/utils/text"
	"os"
	"regexp"
)

func CheckArchitecture() bool {
	returnValue := true
	utils.ExecuteWithTimer(func() {
		configuration, err := config.LoadConfig("arch-go.yml")
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			os.Exit(1)
		} else {
			mainPackage, _ := packages.GetMainPackage()
			pkgs, _ := packages.GetBasicPackagesInfo()
			moduleInfo := &baseModel.ModuleInfo{
				MainPackage: mainPackage,
				Packages:    pkgs,
			}

			verifications := resolveVerifications(configuration, moduleInfo)
			for _, v := range verifications {
				v.Verify()
			}
			for _, v := range verifications {
				v.PrintResults()
			}
		}
	})
	return returnValue
}

func resolveVerifications(configuration *config.Config, moduleInfo *baseModel.ModuleInfo) []model.RuleVerification {
	var verifications []model.RuleVerification
	for _, dependencyRule := range configuration.DependenciesRules {
		verificationInstance := dependencies.NewDependencyRuleVerification(moduleInfo.MainPackage, dependencyRule)
		packageRegExp, _ := regexp.Compile(text.PreparePackageRegexp(dependencyRule.Package))
		for _, pkg := range moduleInfo.Packages {
			if packageRegExp.MatchString(pkg.Path) {
				verificationInstance.PackageDetails = append(verificationInstance.PackageDetails, model.PackageVerification{
					Package: pkg,
					Passes:  false,
				})
			}
		}
		verifications = append(verifications, verificationInstance)
	}
	for _, functionRule := range configuration.FunctionsRules {
		verificationInstance := functions.NewFunctionsRuleVerification(moduleInfo.MainPackage, functionRule)
		packageRegExp, _ := regexp.Compile(text.PreparePackageRegexp(functionRule.Package))
		for _, pkg := range moduleInfo.Packages {
			if packageRegExp.MatchString(pkg.Path) {
				verificationInstance.PackageDetails = append(verificationInstance.PackageDetails, model.PackageVerification{
					Package: pkg,
					Passes:  false,
				})
			}
		}
		verifications = append(verifications, verificationInstance)
	}
	for _, contentRule := range configuration.ContentRules {
		verificationInstance := contents.NewContentsRuleVerification(moduleInfo.MainPackage, contentRule)
		packageRegExp, _ := regexp.Compile(text.PreparePackageRegexp(contentRule.Package))
		for _, pkg := range moduleInfo.Packages {
			if packageRegExp.MatchString(pkg.Path) {
				verificationInstance.PackageDetails = append(verificationInstance.PackageDetails, model.PackageVerification{
					Package: pkg,
					Passes:  false,
				})
			}
		}
		verifications = append(verifications, verificationInstance)
	}
	for _, cycleRule := range configuration.CyclesRules {
		verificationInstance := cycles.NewCyclesRuleVerification(moduleInfo.MainPackage, moduleInfo.Packages, cycleRule)
		packageRegExp, _ := regexp.Compile(text.PreparePackageRegexp(cycleRule.Package))
		for _, pkg := range moduleInfo.Packages {
			if packageRegExp.MatchString(pkg.Path) {
				verificationInstance.PackageDetails = append(verificationInstance.PackageDetails, model.PackageVerification{
					Package: pkg,
					Passes:  false,
				})
			}
		}
		verifications = append(verifications, verificationInstance)
	}

	return verifications
}
