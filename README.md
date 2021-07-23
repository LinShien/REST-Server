# REST-Server
Implementation of REST server. [Tutorial here](https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/)

* [Use Just Standard Library](#StandardLib)

## <a name="StandardLib"> Use Just Standard Library </a>
    POST   /task/              :  creates a new task, and then returns ID
    GET    /task/<taskid>      :  returns a single task by <taskid> 
    GET    /task/              :  returns all tasks
    DELETE /task/<taskid>      :  deletes a task by <taskid>
    GET    /tag/<tagname>      :  returns list of tasks with <tagname> tag
    GET    /due/<yy>/<mm>/<dd> :  returns list of tasks due by date <yy>/<mm>/<dd>

## How to test the REST server
* Write golang programs and import standard lib "testing".

    `> go test name_of_the_testing.go`
    
* Public testing API like Advanced Rest Client Application.

    <img src="https://i.imgur.com/UzxI6P9.png">

---

## Useful golang references 
* [Synchronizing Structs for Safe Concurrency in Go](https://bbengfort.github.io/2017/02/synchronizing-structs/)
* [make vs new in Golang](https://medium.com/d-d-mag/golang-%E7%AD%86%E8%A8%98-make-%E8%88%87-new-%E7%9A%84%E5%B7%AE%E5%88%A5-68b05c7ce016)
* [HTTP service provided by Golang std lib](https://www.jianshu.com/p/16210100d43d)
* [How to Parse a JSON Request Body in Go](https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)
* [Mutex in Golang](https://tour.golang.org/concurrency/9)
* [gorilla/mux](https://github.com/gorilla/mux)
