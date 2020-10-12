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

Aggregates must be transactional consistent, what means this? For each use case we will just mutate one aggregate. This will help us to not attach the service to transactions.

e.g.
In create space command we will get the warehouse to attach to space we will build the space attaching the warehouse and then we will just Save the Space aggregate.

## Events

As we have to be transactional consistent we cannot store the events in the application, the events will be "published" in the same transaction as we save the aggregate. This means the repository will be the responsible of "publishing" (saving in the db, to later publish them).


## Example of service

```go

type UsersRepository interface {
    UsersGetter
    Update(context.Context, domain.User) error
}

type UsersGetter interface {
    ByID(context.Context, domain.UserID) (domain.User, error)
}

// some where used by Company service
type CompaniesRepository interface{
    CompaniesGetter
    Add(context.Context, domain.Company) error
}

type CompaniesGetter interface {
    ByID(context.Context, domain.CompanyID) (domain.Company, error)
}

type UserService struct {
    users UsersRepository
    companies CompaniesGetter
}

type UserChangeCompanyRequest {
    UserID string
    CompanyID string
}

func (s UserService) ChangeCompany(ctx context.Context, r UserChangeCompanyRequest) error {
    userID, err := domain.NewUserID(r.UserID)
    if err != nil {
        return errors.WrapWrongInput(err, "changeCompany")
    }
    companyID, err := domain.NewCompanyID(r.UserID)
     if err != nil {
        return errors.WrapWrongInput(err, "changeCompany")
    }

    // get the user
    u, err := s.users.ByID(ctx, userID)
    if err != nil {
        return errors.WrapNotFound(err, "changeCompany")
    }

    // get the company
    c, err := s.companies.ByID(ctx, companyID)
    if err != nil {
        return errors.WrapNotFound(err, "changeCompany")
    }

    // we send the company to update, then the aggregate will do its own validations (like check if user already is in the given company) and generate needed events
    if err := u.ChangeCompany(c); err != nil {
        return errors.WrapWrongInput(err, "changeCompany")
    }

    // Stores the updated user and also the generated events
    if err := s.users.Update(u); err != nil {
        return errors.WrapInternalError(err, "changeCompany")
    }

    return nil
}

```
