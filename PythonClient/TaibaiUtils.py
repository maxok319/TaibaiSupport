from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *

def IsRectInRect_Point(subTopLeft, subBottomRight, parentTopLeft, parentBottomRight):
    top_left_ok = (subTopLeft.x() >= parentTopLeft.x()) and (subTopLeft.y() >= parentTopLeft.y())
    bottom_right_ok = (subBottomRight.x() <= parentBottomRight.x()) and (subBottomRight.y() <= parentBottomRight.y())
    return top_left_ok and bottom_right_ok

def IsRectInRect_Rect(subRect, parentRect):
    return IsRectInRect_Point(subRect.topLeft(), subRect.bottomRight(), parentRect.topLeft(), parentRect.bottomRight())
