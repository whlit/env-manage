# Env Manage

一个在管理开发环境版本的简单工具。从[nvm-windows](https://github.com/coreybutler/nvm-windows)学习并改造而成。

原理是创建一个文件夹快捷方式，指向所安装的文件夹路径，当需要修改环境时，修改快捷方式的指向即可，这样在终端中的环境变量虽然还是原来的，但是实际路径其实已经变了。

例如：`JAVA_HOME=D:\soft\jdk`环境变量中是这样设置的，但是`jdk`这个文件夹其实是个快捷方式，指向`D:\soft\java\jdk-11.0.2`，在终端中打印`echo %JAVA_HOME%`，结果是`D:\soft\jdk\`，但使用`java -version`，结果是`java version "11.0.2"`。此时修改快捷方式的指向，`D:\soft\jdk\`指向`D:\soft\java\jdk-17.0.3`，`java -version`结果就是`java version "17.0.3"`。但终端中的环境变量`JAVA_HOME`还是`D:\soft\jdk\`。所以依赖于`JAVA_HOME`的软件(比如Maven)都会使用对应版本的JDK。

## Installation

解压后安装目录下有`install.exe`和`uninstall.exe`,用于作为全局命令安装和卸载;

也可自行构建可执行文件，构建的包在dist目录下

```powershell
./build.cmd
```

## VM

VM (Version Manager). 软件版本环境变量管理器

```
Usage: vm <name> <action> [args]
Name:                      环境管理的名称
  jdk                      jdk版本管理
  maven                    maven版本管理
  <name>                   用create创建的其他版本管理的名称
Actions:
  create                   创建一个版本管理
  add <version> <path>     添加版本,version: 版本名称(自定义),path: 版本的绝对路径
  rm                       移除版本
  list                     查询所有已添加的版本管理
  use                      使用版本
  install                  在线安装新版本,只支持jdk/maven的在线安装
```

默认支持JDK和MAVEN的版本管理，同时提供自定义版本管理

### 自定义版本管理

Usage: `vm <name> create`

例如：添加golang版本管理

- 创建版本管理

```sh
vm go create
```

之后会弹出输入要创建的环境变量及要添加到Path的值

```
输入要添加的环境变量名称:
(例如:JDK环境变量名称为:JAVA_HOME)
> GOROOT
输入要添加到Path的值:
(例如:JDK添加到Path为:%JAVA_HOME%\bin)
> %GOROOT%\bin
```

- 添加版本

```sh
vm go add v1.21.4 D:\soft\go
```

- 查询版本信息

```sh
vm go list
```

- 使用版本

```sh
vm use
````

- 移除版本

```sh
vm go rm v1.21.4
```

只是从列表中移除，并不会删除文件


## JVM

JVM (Java Version Manager). Java 版本管理，并不是Java虚拟机的意思。这个命令是上面命令的简化版本。

### 设置JAVA_HOME

设置JAVA_HOME环境变量的值，也就是生成的快捷方式的路径。默认使用`安装目录\runtime\jdk`

```powershell
jvm home <path>  # 例如 jvm home D:\soft\jdk
```

### 添加版本

```powershell
jvm add <key> <path>  # 例如 jvm add jdk-11 D:\soft\java\jdk-11.0.2
```

### 切换版本

```powershell
jvm use # 交互式选择使用的版本
```

### 删除版本

```powershell
jvm rm <key>  # 例如 jvm rm jdk-11
```

### 查看版本

```powershell
jvm list  # 查看所有添加的版本
```

### 在线安装JDK

直接在线安装JDK，支持从Oracle官网下载和Adoptium官网下载

```powershell
jvm install # 交互式选择安装的版本
```

## MVM

MVM (Maven Version Manager). Maven 版本管理器。

### 设置M2_HOME

设置M2_HOME环境变量,默认使用`安装目录\runtime\jdk`

```powershell
mvm home <path>  # 例如 mvm home D:\apache-maven
```

### 添加版本

```powershell
mvm add <key> <path>  # 例如 mvm add apache-maven-3.9.4 D:\soft\apache-maven-3.9.4
```

### 切换版本

```powershell
mvm use # 交互式选择使用的版本
```

### 删除版本

```powershell
mvm rm <key>  # 例如 mvm rm apache-maven-3.9.4
```

### 查看版本

```powershell
mvm list  # 查看所有添加的版本
```

### 在线安装Maven

```powershell
mvm install # 交互式选择安装的版本
```

