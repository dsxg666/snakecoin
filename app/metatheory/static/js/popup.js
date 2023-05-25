document.addEventListener('DOMContentLoaded', () => {
    let keyInput = document.getElementById('keyInput');
    let valueInput = document.getElementById('valueInput');
    let saveButton = document.getElementById('saveButton');
    let getButton = document.getElementById('getButton');
    let resultDiv = document.getElementById('result');
    let copyButton = document.getElementById('copyButton');

    let b = $("#mining");
    b.click(() => {
        b.text("挖矿中...");
        $.ajax({
            url: "http://localhost:8080/mine",
            method: "Post",
            data: {},
            success: function (response) {
                b.text("挖矿");
                $("#result5").text("随机数为："+response["nonce"]);
            }
        });
    });
    $("#newBlock").click(() => {
        let nonce = $("#nonce").val();
        let miner = $("#receive").val();
        if (nonce === "") {
            alert("随机数不能为空！");
        } else if (miner === "") {
            alert("奖励地址不能为空！");
        } else if (!Number.isInteger(Number(nonce))) {
            alert("随机数必须为整数！");
        } else if (!isStringOnlyDigits(nonce)) {
            alert("随机数为无效的输入！");
        } else {
            $.ajax({
                url: "http://localhost:8080/newBlock",
                method: "Post",
                data: {
                    "nonce": nonce,
                    "miner": miner,
                },
                success: function (response) {
                    if (response["state"] === "0") {
                        alert("交易池还没有任何交易！");
                    } else if (response["state"] === "1") {
                        alert("奖励地址不存在！");
                    } else if (response["state"] === "2") {
                        alert("随机数无效！");
                    } else if (response["state"] === "3") {
                        alert("交易被篡改！");
                    } else if (response["state"] === "4") {
                        alert("成功提交一个区块！");
                        $("#nonce").val("");
                        $("#receive").val("");
                    }
                }
            });
        }
    })
    function isStringOnlyDigits(input) {
        return /^\d+$/.test(input);
    }

    $("#sendTx").click(() => {
        let from = $("#from").val();
        let to = $("#to").val();
        let amount = $("#amount").val();
        let password = $("#pswd").val();
        if (from === "") {
            alert("发送地址不能为空！");
        } else if (to === "") {
            alert("接收地址不能为空！");
        } else if (amount === "") {
            alert("交易数额不能为空！");
        } else if (!Number.isInteger(Number(amount))) {
            alert("交易数额必须为整数！");
        } else if (parseInt(amount) <= 0) {
            alert("交易数额不能小于等于0！");
        } else if (password === "") {
            alert("密码不能为空！");
        } else {
            $.ajax({
                url: "http://localhost:8080/newTx",
                method: "Post",
                data: {
                    "from": from,
                    "to": to,
                    "amount": amount,
                    "password": password,
                },
                success: function (response) {
                    if (response["state"] === "0") {
                        alert("发送方地址不存在！");
                    } else if (response["state"] === "1") {
                        alert("接收方地址不存在！");
                    } else if (response["state"] === "2") {
                        alert("密码错误！");
                    } else if (response["state"] === "3") {
                        alert("交易池已满！");
                    } else if (response["state"] === "4") {
                        alert("余额不足！");
                    } else if (response["state"] === "5") {
                        alert("交易成功！");
                        $("#from").val("");
                        $("#to").val("");
                        $("#amount").val("");
                        $("#pswd").val("");
                    }
                },
            });
        }
    })

    $("#searchButton").click(() => {
        let addr = $("#getBalanceAccount").val();
        $.ajax({
            url: "http://localhost:8080/getBalance",
            method: "Post",
            data: {
                "addr": addr,
            },
            success: function (response) {
                if (response["balance"] === "unExist") {
                    $("#result2").text("账户不存在，查询失败！");
                } else {
                    $("#result2").text(response["balance"]);
                }
            },
        })
    })

    $("#newAccountButton").click(() => {
        let passwd = $("#password").val();
        $.ajax({
            url: "http://localhost:8080/newAccount2",
            method: "Post",
            data: {
                "password": passwd,
            },
            success: function (response) {
                $("#result3").text(response["addr"]);
            },
        })
    })

    $("#home").click(() => {
        showDiv('div1')
    });
    $("#getBalance").click(() => {
        showDiv('div2')
    });
    $("#newAccount").click(() => {
        showDiv('div3')
    });
    $("#newTx").click(() => {
        showDiv('div4')
    });
    $("#mine").click(() => {
        showDiv('div5')
    })
    function showDiv(divId) {
        let divs = document.getElementsByClassName('div-container');
        for (let i = 0; i < divs.length; i++) {
            divs[i].classList.remove('active');
        }
        let div = document.getElementById(divId);
        div.classList.add('active');
    }

    // 保存数据
    saveButton.addEventListener('click', () => {
        let beforeValue;
        let key = keyInput.value;
        chrome.runtime.sendMessage({action: 'getData', key}, (response) => {
            beforeValue = response.data;
        });
        setTimeout(() => {
            if (beforeValue === undefined) {
                let result = confirm("该键还没有绑定值，执行此操作将会与该值进行绑定，确定要执行此操作吗？");
                if (result) {
                    let key = keyInput.value;
                    let value = valueInput.value;
                    chrome.runtime.sendMessage({action: 'saveData', key, value}, (response) => {
                        resultDiv.innerText = response.message;
                    });
                }
            } else {
                let result = confirm("该键已经绑定了值，执行此操作将会覆盖之前存储的值，确定要执行此操作吗？");
                if (result) {
                    let key = keyInput.value;
                    let value = valueInput.value;
                    chrome.runtime.sendMessage({action: 'saveData', key, value}, (response) => {
                        resultDiv.innerText = response.message;
                    });
                }
            }
        }, 300);
    });

    // 获取数据
    getButton.addEventListener('click', () => {
        let key = keyInput.value;
        chrome.runtime.sendMessage({action: 'getData', key}, (response) => {
            if (response.data === undefined) {
                resultDiv.innerText = "该键还没有与值进行数据存储哦。";
            } else {
                resultDiv.innerText = response.data;
            }
        });
    });

    // 复制数据
    copyButton.addEventListener('click', () => {
        let copyContent = resultDiv.innerText

        // 创建一个临时的 textarea 元素
        let tempTextarea = document.createElement("textarea");
        tempTextarea.value = copyContent;
        document.body.appendChild(tempTextarea);

        // 选中文本并复制
        tempTextarea.select();
        document.execCommand("copy");
        // 删除临时 textarea 元素
        document.body.removeChild(tempTextarea);
    });
});
