from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
import platform
import os
import sys
from cefpython3 import cefpython as cef

# Platforms
WINDOWS = (platform.system() == "Windows")
LINUX = (platform.system() == "Linux")
MAC = (platform.system() == "Darwin")

class QCefApplication(QApplication):
    def __init__(self, args):
        super(QCefApplication, self).__init__(args)
        sys.excepthook = cef.ExceptHook  # To shutdown all CEF processes on error
        settings = {}
        if MAC:
            # Issue #442 requires enabling message pump on Mac
            # in Qt example. Calling cef.DoMessageLoopWork in a timer
            # doesn't work anymore.
            settings["external_message_pump"] = True

        cef.Initialize(settings)
        if not cef.GetAppSetting("external_message_pump"):
            self.timer = self.createTimer()
        # self.setupIcon()

    def createTimer(self):
        timer = QTimer()
        # noinspection PyUnresolvedReferences
        timer.timeout.connect(self.onTimer)
        timer.start(10)
        return timer

    def onTimer(self):
        cef.MessageLoopWork()

    def stopTimer(self):
        # Stop the timer after Qt's message loop has ended
        self.timer.stop()

    def setupIcon(self):
        icon_file = os.path.join(os.path.abspath(os.path.dirname(__file__)),
                                 "resources", "{0}.png".format(sys.argv[1]))
        if os.path.exists(icon_file):
            self.setWindowIcon(QIcon(icon_file))