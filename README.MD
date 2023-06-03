

# TTPing

TTPing 是一个网络 ping 工具，可用于执行 TCP、ICMP 和 HTTP 的网络连接测试。它支持指定协议、连接尝试次数、数据包大小和采样间隔等参数。

## 作者和贡献者

- ttping 工具的作者是 ChatGPT

- chat-gpt 想对大家说的话
  ![示例图片](./chat-gpt-to-say.jpg)

## 安装

1. 通过go install安装
确保已经安装了Go编译器，并配置好了Go环境变量。
```shell
go install github.com/aiyijing/ttping
```

## 用法

```shell
ttping [address] [flags]
```

示例：

```shell
ttping example.com
```

### 参数

- `address`：要测试的目标地址。

### 标志

- `-t, --protocol`: 指定要使用的协议，有效值为 tcp、icmp 或 http，默认为 icmp。
- `-c, --count`: 指定连接尝试的次数，默认为 5。
- `-s, --size`: 指定要发送的数据包大小，默认为 64。
- `-i, --interval`: 指定采样间隔（秒），默认为 1。

### 示例

- 使用 TCP 协议进行连接测试：

  ```shell
  ttping example.com -t tcp
  ```

- 使用 HTTP 协议进行连接测试：

  ```shell
  ttping http://example.com -t http
  ```

- 指定连接尝试次数和数据包大小：

  ```shell
  ttping example.com -c 10 -s 128
  ```

## 贡献

欢迎贡献代码和提出问题。请查看 CONTRIBUTING.md 文件以获取更多信息。

## 许可证

这个项目基于 MIT 许可证。请阅读 LICENSE 文件以获取更多信息。