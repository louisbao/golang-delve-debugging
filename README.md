## Golang Delve Debugging 调试教程

### 1. Delve 安装

如果操作系统是 `OSX` 请先安装开发命令行编译工具：

```bash
xcode-select --install
```

确保已经配置过 `$GOPATH` 环境变量，具体教程可以参考 [Golang 官方文档](https://golang.org/doc/install)。

运行以下命令就可以在 `$GOPATH/bin` 目录下安装 `delve` 的可执行程序 `dlv`：

```bash
go get -u github.com/derekparker/delve/cmd/dlv
```

如果你的 go 版本小于`1.5`，请先设置环境变量 `GO15VENDOREXPERIMENT=1` 再运行以上命令进行安装。

然后就可以运行 `dlv` 命令启动 `delve`，例如执行以下命令行：

```bash
dlv exec ./greeting -- server --config conf/config.toml
```

需要传输到`被调试程序`的各种参数可以通过 `--` 来传入，比如这里我们给 `greeting` 程序传入了 `server --config conf/config.toml`。

### 2. Delve 常用子命令详解

| 子命令 SubCommand | 描述 Description | 用法 Usage |
| ---------------- | ---------------- | ---------------- |
| dlv attach | 挂载(attach)到运行中的进程并开始调试 | dlv attach pid [executable] [flags] |
| dlv connect | 连接到一个运行中的非交互模式 delve 调试服务器 | dlv connect addr [flags] |
| dlv core | 检查一个 core dump | dlv core <executable> <core> [flags] |
| dlv debug | 编译并且开始调试当前文件夹下的 main package，或者指定的 package | dlv debug [package] [flags] |
| dlv exec | 执行一个已编译的可执行程序后，立即 attach 到该进程开始调试 | dlv exec <path/to/binary> [flags] |
| dlv run | Deprecated command. Use 'debug' instead | dlv run [flags] |
| dlv test | 编译 test 可执行程序并且开始调试 | dlv test [package] [flags] |
| dlv trace | 编译并且开始 trace 程序 | dlv trace [package] regexp [flags] |

#### 2.1 `dlv debug` 命令详解

`dlv debug` 命令其实就是在当前目录下临时编译生成一个没有代码优化的可执行文件，文件名是 `debug`。接下来就对这个 `debug` 文件进行调试。而在 `dlv` 程序退出后，会删除这个 `debug` 文件。

例如，我们有一个项目的结构如下：

```bash
.
├── github.com/me/foo
├── cmd
│   └── foo
│       └── main.go
├── pkg
│   └── baz
│       ├── bar.go
│       └── bar_test.go
```

如果当前所在的文件夹位于 `github.com/me/foo/cmd/foo`，你可以直接执行 `dlv debug` 开始调试工作。如果当前所在的位置在其他地方，例如项目的根目录，则需要指定 main package 所在的位置: `dlv debug github.com/me/foo/cmd/foo`，同样`被调试程序`的各种参数可以通过 `--` 来传入:

```bash
dlv debug github.com/me/foo/cmd/foo -- -arg1 value
```

#### 2.2 `dlv attach` 命令详解

挂载(attach)到一个正在运行在 `8080` 端口上的 go web 服务，启动 `hello` 程序，然后访问 `localhost:8080\hello` 会返回机器的名称:

```bash
$ cd examples/hello
$ go build -o hello main.go
$ ./hello
# 通过 ps 命令查看对应程序的 pid
$ ps aux | grep hello
louis            46960   0.0  0.0  4287476    840 s000  S+    8:50PM   0:00.01 grep hello
louis            46958   0.0  0.1 558459028   7092 s011  S+    8:50PM   0:00.03 ./hello
$ dlv attach 46958
```

### 3. Delve 常用的调试命令

进入 `examples/apple` 文件夹，查看第一个例子：

```bash
cd examples/apple
```

执行 `dlv debug` 后 `delve` 将会自动在当前文件夹找到 `main.go` 所在的包：

```bash
$ dlv debug
Type 'help' for list of commands.
(dlv) break main.main
Breakpoint 1 set at 0x10d640b for main.main() ./main.go:20
(dlv) continue
> main.main() ./main.go:20 (hits goroutine(1):1 total:1) (PC: 0x10d640b)
    15:
    16:	func (p *Product) isObsolete() {
    17:		p.Obsolete = true
    18:	}
    19:
=>  20:	func main() {
    21:		products := []Product{
    22:			{
    23:				Name: "iPhoneXS",
    24:				Year: "2018",
    25:			},
(dlv) continue
[{"name":"iPhoneXS","year":"2018","obsolete":false},{"name":"Macbook","year":"2018","obsolete":false},{"name":"iBook","year":"2006","obsolete":true}]
```

这里 `break` 和 `continue` 都是我们常用的调试命令，一般在进入调试交互界面以后，可以使用的命令有以下这些：

Command | Description
--------|------------
[args](#args) | 打印函数参数
[break](#break) | 设置一个断点
[breakpoints](#breakpoints) | 打印激活的断点信息
[call](#call) | 回复进程，注入一个函数调研function call
[check](#check) | 在当前位置创建一个检查点checkpoint
[checkpoints](#checkpoints) | 打印所有checkpoints
[clear](#clear) | 删除断点
[clear-checkpoint](#clear-checkpoint) | 删除checkpoint
[clearall](#clearall) | 删除所有的断点
[condition](#condition) | 设置断点条件
[config](#config) | 修改配置参数
[continue](#continue) | 运行到断点或程序终止
[deferred](#deferred) | 在deferred call环境执行一个命令
[disassemble](#disassemble) | 拆解器
[down](#down) | 往下移动当前frame
[edit](#edit) | 用$DELVE_EDITOR或$EDITOR打开你所在的位置
[exit](#exit) | 退出 debugger
[frame](#frame) | 在不同的框架上执行的命令
[funcs](#funcs) | 打印函数列表
[goroutine](#goroutine) | 显示或更改当前goroutine
[goroutines](#goroutines) | 列出程序的全部goroutines
[help](#help) | 打印出帮助信息
[list](#list) | 显示源代码
[locals](#locals) | 打印局部变量
[next](#next) | 跳到下一行
[on](#on) | 在遇到断点时执行一个命令
[print](#print) | 评估表达式
[regs](#regs) | 打印CPU寄存器的内容
[restart](#restart) | 重启进程
[rewind](#rewind) | 执行backwards到breakpoint或到程序结束
[set](#set) | 更改变量的值
[source](#source) | 执行包含delve命令列表的文件
[sources](#sources) | 打印源文件列表
[stack](#stack) | 打印堆栈跟踪
[step](#step) | 单步执行程序
[step-instruction](#step-instruction) | 单步单个执行cpu指令
[stepout](#stepout) | 退出当前的function
[thread](#thread) | 切换到指定的线程
[threads](#threads) | 打印每一个跟踪线程的信息
[trace](#trace) | 设置跟踪点
[types](#types) | 打印类型列表
[up](#up) | 往上移动当前frame
[vars](#vars) | 打印某个包内的(全局)变量
[whatis](#whatis) | 打印表达式类型

### 4. 在 `API` 模式下运行 `Delve`

Delve 可以运行在 `API` 模式下，这样 IDE 编辑器就可以通过接口与之进行交互。只需在标准命令之后提供 `--headless` 标志，如下所示：

```bash
$ dlv debug --headless --api-version=2 --log --listen=127.0.0.1:8181
```

这将以非交互模式启动调试器，侦听指定的地址和端口，并启用日志记录。当然最后两个标志(log和listen)是可选的。
如果您需要将多个客户端连接到API，也可以指定 `--accept-multi` 客户端标志。然后使用 `connect` 子命令从 `Delve` 本身连接 `headless` 调试器，这对于远程调试很有用：

```bash
$ dlv connect 127.0.0.1:8181
```

### 5. 结合 `Visual Studio Code` 进行调试

用 `VS Code` 打开 `examples\hello` 项目文件夹，可以看到 `.vscode` 文件夹内有一个配置文件 `launch.json`。

打开 `launch.json` 配置文件：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceRoot}",
      "env": {},
      "args": [],
      "showLog": true
    }
  ]
}
```

如果没有找到该文件，也可以通过在 `VS Code` 界面运行命令 `Debug: Open launch.json` 来创建并修改为以上内容。点击 `F5` 就可以开始进行调试。

`launch.json` 配置文件的主要配置规则如下：

一、`name` 是整个配置的名称，可以随意配置；

二、`program` 配置项是必须填写的，值可以是一个文件夹或文件的路径：

* 可以是希望被调试的 package 的路径，或者 package 的路径内的一个文件；
* 该值必须为绝对路径，不能为相对路径；
* 使用 `${workspaceFolder}` 值可以自动探测到 VS Code 当前打开的工作项目 `workspace` 根目录包含的 package，并进行调试；
* 使用 `${file}` 值可以自动探测到 VS Code 当前打开的文件，并进行调试；
* 使用 `${fileDirname}` 值可以自动探测到 VS Code 当前打开的文件所属的 package，并进行调试。

三、`mode` 配置项主要配置如下：

* 使用 `debug` 值可以编译程序，并开启调试器 [default]；
* 使用 `test` 值可以调试测试（test）代码；
* 使用 `auto` 值可以根据当前打开的文件判断是否使用 `debug` 或 `test`；
* 使用 `exec` 值可以运行预编译的二进制程序，并 attach 到该程序的进程，例如 `"program":"${workspaceRoot}/mybin"`；
* 使用 `remote` 值可以 attach 到一个远程的非交互式 delve 服务器. 你必须先手动在远程的机器上运行一个非交互式 delve 服务器，并且在本地连接远程服务器时提供额外的配置参数 `remotePath`, `host` 以及 `port`。

四、`env` 配置项可以用来配置调试启动时所用的环境变量参数：

* 比如 `$GOPATH` 临时设置为某个参数就可以在这里指定；
* 该配置的基本格式如下：

```json
"env": {
  "ENV1": "envValue1",
  "ENV2": "envValue2",
},
```

五、`args` 配置项可以用来向`被调试程序`传入各种参数：

* 这个 `args` 配置项功能类似命令行下的 `--` 标示；
* 该配置的基本格式如下：

```json
"args": [
  "--verbosity",
  "6",
  "--datadir",
  "dataDir"
]
```

六、远程调试的配置项有 `remotePath`, `port` 和 `host` 需要同时配置，详情见第 6 节。

七、`showLog` 配置项设为 `true` 可以打印 delve 执行的日志。

八、`trace` 配置项设为 `verbose` 可以打印 delve 执行的更详细的日志，在控制台会输出这个日志的目录。

### 6. 结合 `Visual Studio Code` 进行远程调试

用 `Visual Studio Code` 进行远程调试，你需要先在远程目标机器上运行一个非交互模式的 delve 服务器，例如：

```bash
$ dlv debug --headless --listen=:2345 --log
```

`被调试程序`的各种参数可以通过 `--` 来传入，例如：

```bash
$ dlv debug --headless --listen=:2345 --log -- -myArg=123
```

> 注意: 请不要传入 `–api-version=2` 参数, `Visual Studio Code` 的 `Go extension` 还不支持 delve API 2 的版本.

在 `Visual Studio Code` 创建一个远程的 `launch.json` 配置：

```json
{
	"name": "Remote",
	"type": "go",
	"request": "launch",
	"mode": "remote",
	"remotePath": "${workspaceRoot}",
	"port": 2345,
	"host": "127.0.0.1", # 这个例子远程 delve 服务器和 VS Code 编辑器在同一台机器
	"program": "${workspaceRoot}",
	"env": {}
}
```

> 注意: `remotePath` 值必须是正确的，否则 breakpoint 不会有任何提示。

当你启动调试器以后，`VS Code` 将会发送调试命令到 `dlv` 服务器，而非在本地针对你的代码程序启动一个 `dlv` 实例。

### 7. `settings.json` 配置文件

以下这些配置项也可以被调试器使用，你可能并不需要去修改这些配置，但最好能有所了解：
- `go.gopath`. 见 [GOPATH in VS Code](https://github.com/Microsoft/vscode-go/wiki/GOPATH-in-the-VS-Code-Go-extension)
- `go.inferGopath`. 见 [GOPATH in VS Code](https://github.com/Microsoft/vscode-go/wiki/GOPATH-in-the-VS-Code-Go-extension)
- `go.delveConfig`
     - `apiVersion`: 可用来控制 delve apis 的启动版本，默认为 2
     - `dlvLoadConfig`: 当 `apiVersion` 为 1 时无效，可以影响变量在调试面板的展示，主要的配置有：
         - `maxStringLen`:  maximum number of bytes read from a string
         - `maxArrayValues`:  maximum number of elements read from an array, a slice or a map
         - `maxStructFields`:  maximum number of fields read from a struct, -1 will read all fields
         - `maxVariableRecurse`:  how far to recurse when evaluating nested types
         - `followPointers`:  requests pointers to be automatically dereferenced

### 8. Official Delve Documentation 官方文档

- [Installation](https://github.com/derekparker/delve/tree/master/Documentation/installation)
  - [Linux](https://github.com/derekparker/delve/blob/master/Documentation/installation/linux/install.md)
  - [macOS](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md)
  - [Windows](https://github.com/derekparker/delve/blob/master/Documentation/installation/windows/install.md)
- [Getting Started](https://github.com/derekparker/delve/blob/master/Documentation/cli/getting_started.md)
- [Documentation](https://github.com/derekparker/delve/tree/master/Documentation)
  - [Command line options](https://github.com/derekparker/delve/blob/master/Documentation/usage/dlv.md)
  - [Command line client](https://github.com/derekparker/delve/blob/master/Documentation/cli/README.md)
  - [Plugins and GUIs](https://github.com/derekparker/delve/blob/master/Documentation/EditorIntegration.md)
- [Contributing](https://github.com/derekparker/delve/blob/master/CONTRIBUTING.md)
  - [Internal documentation](https://github.com/derekparker/delve/blob/master/Documentation/internal)
  - [API documentation](https://github.com/derekparker/delve/tree/master/Documentation/api)
  - [How to write a Delve client](https://github.com/derekparker/delve/blob/master/Documentation/api/ClientHowto.md)

Delve is a debugger for the Go programming language. The goal of the project is to provide a simple, full featured debugging tool for Go. Delve should be easy to invoke and easy to use. Chances are if you're using a debugger, things aren't going your way. With that in mind, Delve should stay out of your way as much as possible.
