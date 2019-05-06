import sys
from QCefApplication import QCefApplication
from QCefWidget import QCefWidget

def main():
    app = QCefApplication(sys.argv)
    widget = QCefWidget()
    widget.resize(800, 600)
    widget.show()
    app.exec_()
    
if __name__ == '__main__':
    main()