# Serverless Functions Go

Scaleway Serverless Functions is a framework to provide a good developer experience to write Serverless Functions.

Serverless Funcions makes it easy to deploy, scale, optimise your workloads on the cloud.

Get started with Scaleway Functions (we support multiple langages :rocket:) : 
- [Scaleway Serverless Functions Documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)
- [Serverless Scaleway plugin](https://github.com/scaleway/serverless-scaleway-functions)
- [Serverless Examples](https://github.com/scaleway/serverless-examples)
- [Scaleway Cloud Provider](https://scaleway.com)

If you are looking for framework about other runtimes refer to the links below :
* [Node]()
* [Python]()
* [Rust]()
* [PHP]()

## Features

This repository aims to provide the best experience : **local testing, utils, documentations etc...**
additionnaly we love to share things with community and we want to expose our tools/

### Local testing

What this packages does :
* **Format Input**: FaaS have specific input format encapsulate the body recieved by functions to add some useful data.
The local testing package let you interact with this data.
* **Advanced debugging**: To improve developer experience you can run your handler locally, on your computer to make
it simpler to debug by running your code step-by-step or reading output directly before deploying it.
* ****

What this packages does not : 
* **Simulate performance**: Scaleway FaaS let you choose different options for CPU/RAM that can have impact
on your developments. This package does not provide specific limits for your function on local testing but you can
add [Profile your application](https://go.dev/blog/pprof) or you can use our metrics available in [Scaleway Console](https://console.scaleway.com/)
to monitor your application.
* **Build functions**: When your function is uploaded we build it in an environment that can be different than yours. Our build pipelines supports
tons of different packages but sometimes it requires specific setup, as example if your function requires specific 3D system
libraries from your GPU card provider. In case of deployment error please check help section

## Roadmap

## Help & support

## Contributing