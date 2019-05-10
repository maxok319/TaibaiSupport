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
from TaibaiUtils import *
from TaibaiConfig import *
import json
import time

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
                    w.userId = userId
                    w.layout.addWidget(QLabel(str(userId)))
                    self.participantWidgetMap[userId] = w
                    self.participantWidgetMap[userId].show()
                    w.positionChanged.connect(self.videoWidgetPositionChanged)
                self.participantWidgetMap[userId].serverRect = participant["rect"]
                self.participantWidgetMap[userId].setMoveableArea(self.courseArea)
                self.participantWidgetMap[userId].setGeometry(self.mapFromServer(participant["rect"]))

    def resizeEvent(self, event):
        self.courseArea = StandardAreaInRect(event.size(), taibai_config["standard_coursearea_size"])
        for participantWidget in self.participantWidgetMap.values():
            participantWidget.setMoveableArea(self.courseArea)
            participantWidget.setGeometry(self.mapFromServer(participantWidget.serverRect))

    # 把本地的换成server端的
    def mapToServer(self, localRect):
        localX = localRect.x()
        localY = localRect.y()
        localWidth = localRect.width()
        localHeight = localRect.height()

        courseAreaX = self.courseArea.x()
        courseAreaY = self.courseArea.y()
        courseAreaWidth = self.courseArea.width()
        courseAreaHeight = self.courseArea.height()

        serverRectX = (localX - courseAreaX) / courseAreaWidth * taibai_config["standard_coursearea_width"]
        serverRectY = (localY - courseAreaY) / courseAreaHeight * taibai_config["standard_coursearea_height"]  
        serverRectWidth = localWidth / courseAreaWidth * taibai_config["standard_coursearea_width"]
        serverRectHeight = localHeight / courseAreaHeight * taibai_config["standard_coursearea_height"]

        serverRect = {}
        serverRect["X"] = int(serverRectX)
        serverRect["Y"] = int(serverRectY)
        serverRect["Width"] = int(serverRectWidth)
        serverRect["Height"] = int(serverRectHeight)

        return serverRect

    # 把server端的换成本地的
    def mapFromServer(self, serverRect):
        serverRectX = serverRect["X"]
        serverRectY = serverRect["Y"]
        serverRectWidth = serverRect["Width"]
        serverRectHeight = serverRect["Height"]

        courseAreaX = self.courseArea.x()
        courseAreaY = self.courseArea.y()
        courseAreaWidth = self.courseArea.width()
        courseAreaHeight = self.courseArea.height()

        localX = courseAreaX + serverRectX * courseAreaWidth / taibai_config["standard_coursearea_width"]
        localY = courseAreaY + serverRectY * courseAreaHeight / taibai_config["standard_coursearea_height"]
        localWidth = serverRectWidth * courseAreaWidth / taibai_config["standard_coursearea_width"]
        localHeight = serverRectHeight * courseAreaHeight / taibai_config["standard_coursearea_height"]

        localRect = QRect(localX, localY, localWidth, localHeight)
        return localRect

    def videoWidgetPositionChanged(self, userId):
        for participantWidget in self.participantWidgetMap.values():
            if participantWidget.userId == userId:
                serverRect = self.mapToServer(participantWidget.geometry())

                event = {
                    "eventTime" : int(time.time()),
                    "eventType" : "videoPositionChanged",
                    "eventProducer": 0,
                    "eventContent" : {
                        "userId" : userId,
                        "rect" : serverRect
                    }
                }
                package = json.dumps(event)
                self.client.sendTextMessage(package)


