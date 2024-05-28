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

## JVM

JVM (Java Version Manager). Java 版本管理，并不是Java虚拟机的意思。

### 设置JAVA_HOME

设置JAVA_HOME环境变量的值，也就是生成的快捷方式的路径。默认使用`安装目录\runtime\jdk`

```powershell
jvm home [path]  # 例如 jvm home D:\soft\jdk
```

### 添加版本

```powershell
jvm add [key] [path]  # 例如 jvm add jdk-11 D:\soft\java\jdk-11.0.2
```

### 切换版本

```powershell
jvm use # 交互式选择使用的版本
```

### 删除版本

```powershell
jvm rm [key]  # 例如 jvm rm jdk-11
```

### 查看版本

```powershell
jvm list  # 查看所有添加的版本
```

### 在线安装JDK

```powershell
jvm install # 交互式选择安装的版本
```


