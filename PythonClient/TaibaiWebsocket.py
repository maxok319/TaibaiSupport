from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *
import platform
import ctypes
import sys
import websocket
import time
try:
    import thread
except ImportError:
    import _thread as thread

class TaibaiWebsocket(websocket.WebSocketApp, QObject):
    def __init__(self, wsurl):
        websocket.WebSocketApp.__init__(self, wsurl)
        QObject.__init__(self, None)

    
    signal_on_message = pyqtSignal('QString')

    def _message(self, message):
        pass
              

    def _error(self, error):
        print(error)

    def _close(self):
        print("### closed ###")

    def _open(self):
        print("### opend ###")
        def run(*args):
            for i in range(300):
                time.sleep(1)
                self.send("Hello %d" % i)

            time.sleep(1)
            self.close()
            print("thread terminating...")
        thread.start_new_thread(run, ())

    def start(self):
        # ws://127.0.0.1:8888/ws?classroomId=123&userId=111
        self.on_message = self._message
        self.on_error = self._error
        self.on_close = self._close
        self.on_open = self._open

        thread.start_new_thread(self.run_forever, ())