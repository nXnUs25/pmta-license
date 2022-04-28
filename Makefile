.PHONY : clean go

go:
	go fmt
	go build -o check-pmta-license *.go
	ls -l check-pmta-license

clean:
	rm ./pmta-license

