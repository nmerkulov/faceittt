Instructions:

 To run everythin in docker-compose use `docker-compose up -d`
 Maybe you will need to run this twice, since i added no wait-scripts (i run everything locally)
 Alternavely you can build cod and run
 To run app first run postgres using `docker-compose  up -d postgres`, then `go run main.go`
 
  
 
 Migrations are built in, but usually i use something like migrate/migrate to have migrations versioned and executed as 12FA recommends
 
Notes:
1) requests.http includes typical requets can be executed against locally running app
2) I had no time to write actual tests but code itself is test-ready. I used dependency injection pattern (please note without any
mysterious library, quite straitforward from main) so each layer/handler/usecase can be tested completely isolated
3) Even if it is not much lines of code i decided to provide "hex/onion/clean"-like project structure for demostration purposes
4) Please pay attention to "eventHandler" funcs at application layer. This is  sort-of hooks for handling events. Lets say - it can be very specific
logging or it can be kafka/rabbit/nats emitting. So by that mechanics we can plug in new notificators  as we want.
5) API build sort-of REST api. Few things to mention:
    * Some properties are immytable - password, email nickname and ID. Thats just my assumption for application logic
    * Editing API stands on a PUT method, which idiomatically often used as "update evertything according to request" 
    (or "create or update" in some cases. not ours). To explain it i'll provide opposite example - JSON-PATCH or RFC6902. So, in 
    our API, if some field is not presented in payload - this field will be flushed in the database as well. Usually i prefer to use
    PAT~CH with partial update, but it is a bit harder to implement
    * There is no Pagination, but, of course, it must be there. I can explain how to do it on interview ;)
    
6) Please note, i provided 2 data access interfaces - UserRepo and UserFinder. Usually all methods (including Find-like) can be included
into repo interface, but for more-or-less rich query parameters it will require to build sort-of DSL with predicates.
Thats...can be quite interesting, but it might become error-prone and not-so-easy to track abstraction and implementation sync.
As alternative i suggest you to use UserFinder. You are welcome to chain calls like "finder.WithName("name"").WithCountry("ru"").Find()"
Main benefit here is quite straightforward contract and..your code won't compile if you didn't implement "With"-method. Main cons here
is that interface can become overwhelming, with a LOT of methods. So anyway - then one should always keep track of balance in code. 
It never becomes perfect, but it must be good enough
7) Also nice improvement point is to have "request-id" in context
8) Nice improvement point - metrics. Ther are no metrics in application so far, but it is one of key things in order to use
application in production
9) I prefer to generate boilerplate code like handlers by using things like goswagger, but in this case it is only few endpoints so 
i decided to use chi