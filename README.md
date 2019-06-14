## 简介

golang websocket 服务可通过http接口push消息到web客户端

## 安装

- 导入db.sql 安装相关push日志表
- 更改config.dev.json中的相关db配置与项目路径配置
- 执行install_package.sh 安装相关包依赖
- 执行go run build 编译并运行程序 


## 注意

本项目websocket使用用户token验证连接，这一块验证的逻辑需要根据自己的业务去更改或删除。

## 使用

例如使用PHP客户端push消息 

```php

<?php
/**
 * Created by PhpStorm.
 * Date: 2019/5/23
 * Time: 17:03
 */

require_once 'XHXPushApi.php';

$title = "通知标题pxy";
$content = "测试通知内容，测试通知内容，测试通知内容，测试通知内容，测试通知内容，测试通知内容。";

$res = XHXPushApi::getInstance()->push(1,'god', 2, $title, $content, [44, 63]);

print_r($res);

```

client.html

```
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <!-- import CSS -->
  <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
</head>
<body>
  <div id="app">
    
  </div>
</body>
  <!-- import Vue before Element -->
  <script src="https://unpkg.com/vue/dist/vue.js"></script>
  <!-- import JavaScript -->
  <script src="https://unpkg.com/element-ui/lib/index.js"></script>
  <script>
    new Vue({
      el: '#app',
      data() {
        return { 
            visible: false,
            ws:'',
            interval:'',
            retryConnect:false,
        }
      },
      created() {
          this.init()
      },
      methods: {
        init() {
            if (!window["WebSocket"]) {
                console.log('not support websocket')
                return
            }

            var that = this;
            this.ws = new WebSocket("ws://127.0.0.1:9002/ws/");
            this.ws.onclose = function(e) {
                clearInterval(that.interval)
                if(!that.retryConnect) {
                    return
                }
                console.log('push connection is close, retry connect after 5 seconds')
                setTimeout(function() {
                    that.init()
                }, 5000);
            }
            this.ws.addEventListener('open', function (e) {
                //登录
                that.ws.send('{"event":"register", "token":"00000063_d10f2dd30c087a0573d54e4767640253279"}');
            });

            this.ws.addEventListener("message", function(e) {
                let res = JSON.parse(e.data)
                
                //token过期
                if(res.error == 100) {
                    console.log(res)
                    that.retryConnect = false
                    return
                }

                if(res.error != 0) {
                    console.log(res.msg)
                    return
                }
                
                //client注册消息
                if(res.event == 'register') {
                    console.log('ws connection register success ')
                    that.interval = setInterval(function() {
                        //保此常连接心跳
                        that.ws.send('{}')
                    }, 60000)
                    that.retryConnect = true
                    return;
                }

                if(res.event == 'message') {
                    let options = JSON.parse(res.data.options);
                    that.$notify.info({
                        title: res.data.title != '' ? res.data.title : '通知',
                        message: res.data.content,
                        duration: options.duration,
                        position: options.position
                    });
                }
            })
        }
      }
    })
  </script>
</html>
```
