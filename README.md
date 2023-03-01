# Serverless Functions Go 💜

[![build-and-test](https://github.com/scaleway/serverless-functions-go/actions/workflows/test.yml/badge.svg)](https://github.com/scaleway/serverless-functions-go/actions/workflows/test.yml)
[![golangci-lint](https://github.com/scaleway/serverless-functions-go/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/scaleway/serverless-functions-go/actions/workflows/golangci-lint.yml)

Scaleway Serverless Functions is a framework to provide a good developer experience to write Serverless Functions.

[Serverless Functions](https://www.scaleway.com/fr/serverless-functions/) make it easy to deploy, scale, and optimize your workloads on the cloud.

Get started with Scaleway Functions (we support multiple languages :rocket:):

- [Scaleway Serverless Functions Documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)
- [Scaleway Serverless Framework plugin](https://github.com/scaleway/serverless-scaleway-functions)
- [Scaleway Serverless Examples](https://github.com/scaleway/serverless-examples)
- [Scaleway Cloud Provider](https://scaleway.com)

Testing frameworks for Scaleway Serverless Functions in other languages can be found here:

- [Node](https://github.com/scaleway/serverless-functions-node)
- [PHP](https://github.com/scaleway/serverless-functions-php)
- [Python](https://github.com/scaleway/serverless-functions-python)
- [Rust](https://github.com/scaleway/serverless-functions-rust)

## ⚙️ Quickstart

```sh
 go get github.com/scaleway/serverless-functions-go
```

`cmd/main.go`:

```go
import "github.com/scaleway/serverless-functions-go/functest"

func main() {
	// Replace "Handle" with your function handler name if necessary
	functest.ServeHandlerLocally(Handle, functest.WithPort(8080))
}
```

For advanced usage please check [usage section](#-advanced-usage).

## 🚀 Features

This repository aims to provide the best experience: **local testing, utils, documentation etc...**
additionally we love to share things with the community and we want to expose receipts to the public. That's why
we make our framework publicly available to help the community!

## 🏡 Local testing

What this package does:

- **Format Input**: FaaS have a specific input format encapsulating the body received by functions to add some useful data.
  The local testing package lets you interact with this data.
- **Advanced debugging**: To improve developer experience you can run your handler locally, on your computer to make
  it simpler to debug by running your code step-by-step or reading output directly before deploying it.

What this package does not:

- **Simulate performance**: Scaleway FaaS lets you choose different options for CPU/RAM that can have an impact
  on your development. This package does not provide specific limits for your function on local testing but you can
  add [Profile your application](https://go.dev/blog/pprof) or you can use our metrics available in [Scaleway Console](https://console.scaleway.com/)
  to monitor your application.
- **Build functions**: When your function is uploaded we build it in an environment that can be different than yours. Our build pipelines support
  tons of different packages but sometimes it requires a specific setup, for example, if your function requires a specific 3D system library.
  If you have compatibility issues, please see the help section.

## 🛟 Help & support

- Scaleway support is available on Scaleway Console.
- Additionally, you can join our [Slack Community](https://www.scaleway.com/en/docs/tutorials/scaleway-slack-community/)

## 🎓 Contributing

We welcome all contributions to our open-source projects, please see our contributing guidelines <link>.

Do not hesitate to raise issues and pull requests we will have a look at them.

## 🔬 Advanced usage

To run the function locally you need to add an entry point to serve your function.

This entrypoint should be put in a directory not required by your handler, e.g. `cmd/main.go`.

In your `run/main.go` add the following code to invoke your function :

```go
package main

import (
  // "localfunc" is the module name located in your go.mod. To generate a go.mod with localfunc as name you
  // can use the following command : go mod init localfunc
  // Or you can replace "localfunc" with your own module name.
	localfunc "github.com/scaleway/serverless-functions-go/examples/handler"
	"github.com/scaleway/serverless-functions-go/functest"
)

func main() {
	// Replace "Handle" with your function handler name if necessary
	functest.ServeHandlerLocally(localfunc.Handle, functest.WithPort(8080))
}

```

This file will expose your handler on a local web server allowing you to test your function.

Some information will be added to requests for example specific headers. For local development, additional header values are hardcoded
to make it easy to differentiate them. In production, you will be able to observe headers with exploitable data.

Local testing part of this framework does not aim to simulate 100% production but it aims to make it easier to work with functions locally.

### Cli

To run the server locally: `go run cmd/main.go`

### VS Code

Open `cmd/main.go` and open the "Run and Debug" pannel to execute or debug your function there is no special
configuration to add to VSCode.

### Goland

The IDE will generate a run configuration for you, open `cmd/main.go` and run the main.

## ❓ FAQ

**Why do I need an additional package to call my function?**

Your Function Handler can be served by a simple HTTP server but Serverless Ecosystem involves a lot of different layers
and this package aims to simulate everything your request will go through.

**How my function will be deployed**

Your function will be deployed in an environment that allows your function to easily Scale up and down and it's wrapped into
different pieces of software with different roles. This stack also changes the headers, input and output of your function, that's why
this tool has been developed to simulate this part.

**Do I need to deploy my function differently?**

No. This framework does not affect deployment nor performance.

## 🏛️ Architecture

To make development and understanding of this repository we tried to keep the path of the request natural.

- [framework](./framework/) folder is used to store all the code that you can import into your project
- [functest](./functest) contains all the cool tools to work locally with your function 😎

## 🧑‍💻 Development

This repository is at its early stage and is still in active development.
If you are looking for a way to contribute please read [CONTRIBUTING.md](./docs/CONTRIBUTING.md).

## 📭 Reach Us

We love feedback.
Don't hesitate to open a [Github issue](https://github.com/scaleway/serverless-functions/go/issues/new) or
feel free to reach us on [Scaleway Slack community](https://slack.scaleway.com/),
we are waiting for you on [#serverless-functions](https://scaleway-community.slack.com/app_redirect?channel=serverless-functions).
