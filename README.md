# Expense Tracker Golang

This API is basically not done yet, i dont implement delete operation on this services, but other service is available.
The most thing that i confuse is about how to test it and get 100% coverage and i think that is imposible since we are not wrap the function into a interface

## Folder Structure

I'm not use overheard folder structure, all model/domain will placed on one place, repository?, interface?, entiry?, all of that are located on one folder
- internals: hold all internal function like the domain of this expense tracker (book, category, account, and record), configuration, and some middleware
- utils: all reusable function

## Configuration
Configuration go in 3 step:
- Default value
- Load from environment
- Load from config *.yml file

## Testing
Maybe we can use faking, create a separated function to fake Database behaviour using `tag` on golang build, we test on integration level not a unit level.

### I'm still learning, so if you have any question or find any mistakes by me please tell me
