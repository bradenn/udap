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
        title: "UDAP Endpoint",
        backgroundColor: '#000000',
        autoHideMenuBar: true,
        show: true,
        kiosk: false,
        fullscreen: true,
        removeMenu: true,
        frame: false,
    });


    // attachTouchMode()
    // win.webContents.debugger.on('detach', (event, reason) => {
    //     console.log('Debugger detached due to: ', reason)
    // })
    // attachTouchMode()
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
    win.webContents.setZoomFactor(1)
    win.loadURL('http://localhost:5002')

}

function attachTouchMode() {

    try {
        // works with 1.1 too
        win.webContents.debugger.attach('1.2')
    } catch (err) {
        console.log('Debugger attach failed: ', err)
    }

    const isDebuggerAttached = win.webContents.debugger.isAttached()
    console.log('debugger attached? ', isDebuggerAttached)

    win.webContents.debugger.on('detach', (event, reason) => {
        console.log('Debugger detached due to: ', reason)
    });

    win.webContents.debugger.sendCommand('Emulation.setEmitTouchEventsForMouse', {enabled: true});
    win.webContents.debugger.sendCommand('Emulation.setTouchEmulationEnabled', {
        enabled: true,
        configuration: 'desktop',
    });
}

app.on('certificate-error', (event, webContents, url, error, certificate, callback) => {
    // On certificate error we disable default behaviour (stop loading the page)
    // and we then say "it is all fine - true" to the callback
    event.preventDefault();
    callback(true);
});

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
