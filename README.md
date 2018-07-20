# SAC 
SAC：全称为Smart Application Chain, 即智能应用链，旨在打造区块链与人工智能链是智能应用链的应用生态系统。本开源项目是SAC链智能应用链的实现和命令行接口。

## 性质
当前SAC归类为联盟链，SAC链主要挖矿节点归深圳阿尔法巫师科技有限公司所有，其他节点可以通过下载源码运行自行接入该链，接入节点不拥有挖矿权限。另外，参与节点可通过授权获取挖矿权限与主要节点形成利益相关联盟链，共同维护SAC区块链的正常运行。

## 搭建环境
### 安装
该项目是基于GO语言进行编译运行的一套程序，故所有节点需要搭建一套GO语言的开发环境并且安装gcc编译器。

### 编译
| Command    | Description |
|:----------:|-------------|
| **`Windows`** |安装git，进入cmd/geth目录，执行go build ./... 。相应的可执行文件在该目录下生成（即geth.exe）。|
| `Linux` |进入根目录执行make，编译完毕可在builg/bin目录下生成可执行文件geth。|



