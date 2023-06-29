/*
Copyright © 2023 Simon Davis simon@breathingdust.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Will output %s", outputFormat))
		// iterate through aws go sdk package names
		// for each one check to see if we have that package imported and note that status
		//   output status and percentage of methods used.
		fmt.Println("Finding references...")

		fmt.Println("Initializing Clients")
		clients := initClients()
		fmt.Printf("Initialized %d service clients\n", len(clients))

		var coverageResult = CoverageResult{}

		for i := 0; i < len(clients); i += 1 {
			coverageResult.Services = append(coverageResult.Services, processService(clients[i]))
		}

		coverageResult.TotalServices = len(coverageResult.Services)

		for i := 0; i < len(coverageResult.Services); i += 1 {
			if coverageResult.Services[i].CreateCoverage == 100 {
				coverageResult.TotalServiceCoverage++
			} else if coverageResult.Services[i].CreateCoverage == 0 {
				coverageResult.NoServiceCoverage++
			} else {
				coverageResult.PartialServiceCoverage++
			}
		}

		pwd, _ := os.Getwd()
		// "cmd/html.tmpl"
		tmplFile := filepath.Join(pwd, fmt.Sprintf("cmd/%s.tmpl", outputFormat))

		tmpl, err := template.ParseFiles(tmplFile)
		if err != nil {
			panic(err)
		}

		var f *os.File

		extension := "html"
		if outputFormat == "markdown" {
			extension = "md"
		}

		f, err = os.Create(fmt.Sprintf("coverage.%s", extension))
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(f, coverageResult)
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	},
}

type ServiceResult struct {
	Name           string
	CreateMethods  map[string]int
	CreateCoverage float32
}

type CoverageResult struct {
	TotalServices          int
	TotalServiceCoverage   int
	PartialServiceCoverage int
	NoServiceCoverage      int
	Services               []ServiceResult
}

func (sr ServiceResult) CalculateCoverage() float32 {
	if len(sr.CreateMethods) == 0 {
		return 0
	}
	coverage := 0
	for _, element := range sr.CreateMethods {
		if element > 0 {
			coverage++
		}
	}
	return (float32(coverage) / float32(len(sr.CreateMethods))) * 100
}

func processService(client interface{}) ServiceResult {

	t := reflect.TypeOf(client)
	pkgName := t.PkgPath()
	sr := ServiceResult{Name: pkgName[strings.LastIndex(pkgName, "/")+1:]}
	fmt.Printf("Examining %s package\n", pkgName)
	createMethods := getCreateMethods(client)
	sr.CreateMethods = make(map[string]int)

	for i := 0; i < len(createMethods); i += 1 {
		sr.CreateMethods[createMethods[i]] = 0
	}

	// Is AWS provider using V1 or V2

	err := filepath.Walk("../terraform-provider-aws/internal/service",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// TODO make this a glob
			if filepath.Ext(path) == ".go" {
				f, err := os.Open(path)
				if err != nil {
					log.Printf("Error loading %s", err)
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)

				line := 1
				// https://golang.org/pkg/bufio/#Scanner.Scan
				for scanner.Scan() {
					for i := 0; i < len(createMethods); i += 1 {
						if strings.Contains(scanner.Text(), createMethods[i]) {
							sr.CreateMethods[createMethods[i]] += 1
						}
					}
					line++
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	sr.CreateCoverage = sr.CalculateCoverage()
	return sr
}

func getCreateMethods(input interface{}) []string {
	t := reflect.TypeOf(input)
	var methods []string
	ptrFooType := reflect.PtrTo(t)
	for i := 0; i < ptrFooType.NumMethod(); i++ {
		method := ptrFooType.Method(i)
		if strings.HasPrefix(method.Name, "Create") {
			methods = append(methods, fmt.Sprintf("conn.%s", method.Name))
		}
	}
	return methods
}

var outputFormat string

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Format of output")
}
