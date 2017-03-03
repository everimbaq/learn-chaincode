# 安装部署
    0.6版本为准，1.0还未正式发布
## 安装docker
## 安装Fabric
从 docker hub 获取镜像
```bash
docker pull hyperledger/fabric-peer:latest
docker pull hyperledger/fabric-membersrvc:latest
```
## 创建docker配置文件
```yaml
membersrvc:
  image: hyperledger/fabric-membersrvc
  ports:
    - "7054:7054"
  command: membersrvc
vp0:
  image: hyperledger/fabric-peer
  ports:
    - "7050:7050"
    - "7051:7051"
    - "7053:7053"
  environment:
    - CORE_PEER_ADDRESSAUTODETECT=true
    - CORE_VM_ENDPOINT=unix:///var/run/docker.sock
    - CORE_LOGGING_LEVEL=DEBUG
    - CORE_PEER_ID=vp0
    - CORE_PEER_PKI_ECA_PADDR=membersrvc:7054
    - CORE_PEER_PKI_TCA_PADDR=membersrvc:7054
    - CORE_PEER_PKI_TLSCA_PADDR=membersrvc:7054
    - CORE_SECURITY_ENABLED=true
    - CORE_SECURITY_ENROLLID=test_vp0
    - CORE_SECURITY_ENROLLSECRET=MwYpmSRjupbT
  links:
    - membersrvc
  command: sh -c "sleep 5; peer node start --peer-chaincodedev"

```
    docker compose up -f docker-compose.yml
## 部署代码
### 说明

链码是在交易被部署时分发到网络上，并被所有验证 peer 通过隔离的沙箱来管理的应用级代码。
fabric的点对点（peer-to-peer）通信是建立在允许双向的基于流的消息gRPC上的。它使用Protocol Buffers来序列化peer之间传输的数据结构。
总账由两个主要的部分组成，一个是区块链，一个是世界状态。
区块链是在总账中的一系列连接好的用来记录交易的区块。世界状态是一个用来存储交易执行状态的键-值(key-value)数据库。
我们的主要工作就是写chaincode， 实现自己的业务逻辑。 共识、验证、加密等事情交给fabric去做 。

### Init()

Init 在首次部署你的链码时被调用。顾名思义，此函数用于链码所需的所有初始化工作。在我们的示例中，我们使用 Init 函数设置总账一个键值对的初始状态。

在你的 `chaincode_start.go` 文件中，修改 `Init` 函数，以便将 `args` 参数中的第一个元素存储到键 “hello_world” 中。

```go
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    return nil, nil
}
```

这是通过 stub 的 `stub.PutState` 函数完成的。该函数将部署请求中发送的第一个参数解释为要存储在分类帐中的键 “hello_world” 下的值。 这个参数是从哪里来的，什么是部署请求？我们将在实现接口后再解释。如果发生错误，例如传入的参数数量错误，或者写入总账时发生错误，则此函数将返回错误。否则，它将完全退出，什么都不返回。

### Invoke()

当你想调用链码函数来做真正的工作时，`Invoke` 就会被调用。这些调用会被当做交易被分组到链上的区块中。当你需要更新总账时，就会通过调用你的链码去完成。`Invoke` 的结构很简单。它接收一个 `function` 以及一个数组参数，基于调用请求中传递的 `function` 参数所指的函数，`Invoke` 将调用这个辅助函数或者返回错误。

在你的 `chaincode_start.go` 文件中，修改 `Invoke` 函数，让它调用一个普通的 `write` 函数。

```go
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // 处理不同的函数
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}
```

现在，它正在寻找 `write` 函数，让我们把这个函数写入你的 `chaincode_start.go` 文件。

```go
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }

    key = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(key, []byte(value))  //把变量写入链码状态中
    if err != nil {
        return nil, err
    }
    return nil, nil
}
```

你可能会认为这个 `write` 函数看起来类似 `Init`。 它们确实很像。 这两个函数检查一定数量的参数，然后将一个键/值对写入总账。然而，你会注意到，`write` 函数使用两个参数，允许你同时传递给调用的 `PutState` 键和值。基本上，该函数允许你向区块链总账上存储任意你想要的键值对。

### Query()

顾名思义，无论何时查询链码状态，`Query` 都会被调用。查询操作不会导致区块被添加到链中。你不能在 `Query` 中使用类似 `PutState` 的函数，也不能使用它调用的任何辅助函数。你将使用 `Query` 读取链码状态中键/值对的值。

在你的 `chaincode_start.go` 文件中，修改 `Query` 函数，让它调用一个普通的 `read` 函数，类似你对 `Invoke` 函数的修改。

```go
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // 处理不同的函数
    if function == "read" {                            //读取变量
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}
```

现在，它正在寻找 `read` 函数，让我们在 `chaincode_start.go` 文件中创建该函数。

```go
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}
```

这个 `read` 函数使用了与 `PutState` 作用相反的 `GetState`。 `PutState` 允许你对一个键/值对赋值，`GetState`允许你读取之前赋值的键的值。你可以看到，该函数使用的唯一参数被作为应该检索的值的键。接下来，此函数将字符数组返回到 `Query`，然后将它返回给 REST 句柄。

### Main()

最后，你需要创建一个简短的 `main` 函数，它在每个节点部署链码实例的时候执行。它仅仅调用了 `shim.Start()`，该函数会在链码与部署链码的节点之间建立通信。你不需要为该函数添加任何代码。`chaincode_start.go` 和 `chaincode_finished.go` 都有一个 `main` 函数，它位于文件的顶部。该函数如下所示：

```go
func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}