from cefpython3 import cefpython as cef
import platform
import sys
from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *


def main():
    check_version()
    sys.excepthook = cef.ExceptHook         # 异常处理
    cef_settings = {}
    cef.Initialize(cef_settings)            # 设置cef
    app = CefApplication(sys.argv)

    main_window = MainWindow()
    main_window.show()

    app.exec_()
    app.stopTimer()
    del main_window
    del app
    cef.Shutdown()


def check_version():
    ver = cef.GetVersion()
    print("the version {ver}".format(ver=ver))
    print("platform: {} {}".format(platform.python_version(), platform.architecture()))

class CefApplication(QApplication):
    def __init__(self, args):
        super(CefApplication, self).__init__(args)
        self.timer = self.createTimer()

    def createTimer(self):
        timer = QTimer()
        # noinspection PyUnresolvedReferences
        timer.timeout.connect(self.onTimer)
        timer.start(10)
        return timer

    def onTimer(self):
        cef.MessageLoopWork()           # 启动cef消息循环

    def stopTimer(self):
        # Stop the timer after Qt's message loop has ended
        self.timer.stop()

class MainWindow(QMainWindow):
    def __init__(self):
        super(MainWindow, self).__init__(None)
        self.resize(1200, 800)

        # Layout
        center_frame = QFrame()
        layout = QVBoxLayout(center_frame)
        layout.setSpacing(0)
        layout.setContentsMargins(0,0,0,0)
        self.setCentralWidget(center_frame)
        cef_widget = CefWidget()
        layout.addWidget(cef_widget)

class CefWidget(QWidget):
    def __init__(self, parent=None):
        super(CefWidget, self).__init__(parent)
        self.parent = parent
        self.hidden_window = QWindow()
        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0,0,0,0)

        cef_window_info = cef.WindowInfo()
        rect = [0,0, self.width(), self.height()]
        cef_window_info.SetAsChild(self.hidden_window.winId(), rect)
        self.browser = cef.CreateBrowserSync(cef_window_info, 
                                url="https://www.baidu.com/")
        container = QWidget.createWindowContainer(self.hidden_window)
        self.layout.addWidget(container)    # linux下用container包裹browser

    def moveEvent(self, _):                 # move事件通知browser
        self.x = 0
        self.y = 0
        if self.browser:
            # self.browser.SetBounds(self.x, self.y,self.width(), self.height())
            self.browser.NotifyMoveOrResizeStarted()

    def resizeEvent(self, event):           # resize事件通知browser
        size = event.size()
        if self.browser:
            # self.browser.SetBounds(self.x, self.y, size.width(), size.height())
            self.browser.NotifyMoveOrResizeStarted()

    def closeEvent(self, event):
        if self.browser:
            self.browser.CloseBrowser(True)
            self.browser = None

if __name__ == "__main__":
    main()