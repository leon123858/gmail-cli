all:
	go run ./main.go

deploy:
	git ls-remote --tags origin
	git tag v0.1.0
	git push origin v0.1.0
	git tag lastest -f
	git push origin lastest -f