default: run

install:
	go install

build: install
	go build

run: build
	./USDT

clean:
	rm *.txt USDT *.log *.csv