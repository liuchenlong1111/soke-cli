BINARY := soke-cli
VERSION := dev
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: build install test e2e-test test-all clean \
	test-contact test-course test-exam test-certificate test-credit \
	test-point test-training test-learning-map test-news test-clock test-file

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install: build
	install -m755 $(BINARY) /usr/local/bin/$(BINARY)

test:
	go test -v ./...

e2e-test: build
	./scripts/e2e-test.sh

test-all: test e2e-test

# 分模块测试
test-contact: build
	@echo "测试通讯录模块..."
	@./scripts/e2e-test.sh contact

test-course: build
	@echo "测试课程模块..."
	@./scripts/e2e-test.sh course

test-exam: build
	@echo "测试考试模块..."
	@./scripts/e2e-test.sh exam

test-certificate: build
	@echo "测试证书模块..."
	@./scripts/e2e-test.sh certificate

test-credit: build
	@echo "测试学分模块..."
	@./scripts/e2e-test.sh credit

test-point: build
	@echo "测试积分模块..."
	@./scripts/e2e-test.sh point

test-training: build
	@echo "测试线下培训模块..."
	@./scripts/e2e-test.sh training

test-learning-map: build
	@echo "测试学习地图模块..."
	@./scripts/e2e-test.sh learning-map

test-news: build
	@echo "测试新闻模块..."
	@./scripts/e2e-test.sh news

test-clock: build
	@echo "测试作业模块..."
	@./scripts/e2e-test.sh clock

test-file: build
	@echo "测试素材库模块..."
	@./scripts/e2e-test.sh file

clean:
	rm -f $(BINARY)
