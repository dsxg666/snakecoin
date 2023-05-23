chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === 'saveData') {
        // 存储数据
        chrome.storage.local.set({ [request.key]: request.value }, () => {
            sendResponse({ message: '数据存储成功！' });
        });
        return true;
    } else if (request.action === 'getData') {
        // 获取数据
        chrome.storage.local.get(request.key, (result) => {
            if (result[request.key] === undefined) {
                sendResponse({ data: "该键还没有与值进行数据存储哦。" });
            } else {
                sendResponse({ data: result[request.key] });
            }
        });
        return true;
    }
});
