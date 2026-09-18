package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	model "github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	sourceType      string
	sinkType        string
	resourceID      string
	resourceName    string
	dataFormat      string
	dataCompression string
	integrationID   string // For associating with an integration
	url             string // For HEC sink
	token           string // For HEC sink

	processorName string
	routerName    string
	pipeName      string

	sourceID string // For connecting source to router
	routerID string // For connecting source to router
	sinkID   string // For connecting source to router

)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Manage streams",
}

var connectSinkCmd = &cobra.Command{
	Use:   "connect-sink",
	Short: "Connect a stream router to a sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Connecting stream router to a sink...")

		err := AppAPI.SetRouterSink(routerID, sinkID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream connection added successfully: router ", routerID, " to sink ", sinkID)
		//cmd.Println(resourceID)
		return nil
	},
}

var connectRouterCmd = &cobra.Command{
	Use:   "connect-router",
	Short: "Connect a stream source to a router",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Connecting source to a stream router...")

		err := AppAPI.SetSourceRouter(sourceID, routerID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream connection added successfully: source ", sourceID, " to router ", routerID)
		//cmd.Println(resourceID)
		return nil
	},
}

var updatePipeProcessorCmd = &cobra.Command{
	Use:   "update-pipe-processor",
	Short: "Update the processor for a pipe within a router",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Updating pipe processor...")

		err := AppAPI.UpdatePipeProcessor(routerName, pipeName, processorName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Pipe processor updated successfully: router ", routerName, " pipe ", pipeName, " processor ", processorName)
		return nil
	},
}

var addRouterCmd = &cobra.Command{
	Use:   "add-router",
	Short: "Add a stream router",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream router...")

		resourceID, err := AppAPI.AddSimpleRouter(processorName, routerName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream router added successfully: ", resourceID)
		fmt.Println(resourceID)
		return nil
	},
}

var addSourceCmd = &cobra.Command{
	Use:   "add-source",
	Short: "Add a stream source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream datasource...")

		source := &model.DataSourceConfig{
			Type: sourceType,
			Name: resourceName,
			// Add other necessary fields for DataSourceConfig
			Format: "json", // Example format, adjust as needed
			//Compression
		}

		if sourceType == "plugin" {
			if integrationID == "" {
				return fmt.Errorf("integration-id is required for plugin source type")
			}
			source.Plugin = &model.PluginSourceConfig{
				ID: integrationID,
			}
		} else if sourceType == "hec" {
			source.Hec = &model.HecSourceConfig{
				URL: url,
			}
			s := &model.HecSecret{
				Token: token,
			}
			b, _ := json.Marshal(s)
			source.Secret = b
		}

		response, err := AppAPI.AddDataSource(source)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Stream source added successfully: ", response.ID)
		if response.URL != "" {
			cmd.PrintErrln("Access URL:", response.URL)
		}
		if len(response.Secret) > 0 {
			b, _ := response.Secret.MarshalJSON()
			cmd.PrintErrln("Secret:", string(b))
		}
		fmt.Println(response.ID)
		return nil
	},
}

var delSourceCmd = &cobra.Command{
	Use:   "del-source",
	Short: "Delete a stream source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Deleting stream datasource...")

		err := AppAPI.DeleteDataSource(resourceID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream source deleted successfully: ", resourceID)
		return nil
	},
}

var listSourceCmd = &cobra.Command{
	Use:   "list-source",
	Short: "List all stream sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing stream datasource...")

		entries, err := AppAPI.ListDataSource()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No stream source found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Type: %s\n", entry.ID, entry.Name, entry.Type)
		}
		return nil
	},
}

var delSinkCmd = &cobra.Command{
	Use:   "del-sink",
	Short: "Delete a stream sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Deleting stream datasink...")

		err := AppAPI.DeleteDataSink(resourceID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream sink deleted successfully: ", resourceID)
		return nil
	},
}

var listSinkCmd = &cobra.Command{
	Use:   "list-sink",
	Short: "List all stream sinks",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing stream datasink...")

		entries, err := AppAPI.ListDataSink()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No stream sink found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Type: %s\n", entry.ID, entry.Name, entry.Type)
		}
		return nil
	},
}

// Example leaf command: source
var addSinkCmd = &cobra.Command{
	Use:   "add-sink",
	Short: "Add a stream sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream sink...")
		// Just call the global interface
		sink := &model.DataSinkConfig{
			Type: sinkType,
			Name: resourceName,
			// Add other necessary fields for DataSourceConfig
			//Format: "json", // Example format, adjust as needed
			//Compression
		}

		if sinkType == "datalake" {
			if index == "" {
				return fmt.Errorf("index is required for datalake sink type")
			}
			sink.DataLake = &model.DataLakeSinkConfig{
				Datalake:      datalake,
				DatalakeIndex: index,
			}
		} else if sinkType == "hec" {
			sink.Hec = &model.HecSinkConfig{
				URL:   url,
				Token: token,
				// Add necessary fields for HEC sink
			}
		} else if sinkType == "webhook" {
			sink.Webhook = &model.WebhookSinkConfig{
				URL: url,
				//Token: token,
				// Add necessary fields for HEC sink
			}
		} else if sinkType == "redis" {
			// The queue is what makes the sink: queue:BehaviorSummary:EventQueue
			// is the behavior service, queue:LVDBService:JobQueue the datalake
			// writer (which also needs --datalake/--index).
			if redisQueue == "" {
				return fmt.Errorf("--queue is required for redis sink type")
			}
			sink.Redis = &model.RedisSinkConfig{
				Redis: &model.RedisConfig{
					Host:  "localhost",
					Port:  6379,
					Queue: redisQueue,
				},
			}
			if index != "" {
				sink.Redis.Datalake = datalake
				sink.Redis.DatalakeIndex = index
			}
			// the flush defaults every app-installed sink carries
			sink.FlushCount = 1024
			sink.FlushBuffer = 1048576
			sink.FlushInterval = 60
		}

		tags, err := parseTags(pipeTags)
		if err != nil {
			return err
		}
		sink.Tags = tags

		response, err := AppAPI.AddDataSink(sink)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Stream sink added successfully: ", response.ID)
		fmt.Println(response.ID)

		return nil
	},
}

// ... Repeat for sink, router, connection ...

func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.AddCommand(addSourceCmd, delSourceCmd, listSourceCmd, addSinkCmd, delSinkCmd, listSinkCmd, addRouterCmd, connectRouterCmd, connectSinkCmd, updatePipeProcessorCmd) // Add del/update similarly
	streamCmd.AddCommand(listRouterCmd, addPipeCmd, delPipeCmd, delRouterCmd)

	addSourceCmd.Flags().StringVar(&sourceType, "source-type", "", "data source type: plugin, s3, hec, webhook ")
	addSourceCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSourceCmd.Flags().StringVar(&dataFormat, "format", "json", "Data Format")
	addSourceCmd.Flags().StringVar(&dataCompression, "compression", "", "Data Compression")

	addSourceCmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID")

	_ = addSourceCmd.MarkFlagRequired("source-type")
	_ = addSourceCmd.MarkFlagRequired("name")

	addSinkCmd.Flags().StringVar(&sinkType, "sink-type", "", "data sink type: datalake, redis, hec, webhook, drop")
	addSinkCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSinkCmd.Flags().StringVar(&url, "url", "", "URL for HEC sink")
	addSinkCmd.Flags().StringVar(&token, "token", "", "Token")

	addSinkCmd.Flags().StringVar(&datalake, "datalake", "managed", "datalake name")
	addSinkCmd.Flags().StringVar(&index, "index", "", "datalake index name")
	addSinkCmd.Flags().StringVar(&redisQueue, "queue", "", "redis queue for a redis sink, e.g. queue:BehaviorSummary:EventQueue")
	addSinkCmd.Flags().StringArrayVar(&pipeTags, "tag", nil, "ownership tag name=value (repeatable)")

	addPipeCmd.Flags().StringVar(&routerName, "router", "", "Router name or ID")
	addPipeCmd.Flags().StringVar(&resourceName, "name", "", "Pipe name")
	addPipeCmd.Flags().StringVar(&pipeProcessor, "processor", "", "Processor name (a pipe carries exactly one)")
	addPipeCmd.Flags().StringArrayVar(&pipeSinks, "sink", nil, "Sink name or ID (repeatable)")
	addPipeCmd.Flags().StringArrayVar(&pipeTags, "tag", nil, "ownership tag name=value (repeatable)")
	addPipeCmd.Flags().IntVar(&pipePriority, "priority", 0, "pipe priority within the router (templates: main 1000, behavior 2000)")
	addPipeCmd.Flags().StringVar(&pipeSelector, "selector", "", "event selector")
	addPipeCmd.Flags().BoolVar(&pipeMatchAll, "match-all", false, "match all events")
	_ = addPipeCmd.MarkFlagRequired("router")
	_ = addPipeCmd.MarkFlagRequired("name")
	_ = addPipeCmd.MarkFlagRequired("processor")

	delPipeCmd.Flags().StringVar(&routerName, "router", "", "Router name or ID")
	delPipeCmd.Flags().StringVar(&pipeName, "pipe", "", "Pipe name or ID")
	_ = delPipeCmd.MarkFlagRequired("router")
	_ = delPipeCmd.MarkFlagRequired("pipe")

	delRouterCmd.Flags().StringVar(&routerName, "router", "", "Router name or ID")
	_ = delRouterCmd.MarkFlagRequired("router")

	//addSinkCmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID")

	_ = addSinkCmd.MarkFlagRequired("sink-type")
	_ = addSinkCmd.MarkFlagRequired("name")

	delSourceCmd.Flags().StringVar(&resourceID, "id", "", "data source ID")
	_ = delSourceCmd.MarkFlagRequired("id")

	delSinkCmd.Flags().StringVar(&resourceID, "id", "", "data sink ID")
	_ = delSinkCmd.MarkFlagRequired("id")

	addRouterCmd.Flags().StringVar(&processorName, "processor", "", "processor name")
	addRouterCmd.Flags().StringVar(&routerName, "router-name", "", "Router Name")

	_ = addRouterCmd.MarkFlagRequired("processor")

	connectRouterCmd.Flags().StringVar(&sourceID, "source-id", "", "source ID")
	connectRouterCmd.Flags().StringVar(&routerID, "router-id", "", "Router ID")

	_ = connectRouterCmd.MarkFlagRequired("source-id")
	_ = connectRouterCmd.MarkFlagRequired("router-id")

	connectSinkCmd.Flags().StringVar(&sinkID, "sink-id", "", "sink ID")
	connectSinkCmd.Flags().StringVar(&routerID, "router-id", "", "Router ID")

	_ = connectSinkCmd.MarkFlagRequired("sink-id")
	_ = connectSinkCmd.MarkFlagRequired("router-id")

	updatePipeProcessorCmd.Flags().StringVar(&routerName, "router", "", "Router name")
	updatePipeProcessorCmd.Flags().StringVar(&pipeName, "pipe", "", "Pipe name")
	updatePipeProcessorCmd.Flags().StringVar(&processorName, "processor", "", "Processor name")

	_ = updatePipeProcessorCmd.MarkFlagRequired("router")
	_ = updatePipeProcessorCmd.MarkFlagRequired("pipe")
	_ = updatePipeProcessorCmd.MarkFlagRequired("processor")

}

// --- pipes and routers -------------------------------------------------------
//
// platform_router_add_pipe / platform_router_delete_pipe / platform_router_dao
// have always existed on the API; these are the missing CLI bindings. Without
// them the only way to attach a second pipe (a behavior pipe, say) to an
// existing application router was to POST the api/ds call by hand.

var (
	pipeProcessor  string
	pipeSinks      []string
	pipeTags       []string
	pipePriority   int
	pipeSelector   string
	pipeMatchAll   bool
	redisQueue     string
)

// parseTags turns --tag application=Varonis into the ownership tags the app
// installer writes.
func parseTags(pairs []string) ([]*model.Tag, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	tags := make([]*model.Tag, 0, len(pairs))
	for _, p := range pairs {
		name, value, found := strings.Cut(p, "=")
		if !found || name == "" {
			return nil, fmt.Errorf("invalid --tag %q: expected name=value", p)
		}
		tags = append(tags, &model.Tag{Name: name, Value: value})
	}
	return tags, nil
}

// resolveRouter accepts a router id or a router name and returns the router.
func resolveRouter(ref string) (*model.RouterConfig, error) {
	if ref == "" {
		return nil, fmt.Errorf("--router is required (name or id)")
	}
	configs, err := AppAPI.ListConfigs()
	if err != nil {
		return nil, err
	}
	for _, r := range configs.Routers {
		if r != nil && (r.ID == ref || r.Name == ref) {
			return r, nil
		}
	}
	names := make([]string, 0, len(configs.Routers))
	for _, r := range configs.Routers {
		if r != nil {
			names = append(names, r.Name)
		}
	}
	sort.Strings(names)
	return nil, fmt.Errorf("no router named or numbered %q; available routers: %s", ref, strings.Join(names, ", "))
}

// resolveSinkIDs accepts sink ids or sink names.
func resolveSinkIDs(refs []string) ([]string, error) {
	if len(refs) == 0 {
		return nil, nil
	}
	configs, err := AppAPI.ListConfigs()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(refs))
	for _, ref := range refs {
		found := ""
		for _, s := range configs.Sinks {
			if s != nil && (s.ID == ref || s.Name == ref) {
				found = s.ID
				break
			}
		}
		if found == "" {
			return nil, fmt.Errorf("no sink named or numbered %q", ref)
		}
		ids = append(ids, found)
	}
	return ids, nil
}

var listRouterCmd = &cobra.Command{
	Use:   "list-router",
	Short: "List routers with their pipes",
	Long: `List every router of the account with the pipes hanging off it.

platform_router_dao has no "list" action, so the whole topology is read with
platform_list_configs and printed here.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configs, err := AppAPI.ListConfigs()
		if err != nil {
			return err
		}
		if len(configs.Routers) == 0 {
			cmd.Println("No router found.")
			return nil
		}
		sinkNames := map[string]string{}
		for _, s := range configs.Sinks {
			if s != nil {
				sinkNames[s.ID] = s.Name
			}
		}
		pipesOf := map[string][]*model.StreamPipeConfig{}
		for _, p := range configs.Pipes {
			if p != nil {
				pipesOf[p.RouterID] = append(pipesOf[p.RouterID], p)
			}
		}
		for _, r := range configs.Routers {
			if r == nil {
				continue
			}
			cmd.PrintErrf("Router: %s, ID: %s, Workers: %d\n", r.Name, r.ID, r.WorkerCount)
			for _, p := range pipesOf[r.ID] {
				sinks := make([]string, 0, len(p.SinkIDs))
				for _, id := range p.SinkIDs {
					if n, ok := sinkNames[id]; ok {
						sinks = append(sinks, n+" ("+id+")")
					} else {
						sinks = append(sinks, id)
					}
				}
				cmd.PrintErrf("  Pipe: %s, ID: %s, Priority: %d, Processors: %s, Sinks: %s\n",
					p.Name, p.ID, p.Priority, strings.Join(p.ProcessorNames, ","), strings.Join(sinks, ", "))
			}
		}
		return nil
	},
}

var addPipeCmd = &cobra.Command{
	Use:   "add-pipe",
	Short: "Add a pipe to an existing router",
	Long: `Attach a new pipe to a router that already exists, leaving its other pipes alone.

This is how a second consumer -- a behavior pipe beside an application's main
pipe -- is added without reinstalling the application template:

  ingext stream add-pipe --router Varonis-default --name Varonis-default-Behavior \
      --processor Varonis_Behavior --sink Behavior-Varonis-default --priority 2000 \
      --tag application=Varonis --tag appInstance=default

--router and --sink take a name or an id. A pipe carries exactly one processor:
the API models processorNames as a list, but only the one is supported, so a
chain is built as a second pipe that the first hands on to by returning
"abort". Priority orders pipes within the router: the templates give a main pipe
1000 and a behavior pipe 2000, so a pipe added without --priority (0) runs ahead
of both.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resourceName == "" {
			return fmt.Errorf("--name is required")
		}
		if pipeProcessor == "" {
			return fmt.Errorf("--processor is required")
		}
		router, err := resolveRouter(routerName)
		if err != nil {
			return err
		}
		sinkIDs, err := resolveSinkIDs(pipeSinks)
		if err != nil {
			return err
		}
		tags, err := parseTags(pipeTags)
		if err != nil {
			return err
		}
		pipe := &model.StreamPipeConfig{
			Name:           resourceName,
			RouterID:       router.ID,
			MatchAll:       pipeMatchAll,
			Selector:       pipeSelector,
			// The API models this as a list, but a pipe carries exactly one
			// processor. Chain work by having the processor do it, or by a
			// second pipe that the first hands on to with "abort".
			ProcessorNames: []string{pipeProcessor},
			SinkIDs:        sinkIDs,
			Priority:       pipePriority,
			Tags:           tags,
		}
		id, err := AppAPI.AddRouterPipe(router.ID, pipe)
		if err != nil {
			return err
		}
		cmd.PrintErrf("Pipe added successfully: %s on router %s (%s)\n", id, router.Name, router.ID)
		fmt.Println(id)
		return nil
	},
}

var delPipeCmd = &cobra.Command{
	Use:   "del-pipe",
	Short: "Delete a pipe from a router",
	Long: `Remove one pipe from a router by pipe name or id. The router's other pipes keep
running, so this is the rollback for add-pipe.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if pipeName == "" {
			return fmt.Errorf("--pipe is required (name or id)")
		}
		router, err := resolveRouter(routerName)
		if err != nil {
			return err
		}
		configs, err := AppAPI.ListConfigs()
		if err != nil {
			return err
		}
		pipeID := ""
		for _, p := range configs.Pipes {
			if p != nil && p.RouterID == router.ID && (p.ID == pipeName || p.Name == pipeName) {
				pipeID = p.ID
				break
			}
		}
		if pipeID == "" {
			return fmt.Errorf("no pipe named or numbered %q on router %s", pipeName, router.Name)
		}
		if err := AppAPI.DeleteRouterPipe(router.ID, pipeID); err != nil {
			return err
		}
		cmd.PrintErrf("Pipe deleted successfully: %s from router %s (%s)\n", pipeID, router.Name, router.ID)
		return nil
	},
}

var delRouterCmd = &cobra.Command{
	Use:   "del-router",
	Short: "Delete a router",
	RunE: func(cmd *cobra.Command, args []string) error {
		router, err := resolveRouter(routerName)
		if err != nil {
			return err
		}
		if err := AppAPI.DeleteRouter(router.ID); err != nil {
			return err
		}
		cmd.PrintErrf("Router deleted successfully: %s (%s)\n", router.Name, router.ID)
		return nil
	},
}
