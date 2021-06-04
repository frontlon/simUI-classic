package modules

import (
	"github.com/simulatedsimian/joystick"
	"os"
	"os/exec"
	"runtime"
	"simUI/code/utils"
	"sync"
	"time"
)

var JOYSTICK int8

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
		var btnLock [5]int64
		var dirLock int64
		JOYSTICK = 1
		wg.Done()
	EXIT:
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(40)):

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

				if btn > 0 {

					if btn == 1 || btn == 2 {
						if current-btnLock[btn] < 500 {
							break
						}

					} else if btn == 3 || btn == 4 {
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

	if axis[1] == -32767 {
		return 1
	} else if axis[1] == 32768 {
		return 2
	} else if axis[0] == -32767 {
		return 3
	} else if axis[0] == 32768 {
		return 4
	} else if len(axis) >= 7 && axis[6] == -32767 {
		return 1
	} else if len(axis) >= 7 && axis[6] == 32768 {
		return 2
	} else if len(axis) >= 7 && axis[5] == -32767 {
		return 3
	} else if len(axis) >= 7 && axis[5] == 32768 {
		return 4
	}
	return 0
}

//读取按钮
func GetJoystickButtons(button uint32) int {
	btn := 0
	switch button {
	case 1:
		btn = 1 //A
		break
	case 2:
		btn = 2 //B
		break
	case 4:
		btn = 3 //X
		break
	case 8:
		btn = 4 //Y
		break
	case 16:
		btn = 5 //LB
		break
	case 32:
		btn = 6 //RB
		break
	default:
		btn = 0
	}
	return btn
}

/**
 * 关闭simui
 **/
func killSoft() error {

	switch runtime.GOOS {
	case "darwin":
		c := exec.Command("kill", utils.ToString(os.Getpid()))
		c.Start()
	case "windows":
		c := exec.Command("taskkill.exe", "/T", "/PID", utils.ToString(os.Getpid()))
		c.Start()
	case "linux":
		c := exec.Command("kill", utils.ToString(os.Getpid()))
		c.Start()
	}
	return nil
}