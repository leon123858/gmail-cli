all:
	go run ./main.go

deploy:
	git ls-remote --tags origin
	git tag v0.0.5
	git push origin v0.0.5
	git tag lastest -f
	git push origin lastest -f