const {app, BrowserWindow} = require("electron");
const path = require("path");


function createWindow() {
    const win = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            preload: path.join(__dirname, "preload.js"),
        },
        showCursor: false,
        title: "UDAP Endpoint Terminal @ CONFIDENTIAL COPY",
        acceptFirstMouse: true,
        backgroundColor: '#000000',
        autoHideMenuBar: true,
        fullscreen: true,
    });

    win.webContents.setZoomFactor(1)
    win.loadURL('http://localhost:5002')
}

app.whenReady().then(() => {
    createWindow();

    app.on("activate", () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow();
        }
    });
});

app.on("window-all-closed", () => {
    if (process.platform !== "darwin") {
        app.quit();
    }
});
