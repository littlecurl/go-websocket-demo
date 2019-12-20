### 说明
这里演示了使用Go语言实现WebSocket通讯的方式，最基本的api使用方式

在浏览器地址栏中输入地址无法请求websocket，因此需要建立前端页面，在html中使用JavaScript来建立WebSocket连接

#### 前端代码讲解

为了表达建立连接、断开连接的过程，使用了Open、Close两个按钮。而不是一打开页面就建立连接，关闭页面就断开连接。

另外，为了展示发送信息，接受返回信息，使用了一个input输入框和一个div进行占位。

最后，为了让界面效果稍微好看一点，使用了一个table，1 行 2 列，各占50%宽度。因此，Html代码如下：
```HTML
<body>
    <!-- 表格，1 行 2 列，各占 50%宽度 -->
    <table>
        <tr>
            <td valign="top" width="50%">
                <p>
                    Click "Open" to create a connection to the server,
                    <br />
                    Click "Send" to send a message to the server,
                    <br />
                    Click "Close" to close the connection.
                    <br />
                    You can change the message and send multiple times.
                </p>
                <form>
                    <button id="open">Open</button>
                    <button id="close">Close</button>
                    <input id="input" type="text" value="Hello world!">
                    <button id="send">Send</button>
                </form>
            </td>
            <td valign="top" width="50%">
                <div id="output"></div>
            </td>
        </tr>
    </table>
</body>
```


有了前端界面，剩下的就是编写JS代码了。

首先要做的是监听页面加载事件window.addEventListener("load", function)

因为至少要在open、close两个方法中使用到ws对象，所以把ws声明为全局对象。
```javascript
// 声明变量ws
var ws;
```
open按钮响应事件如下
```javascript
 // 和服务器建立连接
document.getElementById("open").onclick = function (evt) {
    // 如果ws对象存在，则不进行创建（避免多次点击open，导致多次建立连接）
    if (ws) {
        return false;
    }
    // 建立ws连接，注意不能只写ws://localhost:7777，后面还有一个 /ws 呢
    ws = new WebSocket("ws://localhost:7777/ws");
    // ws的4个事件
    ws.onopen = function (evt) {
        print("OPEN");
    };
    ws.onclose = function (evt) {
        print("CLOSE");
        ws = null;
    };
    ws.onmessage = function (evt) {
        print("RESPONSE: " + evt.data);
    };
    ws.onerror = function (evt) {
        print("ERROR: " + evt.data);
    };
    // 中断监听事件
    return false;
};
```
当点击open()按钮后，首先要判断一下当前线程上是否已经初始化了ws对象，如果已经初始化了，就无需再次进行初始化。

如果未初始化，就调用new WebSocket("ws://hosts地址");进行建立连接。

紧接着需要注册4个事件：建立连接时onopen()、关闭连接时onclose()、接收到消息时onmessage()、遇到错误时onerror()

最后，中断当前按钮点击事件 return false;

close按钮响应事件如下：
```javascript
// 关闭按钮
document.getElementById("close").onclick = function (evt) {
    // 如果ws对象不存在，直接返回
    if (!ws) {
        return false;
    }
    ws.close();
    return false;
};
```
close按钮最主要的就是调用了ws.close();

还有一个send按钮响应事件：
```javascript
// 发送按钮
document.getElementById("send").onclick = function (evt) {
    // 如果ws对象不存在，直接返回
    if (!ws) {
        return false;
    }
    print("SEND: " + input.value);
    ws.send(input.value);
    return false;
};
```
send按钮最主要的就是调用了ws.send();

总结：
ws的 7 个方法
```javascript
// 1个构造方法
var ws = new WebSocket("ws://hosts地址");
// 4个事件监听方法
ws.onopen();
ws.onclose();
ws.onmessage();
ws.onerror();
// 1个向服务端发送消息方法
ws.send();
// 1个与服务器断开连接方法
ws.close();
```

#### 后端代码讲解

首先看main方法中
```go
func main() {
    // 监听地址及端口
	_ = http.ListenAndServe("0.0.0.0:7777", nil)
	// 通过wsHandler中的upgrade升级为WebSocket长连接
	http.HandleFunc("/ws", wsHandler)
}
```
其中的http是Go语言的标准库，先调用ListenAndServe()方法监听服务地址及端口，紧接着调用HandleFunc方法，传递一个pattern以及一个回调函数。

所谓pattern就是你要访问的除主机之外的地址。

在回调函数里，首先定义了一个upgrade对象，并在upgrade对象内部允许跨域请求访问。接着调用upgrade.Upgrade(w, r, nil);方法来建立WebSocket连接。

之所以叫upgrade是因为，WebSocket连接时HTTP连接的升级版（具体可以看HTTP协议的请求头与响应体的upgrade字段）。

握手之后，写个死循环，不断地读取数据，写回数据。

其中读写操作都用到了*websocket.Conn类型的对象conn

读方法：conn.ReadMessage();

写方法：conn.WriteMessage();

这样就实现了转发客户端消息的功能。