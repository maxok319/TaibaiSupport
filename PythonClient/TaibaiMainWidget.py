from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
import platform
import ctypes
import sys
from TaibaiClassWidget import TaibaiClassWidget

class TaibaiMainWidget(QWidget):
    def __init__(self):
        super(TaibaiMainWidget, self).__init__()
        self.layout = QHBoxLayout(self)
        self.layout.setSpacing(0)
        self.layout.setContentsMargins(0, 0, 0, 0)
        self.classwidget = TaibaiClassWidget(self)
        self.layout.addWidget(self.classwidget)

