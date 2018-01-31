@title[Introduction]

# Gandalf

### An adventure in contract testing

Gandalf is a tool that leverages the go test framework to make defining contracts flexible and expressive.

---

@title[Problems]

## Problems

The last thing you want is to update your development libraries and have all your code break when you are not expecting it, with service to service communication this happens in production and possibly without warning.

Note: this is the main problem, additional problems described on slides below are optional.

+++

### Design

Designing an API and implementing it seperately can lead to changes creeping in and "back to the drawing board" moments when you finally see it is not solving what you need it too. Communication of the design can be ineffecient or leave room for misinterpretation.

+++

### Parrallel Development

A service that provides its contracts to client developers allows for additional confidence that the consumer will integrate with the provider without it existing. This is assured by the contracts being agreed to be met by the provider when it does get called for real.

---

@title[Solutions]

## Solutions

Contracts describe the interaction between a client (consumer) and a service (provider). Contracts should be slow to change if ever, this means that all versions of the service that meet the contract are able to be used interchangeably.

Note: this is the main solution, additonal solutions described below are optional.

+++

### Design

Formally describe your service before writing any code for your API, with the concept of Consumer Driven Contracts you may even work with your clients to decide how the service should behave.

Note: they are code with everything that implies, version control, code review, etc.

+++

### Parallel/Isolated Development

Gandalf can generate definitions for mocking services out of defined contracts for free, rapidly standing up a full mock implementation of a HTTP API service that can be used for testing or parallel development of eg. backend and frontend code.

Note: Spin up a mock of external services that do not provide testing sandboxes.

+++

### Benchmarking

Utilizing the go test benchmarking capabilities Gandalf can benchmark specific endpoints or groups of endpoints, even as complex as benchmarking a use case scenario for more real world workloads. This uses the same contracts used for test assertions so comes for free.
