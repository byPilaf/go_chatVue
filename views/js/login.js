var app = new Vue({
    delimiters: ['[{', '}]'],
    el: "#app",
    data: {
        //登陆信息
        userInfo: {
            name: "",
            pass: "",
        },
        //弹出框
        dialogVisible: false,
        //提示信息
        message: "",
    },
    methods: {
        login: function () {
            var that = this
            axios.post("", that.userInfo)
                .then(function (response) {
                    //登陆成功, 跳转主页 or 登陆失败,展示错误信
                    that.dialogVisible = true
                    that.message = response.data.data
                    if (response.data.code == 200) {
                        document.cookie = "user_token=" + response.data.from_user_token;
                        window.location.replace("/")
                    }
                }, function (err) {
                    console.log(err)
                })
        },
    },
})