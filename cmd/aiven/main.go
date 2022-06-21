package main

import (
	"flag"
	"fmt"
	goprompt "github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/forando/refactory/pkg/factory"
	"github.com/forando/refactory/pkg/parser"
	"github.com/forando/refactory/pkg/prettyprinter"
	"github.com/forando/refactory/pkg/schema"
	"github.com/forando/refactory/pkg/shellexec"
	"os"
	"strings"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("%s [path/to/tfstate1.json path/to/tfstate2.json...]\n", os.Args[0])
		fmt.Println("Running the program with no args would let you pull the tfState from remote.")
		fmt.Println("It would also let you move resources withing the same state.")
	}

	flag.Parse()

	states := flag.Args()

	if len(states) > 0 {
		goOfflinePath(states)
	} else {
		goOnlinePath()
	}
	println("\nBye :-)")
}

func goOfflinePath(states []string) {
	var err error
	var allConsumers []schema.AivenConsumerModule

	for _, state := range states {
		var producers *map[string]schema.AivenProducerModule
		var consumers *map[string]schema.AivenConsumerModule
		if producers, consumers, err = parser.ParseAivenStateFile(state); err != nil {
			panic(err)
		}
		fmt.Printf("File %s:\n", state)
		prettyprinter.PrintAivenResources(producers, consumers)
		allConsumers = append(allConsumers, *getAllConsumers(producers, consumers)...)
	}

	toolName := getIacToolName()
	if getApproval("Do you want to generate a new module?") {
		newModulePath := getPathForNewModule()
		aivenFactory := factory.NewAivenTerraform(newModulePath)
		if err := aivenFactory.BootstrapNewModule(&allConsumers); err != nil {
			panic(err)
		}
	}
	var runner *shellexec.CmdRunner

	if len(allConsumers) > 0 {
		if getApproval("Do you want to import ConnectionAccepter resources?") {
			projectPath := getProjectPath(toolName)
			if runner == nil {
				r := shellexec.GetCmdRunner(toolName, projectPath)
				runner = &r
			}
			importConsumerStates(&allConsumers, runner)
		}
	}
}

func goOnlinePath() {
	var err error
	var producers *map[string]schema.AivenProducerModule
	var consumers *map[string]schema.AivenConsumerModule
	var allConsumers []schema.AivenConsumerModule

	println("Checking AWS Credentials...")
	if parsed, err := shellexec.AwsGetCallerIdentity(); err != nil {
		panic(err)
	} else {
		fmt.Printf("AWS Account: %s\n", parsed.Account)
	}

	toolName := getIacToolName()
	projectPath := getProjectPath(toolName)

	runner := shellexec.GetCmdRunner(toolName, projectPath)

	if err := runner.Init(); err != nil {
		panic(err)
	}
	var bytes *[]byte
	if bytes, err = runner.StatePull(); err != nil {
		panic(err)
	}

	if producers, consumers, err = parser.ParseAivenStateBytes(bytes); err != nil {
		panic(err)
	}

	prettyprinter.PrintAivenResources(producers, consumers)

	allConsumers = *getAllConsumers(producers, consumers)

	if getApproval("Do you want to generate a new module?") {
		newModulePath := getPathForNewModule()
		aivenFactory := factory.NewAivenTerraform(newModulePath)
		if err := aivenFactory.BootstrapNewModule(&allConsumers); err != nil {
			panic(err)
		}
	}

	println()
	if len(*producers) > 0 {
		if getApproval("Do you want to remove aiven PeeringConnection Resources from the state?") {
			removeProducerStates(producers, &runner)
		}
	}

	if len(allConsumers) > 0 {
		fmt.Println("\nDo you want to move aiven ConnectionAccepter Resources to a new place?")
		fmt.Println("Attention!!!\nYou must have the new module adopted and ready for this step.")
		if getApproval("") {
			moveConsumerStates(&allConsumers, &runner)
		}
	}
}

func getApproval(msg string) bool {
	answer := goprompt.Input(
		fmt.Sprintf("%s [y/n]: ", msg),
		func(d goprompt.Document) []goprompt.Suggest {
			s := []goprompt.Suggest{
				{Text: "yes"},
				{Text: "no"},
				{Text: "exit"},
			}
			return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		},
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)

	if answer == "exit" {
		println("\nBye :-)")
		os.Exit(0)
	}
	if strings.HasPrefix(answer, "y") {
		return true
	}
	return false
}

func getIacToolName() string {
	answer := goprompt.Input(
		"Select IaC tool [terraform]: ",
		func(d goprompt.Document) []goprompt.Suggest {
			s := []goprompt.Suggest{
				{Text: "terraform"},
				{Text: "terragrunt"},
				{Text: "exit"},
			}
			return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		},
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)

	if len(answer) == 0 {
		answer = "terraform"
	}
	if answer == "exit" {
		println("\nBye :-)")
		os.Exit(0)
	}
	if answer != "terragrunt" && answer != "terraform" {
		answer = "terraform"
	}
	return answer
}

func getProjectPath(tool string) string {
	path := goprompt.Input(
		fmt.Sprintf("directory to run %s from [.]: ", tool),
		func(d goprompt.Document) []goprompt.Suggest {
			return dirCompleter.Complete(d)
		},
		goprompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
		goprompt.OptionInputTextColor(goprompt.Cyan),
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)
	if len(path) == 0 {
		path = "."
	}

	return path
}

func getPathForNewModule() string {
	path := goprompt.Input(
		fmt.Sprint("directory where to generate a new module [.]: "),
		func(d goprompt.Document) []goprompt.Suggest {
			return dirCompleter.Complete(d)
		},
		goprompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
		goprompt.OptionInputTextColor(goprompt.Cyan),
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)
	if len(path) == 0 {
		path = "."
	}

	return path
}

func dryRunMode() bool {
	tool := goprompt.Input(
		"Select run mode [dry-run]: ",
		func(d goprompt.Document) []goprompt.Suggest {
			s := []goprompt.Suggest{
				{Text: "dry-run", Description: "Only tests the request, but does not change the state"},
				{Text: "hit-it", Description: "Execute the request mutates the state"},
			}
			return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		},
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)

	if len(tool) == 0 || tool != "hit-it" {
		return true
	}
	return false
}

func getModuleName() string {
	answer := goprompt.Input(
		"What is the new module name? [peering_connections]: ",
		func(d goprompt.Document) []goprompt.Suggest {
			return []goprompt.Suggest{}
		},
	)

	if len(answer) == 0 {
		return "peering_connections"
	}
	return answer
}

var dirCompleter = completer.FilePathCompleter{
	IgnoreCase: true,
	Filter: func(fi os.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		return false
	},
}

func removeProducerStates(producers *map[string]schema.AivenProducerModule, tool *shellexec.CmdRunner) {
	fmt.Println("Trying to remove ConnectionRequester Resources from the state:")
	dryRun := dryRunMode()
	for _, producer := range *producers {
		address := fmt.Sprintf("%s.%s", producer.Name, producer.PeeringConnection.Address)
		if err := (*tool).StateRemove(address, dryRun); err != nil {
			panic(err)
		}
	}
}

func importConsumerStates(consumers *[]schema.AivenConsumerModule, tool *shellexec.CmdRunner) {
	destPrefix := getModuleName()
	importedConsumers := make([]*schema.AivenConsumerModule, 0)
	rollBack := false
	for _, consumer := range *consumers {
		moduleName := buildNewModuleName(&consumer, destPrefix)
		accepterAddr := fmt.Sprintf("%s.%s", moduleName, consumer.ConnectionAccepter.Address)
		if err := (*tool).StateImport(accepterAddr, consumer.ConnectionAccepter.Id); err != nil {
			fmt.Println(err.Error())
			rollBack = true
			break
		}

		aclAddr := fmt.Sprintf("%s.%s", moduleName, consumer.AwsNetworkAclRules[schema.IngressUdp].Address)
		if err := (*tool).StateImport(aclAddr, consumer.AwsNetworkAclRules[schema.IngressUdp].Id); err != nil {
			fmt.Println(err.Error())
			(*tool).StateRemove(accepterAddr, false)
			rollBack = true
			break
		}

		importedRoutes := make([]*schema.AwsRouteResource, 0)
		for _, route := range consumer.AwsRoutResources {
			if err := (*tool).StateImport(fmt.Sprintf("%s.%s", moduleName, route.Address), route.Id); err != nil {
				fmt.Println(err.Error())
				rollBack = true
				break
			}
			importedRoutes = append(importedRoutes, &route)
		}
		if !rollBack {
			importedConsumers = append(importedConsumers, &consumer)
		} else {
			fmt.Println("ROLLING BACK...")
			(*tool).StateRemove(accepterAddr, false)
			(*tool).StateRemove(aclAddr, false)
			for _, route := range importedRoutes {
				(*tool).StateRemove(fmt.Sprintf("%s.%s", moduleName, route.Address), false)
			}
			break
		}
	}
	if rollBack {
		for _, consumer := range importedConsumers {
			moduleName := buildNewModuleName(consumer, destPrefix)
			(*tool).StateRemove(fmt.Sprintf("%s.%s", moduleName, consumer.ConnectionAccepter.Address), false)
			(*tool).StateRemove(fmt.Sprintf("%s.%s", moduleName, consumer.AwsNetworkAclRules[schema.IngressUdp].Address), false)
			for _, route := range consumer.AwsRoutResources {
				(*tool).StateRemove(fmt.Sprintf("%s.%s", moduleName, route.Address), false)
			}
		}

	}
}

func moveConsumerStates(consumers *[]schema.AivenConsumerModule, tool *shellexec.CmdRunner) {
	fmt.Println("Trying to move ConnectionAcceptor Resources to a new location:")
	var src, dest string
	dryRun := dryRunMode()
	destPrefix := getModuleName()
	for _, consumer := range *consumers {
		moduleName := buildNewModuleName(&consumer, destPrefix)
		src = fmt.Sprintf("%s.%s", consumer.Name, consumer.ConnectionAccepter.Address)
		dest = fmt.Sprintf("%s.%s", moduleName, consumer.ConnectionAccepter.Address)
		if err := (*tool).StateMove(src, dest, dryRun); err != nil {
			panic(err)
		}

		src = fmt.Sprintf("%s.%s", consumer.Name, consumer.AwsNetworkAclRules[schema.IngressTcp].Address)
		dest = fmt.Sprintf("%s.%s", moduleName, consumer.AwsNetworkAclRules[schema.IngressTcp].Address)
		if err := (*tool).StateMove(src, dest, dryRun); err != nil {
			panic(err)
		}

		src = fmt.Sprintf("%s.%s", consumer.Name, consumer.AwsNetworkAclRules[schema.IngressUdp].Address)
		dest = fmt.Sprintf("%s.%s", moduleName, consumer.AwsNetworkAclRules[schema.IngressUdp].Address)
		if err := (*tool).StateMove(src, dest, dryRun); err != nil {
			panic(err)
		}

		for _, route := range consumer.AwsRoutResources {
			src = fmt.Sprintf("%s.%s", consumer.Name, route.Address)
			dest = fmt.Sprintf("%s.%s", moduleName, route.Address)
			if err := (*tool).StateMove(src, dest, dryRun); err != nil {
				panic(err)
			}
		}
	}
}

func buildNewModuleName(consumer *schema.AivenConsumerModule, prefix string) string {
	key := fmt.Sprintf("%s/%s", consumer.ConnectionAccepter.VpcId, consumer.ConnectionAccepter.PeeringConnectionId)
	if strings.HasPrefix(prefix, "module") {
		return fmt.Sprintf("%s[%q]", prefix, key)
	}
	return fmt.Sprintf("module.%s[%q]", prefix, key)
}

func getAllConsumers(producers *map[string]schema.AivenProducerModule, consumers *map[string]schema.AivenConsumerModule) *[]schema.AivenConsumerModule {
	allConsumers := make([]schema.AivenConsumerModule, 0)
	for _, c := range *consumers {
		allConsumers = append(allConsumers, c)
	}
	for _, p := range *producers {
		if p.Consumer != nil {
			allConsumers = append(allConsumers, *p.Consumer)
		}
	}
	return &allConsumers
}
