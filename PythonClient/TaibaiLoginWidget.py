from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
from PyQt5.QtWebSockets import *
import platform
import ctypes
import sys
from QCefWidget import QCefWidget
from cefpython3 import cefpython as cef

class TaibaiLoginWidget(QWidget):
    joinClicked = pyqtSignal(str, name="joinClicked")

    def __init__(self, parent):
        super(TaibaiLoginWidget, self).__init__(parent)
        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)
        demofile = "file://" + sys.path[0] + "/login.html"
        self.cefwidget = QCefWidget(demofile)
        self.layout.addWidget(self.cefwidget)

        self.cefwidget.bindCefObject("LoginWidget", self)

    def joinClickedInH5(self, jsonparam=None):
        self.joinClicked.emit(jsonparam)


    


        