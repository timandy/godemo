# routine 编译器

## 基础知识

### 参数 `-toolexec` 作用

例:

> go build -toolexec='routiner go-agent' -a -o main.exe .

`go build` 是个高级的工具，其会执行 `compile.exe`、`asm.exe`和`link.exe` 等一系列过程。

参数 `-toolexec` 主要拦截编译过程，再调用 `compile.exe`，其参数为自带的 `库源码` 和 `用户源码`。

以编译 `runtime` 包为例，编译过程会经历以下的步骤：

```shell
# (1)
routiner go-agent compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
# (2)
go-agent compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
# (3)
compile.exe -o xx.a -trimpath $WORK\xx=> -p runtime ...... -asmhdr runtime.go runtime2.go
```

`go build` 会调用 `(1)`, 但是第 `(2)` 步需要 `(1)` 中的 `routiner` 自身逻辑发起 `(2)` 调用。

然后需要 `(2)` 中的 `go-agent` 发起 `(3)` 调用。

### `-a` 作用

`compile.exe` 会按包缓存编译产物，当编译过一次后会将结果缓存起来。

下一个编译过程如果用到了该包的编译产物，将跳过编译直接使用上次编译产物。

由于 `routiner` 会替换库文件的路径，为确保修改后的文件生效，所以要指定 `-a` 参数。

## 核心逻辑

`routiner` 的核心逻辑如下：

1. 过滤 `runtime` 包；
2. 使用 `ast` 解析 `runtime2.go` 等源文件结构，修改语法并将内容存到一个临时文件；
3. 修改命令行参数，使用新的文件路径替换命令行中 `runtime2.go` 的路径；
4. 使用修改后的参数，调用下一个工具链。

## 使用

### 安装

```shell
# 清理构建物
go clean
# 重新编译, 防止缓存
go build -a .
# 安装到 GOPATH 下
go install
```

### 单独使用

- windows

```shell
@echo off

# 将 %GOPATH%/bin 追加到 PATH, 以搜索到安装的 routiner.exe
SET PATH=%PATH%;%GOPATH%/bin
# 添加参数 -toolexec="routiner" -a
go build -toolexec="routiner" -a -o main.exe .
```

- linux

```shell
#!/bin/bash

# 将 %GOPATH%/bin 追加到 PATH, 以搜索到安装的 routiner
export PATH=$PATH:$GOPATH/bin
# 添加参数 -toolexec="routiner" -a
go build -toolexec="routiner" -a -o main.exe .
```

### 多个工具链组合使用

不保证其他工具链有链式传递功能，所以要将 `routiner` 放在其他工具链前边。

- windows

```shell
@echo off

# 将 %GOPATH%/bin 追加到 PATH, 以搜索到安装的 routiner.exe
SET PATH=%PATH%;%GOPATH%/bin
# 添加参数 -toolexec="routiner go-agent" -a
go build -toolexec="routiner go-agent" -a -o main.exe .
```

- linux

```shell
#!/bin/bash

# 将 %GOPATH%/bin 追加到 PATH, 以搜索到安装的 routiner
export PATH=$PATH:$GOPATH/bin
# 添加参数 -toolexec="routiner go-agent" -a
go build -toolexec="routiner go-agent" -a -o main.exe .
```
