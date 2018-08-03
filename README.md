# SAC 
SAC：全称为Smart Application Chain, 即智能应用链，旨在打造基于区块链与人工智能的应用生态系统。本开源项目是SAC链智能应用链的实现和命令行接口。

## 性质
当前SAC归类为测试公链，SAC链主要挖矿节点归SAC基金会所有，其他节点可以通过下载源码运行自行接入该链，接入节点不拥有挖矿权限。另外，参与节点可通过授权获取挖矿权限与主要节点形成利益相关链，共同维护SAC区块链的正常运行。

## 搭建环境
### 安装
该项目是基于GO语言进行编译运行的一套程序，故所有节点需要搭建一套GO语言的开发环境并且安装gcc编译器。

### 编译
| 系统       |  说明        |
|:----------:|-------------|
| **`Windows`** |安装git，进入cmd/geth目录，执行go build ./... 。相应的可执行文件在该目录下生成（即geth.exe）。|
| `Linux` |进入根目录执行make，编译完毕可在builg/bin目录下生成可执行文件geth。|

### 配置创世文件
SAC链是基于以太坊所编写的一套区块链系统，当前创世区块的所有字段与以太坊的保持一致，读者可以自行参考以太坊相关创世区块文件说明。
SAC拥有独立的网络号以区分其他链，网络号为2673，成功接入网络的条件之一是配置正确的创世文件。SAC唯一创世文件如下（需新建genesis.json）：
```json
{
   "config": {
        "chainId": 2673,
        "homesteadBlock": 0,
        "eip155Block": 0,
        "eip158Block": 0
    },
    "alloc": {
        "0x2a62740d4767b70d7144e6ea232188c5f7cef44e": {
            "balance": "200000000000000000000000000"
        }
    },
    "nonce": "0x0000000000000042",
    "difficulty": "0x0280000",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "coinbase": "0x2a62740d4767b70d7144e6ea232188c5f7cef44e",
    "timestamp": "0x00",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "extraData": "",
    "gasLimit": "0xffffffff"
}
```

用刚初始化创世区块进行初始化,初始化命令如下
```
geth init genesis.json --datadir data0
```

### 配置静态文件
为了更好接入SAC网络，建议用户配置节点静态文件。初始化创世节点，生成相应的数据库文件（如上命令的data0），进入该文件创建static-nodes.json文件
```json
["enode://298089e66789dc51e302dcc921951f912e21e0a60ed8b054d08433d70f01670e9108a526a3a4905f0bb598c5c6d87956b17043d72ef3ce593f64224a08e0c4e1@112.74.43.58:30303?discport=0","enode://125aef2132619a846f2fd41a7aabd400cb78dfd02dc64ceb4b3176e5931d26fa8b46b8ce256aa4714c5d12fde041e3fb4ad6c3ea0233016990a3d4bdc015f4b6@120.79.36.94:30303?discport=0"]
```

### 运行
```
geth --nodiscover  --rpc  --datadir data0 --port "30303" --rpcapi "db,eth,net,web3,miner,personal" --rpcport "8545" --rpcaddr "localhost" --networkid "2673" console 2>>chain.log
```

| 字段       |  说明        |
|:----------:|-------------|
| **`nodiscover`** |安装git，进入cmd/geth目录，执行go build ./... 。相应的可执行文件在该目录下生成（即geth.exe）。|
| `rpc` |进入根目录执行make，编译完毕可在build/bin目录下生成可执行文件geth。
| `datadir `|设置当前区块链网络数据存放的位置|
| `port `|网络监听端口|
| `rpcapi `|设置允许连接的rpc的客户端|
| `rpcport`|设置允许连接的rpc的客户端|
| `rpcaddr`|rpc接口的地址 |
| `networkid`|设置当前区块链的网络ID，用于区分不同的网络|
| `console`|启动命令行模式，|

### 命令
当前所有命令与以太坊相关命令一致，用户可以参考以太坊命令。参考，[以太坊基本命令](https://blog.csdn.net/wo541075754/article/details/53073799)

## 贡献
欢迎广大区块链爱好者人加入我们的项目，参与SAC链的迭代开发。












