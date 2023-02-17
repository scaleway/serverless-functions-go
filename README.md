# Serverless Functions Go

Scaleway Serverless Functions is a framework to provide a good developer experience to write Serverless Functions.

Serverless Funcions makes it easy to deploy, scale, optimise your workloads on the cloud.

Get started with Scaleway Functions (we support multiple langages :rocket:) :

- [Scaleway Serverless Functions Documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)
- [Serverless Scaleway plugin](https://github.com/scaleway/serverless-scaleway-functions)
- [Serverless Examples](https://github.com/scaleway/serverless-examples)
- [Scaleway Cloud Provider](https://scaleway.com)

If you are looking for framework about other runtimes refer to the links below :

- [Node](https://github.com/scaleway/serverless-functions-node)
- [Python](https://github.com/scaleway/serverless-functions-python)
- [Rust](https://github.com/scaleway/serverless-functions-rust)
- [PHP](https://github.com/scaleway/serverless-functions-php)

## 🚀 Features

This repository aims to provide the best experience : **local testing, utils, documentations etc...**
additionnaly we love to share things with community and we want to expose receipts to public.

## 🏡 Local testing

What this packages does :

- **Format Input**: FaaS have specific input format encapsulate the body recieved by functions to add some useful data.
  The local testing package let you interact with this data.
- **Advanced debugging**: To improve developer experience you can run your handler locally, on your computer to make
  it simpler to debug by running your code step-by-step or reading output directly before deploying it.

What this packages does not :

- **Simulate performance**: Scaleway FaaS let you choose different options for CPU/RAM that can have impact
  on your developments. This package does not provide specific limits for your function on local testing but you can
  add [Profile your application](https://go.dev/blog/pprof) or you can use our metrics available in [Scaleway Console](https://console.scaleway.com/)
  to monitor your application.
- **Build functions**: When your function is uploaded we build it in an environment that can be different than yours. Our build pipelines supports
  tons of different packages but sometimes it requires specific setup, as example if your function requires specific 3D system
  libraries from your GPU card provider. In case of deployment error please check help section

## 🛟 Help & support

- Scaleway support is available on Scaleway Console.
- Additionnaly you can join our [Slack Community](https://www.scaleway.com/en/docs/tutorials/scaleway-slack-community/)

## 🎓 Contributing

There are many ways to contribute to Scaleway Serverless FaaS, Serverless ecosystem evolution is amazing and you can take part of this by submitting new proposals, ideas, submit bugs and contribute to documentation.

Do not hesitate to raise issues and pull requests we will have a look at it.

# Usage

In order to run the function locally you need to add an entry point to serve your function.

So create a new file in your project in a folder that is not required by your handler, example : `cmd/main.go`.

In your `run/main.go` add the following code to invoke your function :

```go
package main

import (
  // "localfunc" is the module name located in your go.mod. To generate a go.mod with localfunc as name you
  // can use the following command : go mod init localfunc
  // Or you can replace "localfunc" with your own module name.
	func "localfunc"
)

func main() {
	// Replace "Handle" with your function handler name if necessary
	server.ScalewayRouter(Handle, server.WithPort(8080))
}

```

This file will expose your handler on a local web-server allowing you to test your function.

### Cli

To run the server locally, run `go run cmd/main.go`

### VS Code

Open `cmd/main.go` and open the "Run and Debug" pannel to execute or debug your function, there is no special
configuration to add

### Goland

The IDE will generate a run configuration for you, open `cmd/main.go` and run the main.
