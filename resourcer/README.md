#Resourcer Middleware

This middleware automatically loads Resources from your ORM (**Currently only supports GORM**)
into your handler considering parameters that are specified from you from your URLPATH.

This is inspired by rails' gem hosted [here](https://github.com/josevalim/inherited_resources), so if you liked that you will probably like this as well,
although it is a little bit verbose.

## **IMPORTANT**
This middleware currently needs a little modification to gin's core, so I will put my fork of the framework in the Godeps.
To check the status of my fork see [this](https://github.com/sirlori/gin).

## Examples

```go
  type User struct {
    Id       int
    Name     string
    Password string
  }
 
  type GormResourcer struct {
    resourcer.ResourceGorm
  }
  func (_ *GormResourcer) ResourceClass() (reflect.Type, string) {
                                  //Name of the param withou _id
    return reflect.TypeOf(User{}), "user"
  }

  func main(){
    // router init
                        // Use your favourite db.
                        // So you really want to use Postgres :)
    db, err := gorm.Open("postgres", "database config here")
    group := router.Group("/users")
 
    group.Use(ResourcerGorm(&GormResourcer{}, &db))
    group.GET("/:user_id", func(c *gin.Context) {
      // c.Resource to have the results of the query
      // You should also do some type assertion, in this case:
      // res := c.Resource.(*User)
    })
   
    // Now requests like /users/1 will get the user with id==1 automatically in your c.Resource
}
```

Use the resourcer.COllectionerGorm for collection instead of single resources, in pages where you have to list
the resources. In this case would be in the imaginary path: "/users"

***Enjoy! :D***
