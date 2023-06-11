package modules

import (
	"github.com/simulatedsimian/joystick"
	"simUI/code/utils"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var JOYSTICK int8
var user32 = syscall.NewLazyDLL("user32.dll")
var procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
var getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")

func CheckJoystick() (status int8) {

	if JOYSTICK == 1 {
		//已存在
		return -1
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		jsid := 0

		js, jserr := joystick.Open(jsid)

		if jserr != nil {
			JOYSTICK = 0
			wg.Done()
			return
		}
		btnLock := map[string]int64{}
		var dirLock int64
		JOYSTICK = 1
		wg.Done()
	EXIT:
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(40)):

				/*
					// 判断句柄是否属于当前进程
					hwnd, _, _ := procGetForegroundWindow.Call()
					processID := getProcessIDFromHWND(hwnd)
					if processID != uint32(syscall.Getpid()) {
						break
					}
				*/
				active := utils.CheckWinActive()
				if active == false {
					break
				}

				jinfo, err := js.Read()
				if err != nil {
					break EXIT
				}

				btn := GetJoystickButtons(jinfo.Buttons)
				dir := GetJoystickDirection(jinfo.AxisData)
				current := time.Now().UnixNano() / 1e6

				if dir > 0 {
					if current-dirLock < 400 {
						break
					}
					dirLock = current
					utils.ViewDirection(dir)
					//fmt.Println("Buttons:", dir)
				}

				if btn != "" {
					if _, ok := btnLock[btn]; ok {
						if current-btnLock[btn] < 500 {
							break
						}
					}

					btnLock[btn] = current
					utils.ViewButton(btn)
					//fmt.Println("AxisData:", btn, current-btnLock[btn])
				}
			}
		}

	}()

	wg.Wait()
	return JOYSTICK
}

//读取方向
//1上2下3左4右
func GetJoystickDirection(axis []int) int {
	if len(axis) == 6 {
		//条件位置不能换
		if axis[0] == -32767 || axis[4] == -32767 {
			return 3
		} else if axis[0] == 32768 || axis[4] == 32768 {
			return 4
		} else if axis[1] == -32767 || axis[5] == -32767 {
			return 1
		} else if axis[1] == 32768 || axis[5] == 32768 {
			return 2
		}

	} else {
		if len(axis) >= 1 && axis[0] == 32768 {
			return 4
		} else if len(axis) >= 1 && axis[0] == -32767 {
			return 3
		} else if len(axis) > 1 && axis[1] == 32768 {
			return 2
		} else if len(axis) > 1 && axis[1] == -32767 {
			return 1
		} else if len(axis) > 3 && axis[3] == 32768 {
			return 2
		} else if len(axis) > 3 && axis[3] == -32767 {
			return 1
		} else if len(axis) > 4 && axis[4] == 32768 {
			return 4
		} else if len(axis) > 4 && axis[4] == -32767 {
			return 3
		} else if len(axis) > 5 && axis[5] == 32768 {
			return 4
		} else if len(axis) > 5 && axis[5] == -32767 {
			return 3
		} else if len(axis) > 6 && axis[6] == 32768 {
			return 2
		} else if len(axis) > 6 && axis[6] == -32767 {
			return 1
		}
	}

	return 0
}

//读取按钮
func GetJoystickButtons(button uint32) string {
	btn := ""
	switch button {
	case 1:
		btn = "A" //A
		break
	case 2:
		btn = "B" //B
		break
	case 4:
		btn = "X" //X
		break
	case 8:
		btn = "Y" //Y
		break
	case 16:
		btn = "LB" //LB
		break
	case 32:
		btn = "RB" //RB
		break
	case 64:
		btn = "BACK" //back
		break
	case 128:
		btn = "START" //start
		break
	case 192:
		btn = "START+BACK" //start + back
		break
	}
	return btn
}

// 获取 HWND 所属的进程 ID
func getProcessIDFromHWND(hwnd uintptr) uint32 {
	var processID uint32
	threadID, _, _ := getWindowThreadProcessId.Call(
		hwnd,
		uintptr(unsafe.Pointer(&processID)),
	)
	if threadID == 0 {
		return 0
	}
	return processID
}
