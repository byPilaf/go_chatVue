var app = new Vue({
    delimiters: ['[{', '}]'],
    el: "#app",
    mounted() {
        this.checkUser()
        this.getOnlineUserList()
    },
    data: {
        //弹出框
        dialogVisible: false,
        //提示信息
        message: "",
        //在线用户列表
        onlineUserList: [],
        //user_token
        user_token: "",
        //输入的内容
        input: "",
        //收到的消息列表
        reUserMes: [],
        flag: true,
    },
    methods: {
        checkUser: function () {
            var that = this
            //读取cookie
            var user_token = document.cookie
            user_token = user_token.split("=")[1];
            if (user_token) {
                //有cookie, 请求服务段获取用户数据
                axios.post("", { "user_token": user_token })
                    .then(function (response) {
                        if (response.data.code != "200") {
                            that.$alert(response.data.data, '提示', {
                                confirmButtonText: '确定',
                                callback: function () {
                                    window.location.replace("/login")
                                }
                            });
                            // that.dialogVisible = true
                            // that.message = response.data.data
                        } else {
                            that.user_token = user_token
                            that.ws() //连接websocket服务器as
                        }
                    }, function (err) { })
            } else {
                that.dialogVisible = true
                that.message = "未登录"
            }
        },
        tologin: function () {
            window.location.replace("/login")
        },
        getOnlineUserList: function () {
            var that = this
            //获取当前在线用户列表
            axios.get("/getOnlineUserList")
                .then(function (response) {
                    if (response.data.code == "200") {
                        that.onlineUserList = response.data.data
                    }
                }, function (err) { })
        },
        ws: function () {
            var that = this
            while (that.flag) {
                that.sleep(5000); // 延时函数，单位ms
                var ws = new WebSocket("ws://127.0.0.1:7999/ws?user_token=" + that.user_token)
                ws.onerror = function (err) {
                    //连接失败,重新登陆
                    //尝试重新连接
                    that.$message.error("err websocket")
                    console.log(err)
                    // that.$alert('请重新登陆', '提示', {
                    //     confirmButtonText: '确定',
                    //     callback: function () {
                    //         window.location.replace("/login")
                    //     }
                    // });
                }

                ws.onmessage = function (event) {
                    console.log("asd")
                    var resMes = JSON.parse(event.data)
                    that.flag = false
                    switch (resMes.mes_type) {
                        case 0: //系统消息
                            that.$message("系统消息: " + resMes.data)
                            break
                        case 1: //用户状态消息
                            if (resMes.data == "offline") {
                                if (resMes.from_user_token == that.user_token) {
                                    //自己下线了
                                    //连接失败,重新登陆
                                    that.$alert('您已经在其他地方登录', '提示', {
                                        confirmButtonText: '确定',
                                        callback: function () {
                                            window.location.replace("/login")
                                        }
                                    });
                                } else {
                                    for (var i = 0; i < that.onlineUserList.length; i++) {
                                        if (that.onlineUserList[i] == resMes.from_user_name) {
                                            //下线
                                            that.$message(resMes.from_user_name + "下线了")
                                            that.onlineUserList.splice(i, 1)
                                        }
                                    }
                                }
                            }
                            if (resMes.data == "online") {
                                if (that.onlineUserList.indexOf(resMes.from_user_name) == -1) {
                                    that.onlineUserList.push(resMes.from_user_name)
                                    if (that.user_token != resMes.from_user_token) {
                                        that.$message(resMes.from_user_name + "上线了")
                                    }
                                }
                            }
                            break
                        case 2: //聊天消息
                        case 3: //群聊消息
                            that.reUserMes.push(resMes)
                            that.scrollToBottom()
                            break
                        case 4: //隐藏消息
                            //心跳检测
                            if (resMes.code == 200) {

                            }
                            break
                    }
                }
            }
        },
        sleep: function (d) {
            for (var t = Date.now(); Date.now() - t <= d;);
        },
        sendMes: function () {
            var that = this
            //判断内容
            if (that.input != "") {
                //发送的消息
                var sendData = {
                    "from_user_token": that.user_token,
                    "mes_type": 3,
                    "data": that.input,
                }
                axios.post("/sendMes", sendData)
                    .then(function (response) {
                        if (response.data.code == "200") {
                            that.input = ""
                        } else if (response.data.code == "401") {
                            that.dialogVisible = true
                            that.message = response.data.data
                        } else {
                            that.$message.error(response.data.data + ", 发送失败")
                        }
                    }, function (err) {
                        that.$message.error("发送失败")
                    })
            }
        },
        //自动滚动到底部
        scrollToBottom: function () {
            this.$nextTick(() => {
                var container = this.$el.querySelector(".chat_box");
                container.scrollTop = container.scrollHeight;
            });
        },
    },
})