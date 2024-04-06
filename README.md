# SimActive

One of kostudio helper mircoservice.
gRPC server that implements **Domain-Driven Design** and **Repository** pattern

#### Why am I using additional in-memory storage? 
Because it will minimize to ask SQL server(probably sql microservice in the future) with queries.
And it will limit sections of code which will be locked with `sync` tools.

#### Why is `Generic` have appeared in the code?
Cuz to implement `repository pattern` which is very usefull in `Domain-Driven Design` will generate so much duplicate code.
I prefer to use `type assertion` (without using package `reflect`, cuz go.dev) to CRUD operations if needed, cuz 
otherwise here should be place 2 implementations of repository(may 1, but nvm) for every domain model, for now its over 8, so on..
todo: possibly to remove it in future and use only interfaces where possible