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
                that.$alert("请登录", '提示', {
                    confirmButtonText: '确定',
                    callback: function () {
                        window.location.replace("/login")
                    }
                });
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
        sleep: function (d) {
            for (var t = Date.now(); Date.now() - t <= d;);
        },
        ws: function () {
            var that = this
            var host = window.location.hostname + ":" + window.location.port
            var ws = new WebSocket("ws://" + host + "/ws?user_token=" + that.user_token)
            ws.onerror = function (err) {
                // that.$message.error("连接失败")
                // that.sleep(2000);
                //todo 尝试重连
                //连接失败,重新登陆
                that.$alert('websocket连接失败', '请重新登陆', {
                    confirmButtonText: '确定',
                    callback: function () {
                        window.location.replace("/login")
                    }
                });
            }

            ws.onmessage = function (event) {
                var resMes = JSON.parse(event.data)
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
                                    if (that.onlineUserList[i].name == resMes.from_user_name) {
                                        //下线
                                        that.$message(resMes.from_user_name + "下线了")
                                        that.onlineUserList.splice(i, 1)
                                    }
                                }
                            }
                        }
                        if (resMes.data == "online") {
                            //判断列表中是否有此用户
                            var isSelf = false
                            for (var i = 0; i < that.onlineUserList.length; i++) {
                                if (that.onlineUserList[i].name == resMes.from_user_name) {
                                    isSelf = true
                                    break
                                }
                            }
                            if (!isSelf) {
                                var user = {
                                    "name": resMes.from_user_name,
                                    "token": resMes.from_user_token
                                }
                                that.onlineUserList.push(user)
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
        //选择用户
        selectUser: function (token) {
            alert(token)
        }
    },
})