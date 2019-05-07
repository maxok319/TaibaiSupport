from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
import platform
import ctypes
import sys
from TaibaiClassWidget import TaibaiClassWidget
from TaibaiLoginWidget import TaibaiLoginWidget

class TaibaiMainWidget(QWidget):
    def __init__(self):
        super(TaibaiMainWidget, self).__init__()
        self.layout = QStackedLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)
        self.loginwidget = TaibaiLoginWidget(self)
        self.classwidget = TaibaiClassWidget(self)

        self.layout.addWidget(self.loginwidget)
        self.layout.addWidget(self.classwidget)

        self.layout.setCurrentWidget(self.loginwidget)

        self.loginwidget.joinClicked.connect(self.startJoinClass)

    def startJoinClass(self, jsonparams):

        wsurl = "ws://127.0.0.1:8888/ws?classroomId=123&userId=111"
        self.layout.setCurrentWidget(self.classwidget)
        self.classwidget.startWS(wsurl)

