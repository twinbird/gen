PROGRAM = gen

all:	$(PROGRAM)

$(PROGRAM):
	CGO_ENABLED=0 go build

clean:
	rm -f $(PROGRAM)
