package resourceimport

import (
	"fmt"
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/parser"
	"github.com/forando/refactory/pkg/schema"
	"github.com/forando/refactory/pkg/shellexec"
	"path/filepath"
	"time"
)

type merge struct {
	resources chan schema.ResourceCount
	done      chan schema.Done
	result    chan schema.Result
	output    chan schema.Done
}

// We have to limit max GoRoutines running in parallel
// Otherwise `timeout` error could happen
const maxGoroutines = 8

func Start(names []string) (<-chan schema.Done, *[]schema.ResourceCount) {
	m := &merge{
		resources: make(chan schema.ResourceCount),
		done:      make(chan schema.Done),
		result:    make(chan schema.Result),
		output:    make(chan schema.Done),
	}
	delay := 10000
	guard := make(chan struct{}, maxGoroutines)
	for index, name := range names {
		if index == 0 {
			go importResources(name, m.resources, m.result, m.done, 0, guard)
		} else {
			// give aws some time to create s3 bucket and mongoDB table
			// need it only if using shellexec.ExecTerragruntInitWithStdIn
			go importResources(name, m.resources, m.result, m.done, delay, guard)
		}
	}

	rest := make(map[string]interface{})
	dirsToProcess := make(map[string]interface{})
	for _, name := range names {
		rest[name] = nil
		dirsToProcess[name] = nil
	}

	resources := 0

	emptyResources := make([]schema.ResourceCount, 0)

	for res := range m.resources {
		if res.Total == 0 {
			emptyResources = append(emptyResources, res)
			delete(dirsToProcess, res.Dir)
		} else {
			resources += res.Total
		}
		delete(rest, res.Dir)
		if len(rest) == 0 {
			close(m.resources)
			break
		}
	}

	go processMessages(dirsToProcess, m, resources, guard)
	return m.output, &emptyResources
}

func processMessages(names map[string]interface{}, m *merge, totalResources int, guard chan struct{}) {
	okResults := 0
	errorResults := 0
	totalDirs := len(names)
	processedDirs := 0
	var kickedInDirs int
	if totalDirs > maxGoroutines {
		kickedInDirs = maxGoroutines
	} else {
		kickedInDirs = totalDirs
	}

	if len(names) == 0 {
		close(m.done)
		close(m.result)
		close(m.output)
		return
	}

	for {
		select {
		case res := <-m.result:
			if res.GetStatus() == schema.Ok {
				okResults++
			} else {
				errorResults++
			}
			fmt.Printf("\033[2K\rDirs: %d/%d Resources out of %d: ok %d errors %d", totalDirs, kickedInDirs, totalResources, okResults, errorResults)
		case done := <-m.done:
			delete(names, done.Dir)
			m.output <- done

			if totalDirs-processedDirs > maxGoroutines {
				<-guard
				kickedInDirs++
			}
			processedDirs++
			if len(names) == 0 {
				fmt.Println("")
				close(m.done)
				close(m.result)
				close(m.output)
				close(guard)
				return
			}
		}
	}
}

const ImportsFileName string = "imports.csv"

func importResources(dir string, resources chan<- schema.ResourceCount, res chan<- schema.Result, done chan<- schema.Done, delay int, guard chan<- struct{}) {
	imports, err := parser.ParseImports(filepath.Join(dir, ImportsFileName))
	if err != nil {
		resources <- schema.ResourceCount{Dir: dir, Total: 0, Message: err.Error()}
		return
	}

	if imports == nil || len(*imports) == 0 {
		resources <- schema.ResourceCount{Dir: dir, Total: 0, Message: fmt.Sprintf("%s contains no resources", ImportsFileName)}
		return
	}

	resourceCount := len(*imports)

	resources <- schema.ResourceCount{Dir: dir, Total: resourceCount}

	time.Sleep(time.Duration(delay) * time.Millisecond)

	guard <- struct{}{} //would block if guard channel is already filled

	terragrunt := shellexec.Terragrunt{Dir: dir}

	if err := terragrunt.Init(); err != nil {
		cleanErrors := filesystem.CleanDir(dir)
		res <- &schema.ErrResult{Status: schema.Err, Dir: dir, ResourceCount: resourceCount, Message: err.Error()}
		done <- schema.Done{Status: schema.Err, Dir: dir, FailedResource: &schema.ImportResource{Address: "", Id: ""}, ResourceCount: resourceCount, Message: err.Error(), CleanErrors: cleanErrors}
		return
	}

	doneImports := make([]schema.Import, 0)

	for _, imp := range *imports {
		doneImports = append(doneImports, imp)
		if output, err := terragrunt.Import(imp.Address, imp.Id); err != nil {
			terragrunt.RollBackImports(&doneImports)
			cleanErrors := filesystem.CleanDir(dir)
			res <- &schema.ErrResult{Status: schema.Err, Dir: dir, ResourceCount: resourceCount, Message: err.Error()}
			done <- schema.Done{Status: schema.Err, Dir: dir, FailedResource: &schema.ImportResource{Address: imp.Address, Id: imp.Id}, ResourceCount: resourceCount, Message: err.Error(), CleanErrors: cleanErrors}
			return
		} else {
			res <- &schema.OkResult{Status: schema.Ok, Dir: dir, ResourceCount: resourceCount, Message: *output}
		}
	}
	cleanErrors := filesystem.CleanDir(dir)
	done <- schema.Done{Status: schema.Ok, Dir: dir, ResourceCount: resourceCount, Message: "ok", CleanErrors: cleanErrors}
}
