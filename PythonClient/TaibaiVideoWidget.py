from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
from PyQt5.QtWebSockets import *
import platform
import ctypes
import sys
from QCefWidget import QCefWidget
from TaibaiWebsocket import TaibaiWebsocket
import json
from TaibaiMoveableSubWidget import TaibaiMoveableSubWidget
from TaibaiConfig import *

class TaibaiVideoWidget(TaibaiMoveableSubWidget):
    def __init__(self, parent):
        super(TaibaiVideoWidget, self).__init__(parent)
        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)
    
    # 按照server返回的位置 和 现在客户端的大小 调整自己
    def adjustPosition(self):
        
        rectX = self.serverRect["X"]
        rectY = self.serverRect["Y"]
        rectWidth = self.serverRect["Width"]
        rectHeight = self.serverRect["Height"]

        areaX = self.allowarea.x()
        areaY = self.allowarea.y()
        areaWidth = self.allowarea.width()
        areaHeight = self.allowarea.height()

        currX = rectX * areaWidth / taibai_config["standard_coursearea_width"]
        currY = rectY * areaHeight / taibai_config["standard_coursearea_height"]
        currWidth = rectWidth * areaWidth / taibai_config["standard_coursearea_width"]
        currheight = rectHeight * areaHeight / taibai_config["standard_coursearea_height"]

        self.resize(currWidth, currheight)
        self.move(areaX + currX, areaY + currY)
    
    # 按照客户端的大小 以及自己当前的位置和大小 给server发送位置
    def mousePressEvent(self, event):
        super(TaibaiVideoWidget, self).mousePressEvent(event)
        if event.button() == Qt.LeftButton:
            self.aaa = self.pos()

    def mouseReleaseEvent(self, event):
        super(TaibaiVideoWidget, self).mouseReleaseEvent(event)
        if event.button() == Qt.LeftButton:
            self.bbb = self.pos
            if self.aaa != self.bbb:
                print("fa song")

    


        

    


        