.PHONY: test-lib
# It should show something similar:
# -I/usr/include/qt6/QtCore -I/usr/include/qt6 -DQT_CORE_LIB -I/usr/lib64/qt6/mkspecs/linux-g++ -lQt6Core
# If pkg-config not found then need install
test-lib:
	pkg-config --cflags --libs Qt6Core
	pkg-config --modversion Qt6Core
	
test:
	go test ./... -v
