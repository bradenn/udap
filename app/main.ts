const {app, BrowserWindow, Menu} = require("electron");
const path = require("path");

let win;

function createWindow() {
    win = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            preload: path.join(__dirname, "preload.js"),
        },
        showCursor: false,
        title: "Udap-Nexus Terminal Endpoint",
        backgroundColor: '#000000',
        autoHideMenuBar: true,
        show: true,
        kiosk: false,
        fullscreen: true,
        removeMenu: true,
        frame: false,
    });


    win.webContents.setZoomFactor(1)
    win.loadURL('http://localhost:5002')

}


app.whenReady().then(() => {
    createWindow();
    app.on("activate", () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow();
        } else {
            Menu.setApplicationMenu(null);
            win.setMenu(null)
        }
    });
});

app.on("window-all-closed", () => {
    if (process.platform !== "darwin") {
        app.quit();
    }
});
