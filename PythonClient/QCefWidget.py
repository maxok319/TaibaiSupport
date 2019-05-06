from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
import platform
from cefpython3 import cefpython as cef
import ctypes

# Platforms
WINDOWS = (platform.system() == "Windows")
LINUX = (platform.system() == "Linux")
MAC = (platform.system() == "Darwin")

class QCefWidget(QWidget):
    def __init__(self, parent=None):
        # noinspection PyArgumentList
        super(QCefWidget, self).__init__(parent)
        self.resize(800,600)
        self.browser = None
        self.WindowUtils = cef.WindowUtils()

        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)

        window_info = cef.WindowInfo()
        rect = [0, 0, self.width(), self.height()]

        if LINUX:
            self.hidden_window = QWindow()
            window_info.SetAsChild(self.hidden_window.winId(), rect)
            container = QWidget.createWindowContainer(self.hidden_window)
            self.layout.addWidget(container)
        else:
            window_info.SetAsChild(self.winId(), rect)

        self.browser = cef.CreateBrowserSync(window_info, url="https://www.baidu.com")

    def moveEvent(self, _):
        self.x = 0
        self.y = 0
        if self.browser:
            if WINDOWS:
                self.WindowUtils.OnSize(self.winId(), 0, 0, 0)
            elif LINUX:
                self.browser.SetBounds(self.x, self.y,
                                       self.width(), self.height())
            self.browser.NotifyMoveOrResizeStarted()

    def resizeEvent(self, event):
        size = event.size()
        if self.browser:
            if WINDOWS:
                self.WindowUtils.OnSize(self.winId(), 0, 0, 0)
            elif LINUX:
                self.browser.SetBounds(self.x, self.y,
                                       size.width(), size.height())
            self.browser.NotifyMoveOrResizeStarted()
