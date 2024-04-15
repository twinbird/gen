PROGRAM = gen
SRC = main.go

gen: $(SRC)
	CGO_ENABLED=0 go build

clean:
	rm -f $(PROGRAM)
