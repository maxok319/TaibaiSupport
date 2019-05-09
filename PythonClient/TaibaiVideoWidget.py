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
class TaibaiVideoWidget(TaibaiMoveableSubWidget):
    def __init__(self, parent):
        super(TaibaiVideoWidget, self).__init__(parent)
        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)
    


        

    


        