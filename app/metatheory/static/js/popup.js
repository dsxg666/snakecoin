document.addEventListener('DOMContentLoaded', () => {
    const keyInput = document.getElementById('keyInput');
    const valueInput = document.getElementById('valueInput');
    const saveButton = document.getElementById('saveButton');
    const getButton = document.getElementById('getButton');
    const resultDiv = document.getElementById('result');
    const copyButton = document.getElementById('copyButton');


    // 保存数据
    saveButton.addEventListener('click', () => {
        let beforeValue;
        let key = keyInput.value;
        chrome.runtime.sendMessage({ action: 'getData', key }, (response) => {
            beforeValue = response.data;
        });
        setTimeout(() => {
            if (beforeValue === undefined) {
                let result = confirm("该键还没有绑定值，执行此操作将会与该值进行绑定，确定要执行此操作吗？");
                if (result) {
                    const key = keyInput.value;
                    const value = valueInput.value;
                    chrome.runtime.sendMessage({ action: 'saveData', key, value }, (response) => {
                        resultDiv.innerText = response.message;
                    });
                }
            } else {
                let result = confirm("该键已经绑定了值，执行此操作将会覆盖之前存储的值，确定要执行此操作吗？");
                if (result) {
                    const key = keyInput.value;
                    const value = valueInput.value;
                    chrome.runtime.sendMessage({ action: 'saveData', key, value }, (response) => {
                        resultDiv.innerText = response.message;
                    });
                }
            }
        }, 300);
    });

    // 获取数据
    getButton.addEventListener('click', () => {
        const key = keyInput.value;
        chrome.runtime.sendMessage({ action: 'getData', key }, (response) => {
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
