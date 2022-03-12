# BlogApi in Golang

This is the simple blog api in which user have to create account and then login to create a blog post and can comment to anyone's blog post.

## API Methods

1. POST - 

    To create an account.
    
    To create a blog post.
    
    To create comment on post.
    
    To login.

2. GET - 

     To get the list of person.
     
     To get the list of blog post.
     
     To get the list of comments.


## Quick-setup in some steps
 
1. Install Golang on your machine. 

    1.1. Make sure you have GOPATH set in your environment variables. 

    1.2. Ensure it using `echo %GOPATH%`
 
2. Get this project by this command: `go get -u https://gitlab.com/htp22/mashiyatblog.git`

3. This will take some time because it downloads this project and downloads all the imported dependencies.

4. Now, `cd mashiyatblog`

5. Now, run a mogodb server on your local machine.

6. Run `go build` to build the go project in a executable file.


## Endpoints 
(Testing can be done using POSTMAN)

1. Create an account > POST `"http://localhost:8000/person"` > Enter three values in Body (raw) `firstname, lastname, email, password`.

2. Login > POST `"http://localhost:8000/login"` > Enter two values in Body (raw) `email password`.
 
3. Create a blog > POST `"http://localhost:8000/blogpage"` > Enter three key values in Body (raw) `title, description, user_id`.

4. Get all blogs > GET `"http://localhost:8000/blogpage"` > This get an array of all blogs.

5. Get a specific blog > GET `"http://localhost:8000/person/{_id}"` > This gets the object of the specified blog id.
