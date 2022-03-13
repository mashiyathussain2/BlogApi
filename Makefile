run: .env create bin/blog
	@PATH="$(PWD)/bin:$(PATH)" heroku local

bin/blog: main.go
	go build -o bin/blog main.go