# gcp-cloud-native-stack
An experimental repo building demo cloud native app using the Google cloud stack.

A few selected services from the GCP microservices demo app are used to build this app for learning purpose. Each component below is a microservice. All components are made to be standalone.

### Catalog


### Recommendation Engine


### Frontend


### Checkout


### Shipping


### Email


### Loadgen


# Dev Guide
Use the `Dockerfile-dev` container image (an alpine base image) for developing and testing the app. It has the `cgo enabled` environment variable so that we can run `go test`.
```
docker run -itd --rm -v /<path to>/gcp-cloud-native-stack:/go/src/github.com/tony-yang/gcp-cloud-native-stack go-dev-alpine
```

Within each component, run `make test` to start the unit test and `make build` to update the generated proto, fetch and update go modules, and build the packages.

Run `make cover` to run the unit test with coverage report breakdown by function.

# Design Principles
- This app can be run locally (k8s on Docker) or on GCP without any modification.
- This app will mimic the regular dev flow, init, test, run, and deploy locally for fast iteration, and can use the same pipeline to deploy to prod.
- Code may not be optimal (runtime-wise) but should be easy to read and short for demo and learning purpose.

# References
This project references the GCP Microservices demo to learn how to build a cloud-native app. Most of the code comes from that repo but is modified as needed.
- https://github.com/GoogleCloudPlatform/microservices-demo
