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

	stateFileFlag := flag.String("state", "", "path to a tfState.json file")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("%s [FLAGS...]\n", os.Args[0])
		fmt.Println("Running the program with no args would let you pull the tfState from remote.")
		fmt.Println("It would also let you do state migration operations.")
		fmt.Println("FLAGS:")
		flag.PrintDefaults()
	}

	flag.Parse()

	var err error
	var producers *map[string]schema.AivenProducerModule
	var consumers *map[string]schema.AivenConsumerModule
	var allConsumers *map[string]schema.AivenConsumerModule

	if len(*stateFileFlag) > 0 {
		if producers, consumers, err = parser.ParseAivenStateFile(*stateFileFlag); err != nil {
			panic(err)
		}
		prettyprinter.PrintAivenResources(producers, consumers)

		allConsumers = getAllConsumers(producers, consumers)

		if getApproval("Do you want to generate a new module?") {
			dir := getPathForNewModule()
			factory.BootstrapAivenModule(allConsumers, dir)
		}
		println("\nBye :-)")
		return
	}

	println("Checking AWS Credentials...")
	if parsed, err := shellexec.AwsGetCallerIdentity(); err != nil {
		panic(err)
	} else {
		fmt.Printf("AWS Account: %s\n", parsed.Account)
	}

	toolName := getIacToolName()
	path := getProjectPath(toolName)

	tool := getIaCTool(toolName, path)

	fmt.Printf("Initialazing %s...\n", tool.Name())
	if err := tool.Init(); err != nil {
		panic(err)
	}
	var bytes *[]byte
	if bytes, err = tool.StatePull(); err != nil {
		panic(err)
	}

	if producers, consumers, err = parser.ParseAivenStateBytes(bytes); err != nil {
		panic(err)
	}

	prettyprinter.PrintAivenResources(producers, consumers)

	allConsumers = getAllConsumers(producers, consumers)

	if getApproval("Do you want to generate a new module?") {
		dir := getPathForNewModule()
		factory.BootstrapAivenModule(allConsumers, dir)
	}

	println()
	if getApproval("Do you want to remove aiven PeeringConnection Resources from the state?") {
		removeProducerStates(producers, &tool)
	}

	fmt.Println("\nDo you want to move aiven ConnectionAccepter Resources to a new place?")
	fmt.Println("Attention!!!\nYou must have the new module adopted and ready for this step.")
	if getApproval("") {
		moveConsumerStates(allConsumers, &tool)
	}
	println("\nBye :-)")
}

func getApproval(msg string) bool {
	answer := goprompt.Input(
		fmt.Sprintf("%s [y/n]: ", msg),
		func(d goprompt.Document) []goprompt.Suggest {
			s := []goprompt.Suggest{
				{Text: "yes"},
				{Text: "no"},
			}
			return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		},
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)

	if strings.HasPrefix(answer, "y") {
		return true
	}
	return false
}

func getIacToolName() string {
	tool := goprompt.Input(
		"Select IaC tool [terraform]: ",
		func(d goprompt.Document) []goprompt.Suggest {
			s := []goprompt.Suggest{
				{Text: "terraform"},
				{Text: "terragrunt"},
			}
			return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		},
		goprompt.OptionShowCompletionAtStart(),
		goprompt.OptionCompletionOnDown(),
	)

	if len(tool) == 0 || tool != "terragrunt" {
		tool = "terraform"
	}
	return tool
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
		"What is the new module name? [peering_connection]: ",
		func(d goprompt.Document) []goprompt.Suggest {
			return []goprompt.Suggest{}
		},
	)

	if len(answer) == 0 {
		return "peering_connection"
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

func getIaCTool(name string, path string) shellexec.IaC {
	if name == "terragrunt" {
		return shellexec.NewTerragrunt(path)
	}
	return shellexec.NewTerraform(path)
}

func removeProducerStates(producers *map[string]schema.AivenProducerModule, tool *shellexec.IaC) {
	fmt.Println("Trying to remove ConnectionRequester Resources from the state:")
	dryRun := dryRunMode()
	for _, producer := range *producers {
		address := fmt.Sprintf("%s.%s", producer.Name, producer.PeeringConnection.Address)
		if err := (*tool).StateRemove(address, dryRun); err != nil {
			panic(err)
		}
	}
}

func moveConsumerStates(consumers *map[string]schema.AivenConsumerModule, tool *shellexec.IaC) {
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
	return fmt.Sprintf("modue.%s[%q]", prefix, key)
}

func getAllConsumers(producers *map[string]schema.AivenProducerModule, consumers *map[string]schema.AivenConsumerModule) *map[string]schema.AivenConsumerModule {
	allConsumers := make(map[string]schema.AivenConsumerModule)
	for key, c := range *consumers {
		allConsumers[key] = c
	}
	for _, p := range *producers {
		if p.Consumer != nil {
			allConsumers[p.Consumer.Name] = *p.Consumer
		}
	}
	return &allConsumers
}
