# golang-programming-task

Welcome to the test, purpose of this exercise is to see your coding style and overall approach to software engineering and planning. In order to start the test please fork this repo, choose tasks to complete and once you are done send pull request.

### Project overview

The project is a simple http api server that provides a search endpoint for a document index.
Document index is implemented as in memory data structure. 
Index is built from a dataset of documents, where each document has a text description and link to an image. 
The search endpoint accepts any number of words in english as input in the query string parameter q and returns a list of all of the documents containing ALL words in the query.

### Project setup

1. Download the [Midjourney v5 dataset](https://huggingface.co/datasets/tarungupta83/MidJourney_v5_Prompt_dataset/resolve/main/Midjourney_v5_Prompt.zip) from Huggingface
2. Extract archive contents into `/tmp/data` directory. (Or any other directory and change the `IndexDataDir` value in config file later)
3. Clone the repository and enter to the top level directory
4. Run the api server: `make run_server`

If everything is set up correctly the output should look like this:
```
# make run_server
go build -o out/ github.com/getnexar/golang-programming-task/doc-index/cmd/... && out/doc-index-search-api
[Fx] PROVIDE	*config.Config <= main.provideConfig()
[Fx] PROVIDE	*zap.SugaredLogger <= main.provideLogger()
[Fx] PROVIDE	*index.Index <= main.provideIndex()
[Fx] PROVIDE	handlers.HandlersInterface <= main.provideHandlers()
[Fx] PROVIDE	http.Handler <= main.provideHttpRouter()
[Fx] PROVIDE	*server.HttpServer <= main.provideHttpServer()
[Fx] PROVIDE	fx.Lifecycle <= go.uber.org/fx.New.func1()
[Fx] PROVIDE	fx.Shutdowner <= go.uber.org/fx.(*App).shutdowner-fm()
[Fx] PROVIDE	fx.DotGraph <= go.uber.org/fx.(*App).dotGraph-fm()
[Fx] INVOKE		main.registerLifecycleHooks()
[Fx] RUN	provide: go.uber.org/fx.New.func1()
[Fx] RUN	provide: main.provideConfig()
[Fx] RUN	provide: main.provideLogger()
{"level":"info","ts":"2023-11-10T23:56:49.175-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_7.csv"}
{"level":"info","ts":"2023-11-10T23:56:49.300-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_6.csv"}
{"level":"info","ts":"2023-11-10T23:56:49.428-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_4.csv"}
{"level":"info","ts":"2023-11-10T23:56:49.535-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_18.csv"}
{"level":"info","ts":"2023-11-10T23:56:49.624-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_19.csv"}
{"level":"info","ts":"2023-11-10T23:56:49.660-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_5.csv"}
{"level":"info","ts":"2023-11-10T23:56:49.792-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_1.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.125-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_2.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.267-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_3.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.389-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_12.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.494-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_13.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.634-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_11.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.715-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_10.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.823-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_8.csv"}
{"level":"info","ts":"2023-11-10T23:56:50.941-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_14.csv"}
{"level":"info","ts":"2023-11-10T23:56:51.118-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_15.csv"}
{"level":"info","ts":"2023-11-10T23:56:51.305-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_9.csv"}
{"level":"info","ts":"2023-11-10T23:56:51.438-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_17.csv"}
{"level":"info","ts":"2023-11-10T23:56:51.481-0700","caller":"index/index.go:54","msg":"Loading data file","filepath":"/tmp/data/MJ_v5.1_Part_16.csv"}
{"level":"info","ts":"2023-11-10T23:56:51.812-0700","caller":"index/index.go:80","msg":"Loading index complete...","duration":2.637016}
[Fx] RUN	provide: main.provideIndex()
[Fx] RUN	provide: main.provideHandlers()
[Fx] RUN	provide: main.provideHttpRouter()
[Fx] RUN	provide: main.provideHttpServer()
[Fx] HOOK OnStart		main.registerLifecycleHooks.func1() executing (caller: main.registerLifecycleHooks)
[Fx] HOOK OnStart		main.registerLifecycleHooks.func1() called by main.registerLifecycleHooks ran successfully in 3.001364125s
[Fx] RUNNING
```

### Tasks

1. Change the search response format so it looks like this:

    ```
    {
    "results": [
        {
        "description": "**hello world** - Upscaled by @NeonNebula (fast)",
        "imageUrl": "https://cdn.discordapp.com/attachments/941971306004504638/1099717251684376596/NeonNebula_hello_world_499f2603-2ee9-4d1d-bdb2-6fce9af4dd55.png"
        },
        {
        "description": "**hello world blue technology toy style --v 5** - @TechnoKing (fast)",
        "imageUrl": "https://cdn.discordapp.com/attachments/995431233121161246/1101927504731701401/TechnoKing_hello_world_blue_technology_toy_style_f37ac2e2-94e4-4c8c-b0b0-1574841db1a3.png"
        }
    ]
    }
    ```

2. Optimize the `Index` structure and search api endpoint so it works much faster and outputs only full keyword matches
   
3. Add tests for `Index`

4. [Optional] Add a new endpoint which deletes documents from the Index for a given keyword / set of keywords

* You should not replace `Index` implementation with any database solution

### Benchmarking.

You can test the Index implementation performance by running perf tool like `hey`

Initial implementation:
```
# hey -n 30 -c 1 'http://localhost:8080/search?q=hello&q=world&q=fast'

Summary:
  Total:	12.0965 secs
  Slowest:	0.4220 secs
  Fastest:	0.3952 secs
  Average:	0.4032 secs
  Requests/sec:	2.4801


Response time histogram:
  0.395 [1]	|■■■■■
  0.398 [4]	|■■■■■■■■■■■■■■■■■■■■
  0.401 [3]	|■■■■■■■■■■■■■■■
  0.403 [8]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.406 [8]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.409 [4]	|■■■■■■■■■■■■■■■■■■■■
  0.411 [1]	|■■■■■
  0.414 [0]	|
  0.417 [0]	|
  0.419 [0]	|
  0.422 [1]	|■■■■■


Latency distribution:
  10% in 0.3976 secs
  25% in 0.4009 secs
  50% in 0.4031 secs
  75% in 0.4054 secs
  90% in 0.4079 secs
  95% in 0.4220 secs
  0% in 0.0000 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0002 secs, 0.3952 secs, 0.4220 secs
  DNS-lookup:	0.0001 secs, 0.0000 secs, 0.0021 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0008 secs
  resp wait:	0.4028 secs, 0.3951 secs, 0.4135 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0004 secs

Status code distribution:
  [200]	30 responses
```

Target implementation performance:
```
# hey -n 30 -c 1 'http://localhost:8080/search?q=hello&q=world&q=test'

Summary:
  Total:	0.0964 secs
  Slowest:	0.0186 secs
  Fastest:	0.0024 secs
  Average:	0.0032 secs
  Requests/sec:	311.0577

  Total data:	450 bytes
  Size/request:	15 bytes

Response time histogram:
  0.002 [1]	|■
  0.004 [28]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.006 [0]	|
  0.007 [0]	|
  0.009 [0]	|
  0.010 [0]	|
  0.012 [0]	|
  0.014 [0]	|
  0.015 [0]	|
  0.017 [0]	|
  0.019 [1]	|■


Latency distribution:
  10% in 0.0024 secs
  25% in 0.0025 secs
  50% in 0.0026 secs
  75% in 0.0028 secs
  90% in 0.0033 secs
  95% in 0.0186 secs
  0% in 0.0000 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0003 secs, 0.0024 secs, 0.0186 secs
  DNS-lookup:	0.0001 secs, 0.0000 secs, 0.0027 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0007 secs
  resp wait:	0.0028 secs, 0.0023 secs, 0.0071 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0007 secs

Status code distribution:
  [200]	30 responses
```
