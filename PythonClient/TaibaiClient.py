import sys
from QCefApplication import QCefApplication
from TaibaiMainWidget import TaibaiMainWidget
import websocket

def main():
    websocket.enableTrace(True)
    app = QCefApplication(sys.argv)
    widget = TaibaiMainWidget()
    widget.resize(800, 600)
    widget.show()
    app.exec_()
    
if __name__ == '__main__':
    main()