# Domain

## Aggregates

### Race

It is an event that will happen in a 
The race that competitors will run.

### User

A person in the system

### Team

A group of Users

## Types

We have multiple types defined in the domain, not only entities or aggregates, also VOs. This has a reason, each type needs to take care of his own integrity.

e.g:
PackageName, on NewPackage name we validate that name is not empty. This allows us to have more meaning full validations on the domain, and also delegating the responsibilities of validations.

## References

- [what is an aggregate?](https://stackoverflow.com/a/1958722)
- [aggregates relations, entities or ids](https://enterprisecraftsmanship.com/posts/link-to-an-aggregate-reference-or-id)
- [primitive obsesion](https://enterprisecraftsmanship.com/posts/functional-c-primitive-obsession)
- [relationship](https://blog.sapiensworks.com/post/2016/08/24/DDD-Relationships)
- [bounded contexts](https://www.baeldung.com/java-modules-ddd-bounded-contexts)
