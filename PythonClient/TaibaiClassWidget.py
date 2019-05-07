from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
from PyQt5.QtWebSockets import *
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
        demofile = "file://" + sys.path[0] + "/classroom.html"
        self.cefwidget = QCefWidget(demofile)
        self.layout.addWidget(self.cefwidget)

        self.client = QWebSocket("",QWebSocketProtocol.Version13,None)
        self.client.connected.connect(self.connected)
        self.client.error.connect(self.error)
        self.client.disconnected.connect(self.disconnected)
        self.client.textMessageReceived.connect(self.textMessageReceived)
    
    def startWS(self, wsurl):
        self.client.open(QUrl(wsurl))

    def connected(self):
        print("connected")
    
    def disconnected(self):
        print("disconnected")

    def error(self, error_code):
        print(self.client.errorString())
    
    def textMessageReceived(self, message):
        print(message)

    


        