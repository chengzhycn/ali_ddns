APP:=ali_ddns

$(APP):
	GOOS=linux go build

docker: $(APP)
	docker build -t $(APP) .

.PHONY: clean
clean:
	rm -rf $(APP)