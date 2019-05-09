from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
from TaibaiUtils import *

class TaibaiMoveableSubWidget(QWidget):
    def __init__(self, parent):
        super(TaibaiMoveableSubWidget, self).__init__(parent)
        
    def mousePressEvent(self, event):
        if event.button() == Qt.LeftButton:
            self.origin_position = self.pos()
            self.dragPosition = event.globalPos()
            event.accept()

    def mouseMoveEvent(self, event):
        if event.buttons() & Qt.LeftButton:
            delta = event.globalPos() - self.dragPosition

            target_topleft = self.origin_position + delta
            target_bottomRight = target_topleft + QPoint(self.size().width(), self.size().height())
            parent_topleft = QPoint(0, 0)
            parent_bottomRight = QPoint(self.parent().size().width(), self.parent().size().height())

            subRect = QRect(target_topleft, target_bottomRight)
            parentRect = QRect(parent_topleft, parent_bottomRight)

            if not IsRectInRect_Rect(subRect, parentRect):
                if target_topleft.x() < 0:
                    target_topleft.setX(0)
                if target_topleft.y() < 0:
                    target_topleft.setY(0)
                if target_bottomRight.x() > parent_bottomRight.x():
                    target_topleft.setX(parent_bottomRight.x() - self.size().width())
                if target_bottomRight.y() > parent_bottomRight.y():
                    target_topleft.setY(parent_bottomRight.y() - self.size().height())
                self.move(target_topleft)
                event.accept()

            self.move(target_topleft)
            event.accept()
