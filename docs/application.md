# Application

Application is divided by services, depending on the aggregates.

Services are just the bridge between infrastructure and domain layers.
They are responsible of building aggregates and entities needed by the domain and store them.

As go is created for simplicity, use cases are grouped by service, not like other languages that we will need to create one struct (class) per use case.

As we don't want to attach services to the different infrastructure implementations we use interfaces to communicate with it.

## Requests

The requests that gets the command must be DTOs of non domain values, like primitive values or Application/Service values. Then the command will be responsible to build the domain id.

## Repositories

Repository interfaces are composed by the mutate functions like `Save` and `Update` that are used only by the aggregate service to store changes of the aggregate, and by `{entity}Getter` interface, that will define the functions to get entities by different criteria (`ByID`, `ByExternalID`, `All`...). This composition is done in purpose of allow other aggregate services to just use the getter interface.

e.g.
SpaceService will use SpaceRepository as it needs to save Space aggregate and will use WarehouseGetter as it needs to retrieve the warehouse on space create to attach the warehouse to the space. It doesn't make sense to use WarehouseRepository as SpaceService must not update a warehouse.

## Transactions

Service use case will be responsible of the transactional consistent, for that it will use unit of work, which everything that runs inside it will run in the same transaction.

## Events

In order to have a log and to be detached from other services, use cases, for each action in the system we will raise an event.
