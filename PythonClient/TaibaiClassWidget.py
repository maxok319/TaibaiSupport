from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
import platform
import ctypes
import sys
from QCefWidget import QCefWidget
from TaibaiWebsocket import TaibaiWebsocket

class TaibaiClassWidget(QWidget):
    def __init__(self, parent):
        super(TaibaiClassWidget, self).__init__(parent)
        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)
        self.cefwidget = QCefWidget()
        self.layout.addWidget(self.cefwidget)

        self.wsurl = "ws://127.0.0.1:8888/ws?classroomId=123&userId=111"
        self.ws = TaibaiWebsocket(self.wsurl)
        self.ws.signal_on_message.connect(self.on_ws_message)

        self.ws.start()
    
    @pyqtSlot('QString')
    def on_ws_message(self, message):
        print(message)

        