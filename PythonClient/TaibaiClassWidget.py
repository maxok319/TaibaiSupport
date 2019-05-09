from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
from PyQt5.QtWebSockets import *
import platform
import ctypes
import sys
from QCefWidget import QCefWidget
from TaibaiWebsocket import TaibaiWebsocket
from TaibaiVideoWidget import TaibaiVideoWidget
import json

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

        self.participantWidgetMap={}
    
    def startWS(self, wsurl):
        self.client.open(QUrl(wsurl))

    def connected(self):
        print("connected")
    
    def disconnected(self):
        print("disconnected")

    def error(self, error_code):
        print(self.client.errorString())
    
    def textMessageReceived(self, wspackage):
        print(wspackage)
        messageObject = json.loads(wspackage)
        messageType = messageObject["messageType"]
        messageTime = messageObject["messageTime"]
        messageContent = messageObject["messageContent"]

        if messageType=="classroomStatus":
            
            for participant in messageContent["participantList"]:
                userId = participant["userId"]
                if userId not in self.participantWidgetMap:
                    w = TaibaiVideoWidget(self)
                    w.layout.addWidget(QLabel(str(userId)))
                    w.setMoveableArea(QRect(200, 100, 600, 500))
                    self.participantWidgetMap[userId] = w
                    self.participantWidgetMap[userId].show()
                rect = participant["rect"]
                self.participantWidgetMap[userId].resize(rect["Width"], rect["Height"])
                self.participantWidgetMap[userId].move(rect["X"], rect["Y"])
    
    

        

    


        