package modules

import (
	"fmt"
	"github.com/simulatedsimian/joystick"
	"simUI/code/utils"
	"time"
)

func CheckJoystick() {
	time.Sleep(2 * time.Second)

	go func() {

		jsid := 0

		js, jserr := joystick.Open(jsid)

		if jserr != nil {
			fmt.Println(jserr)
			return
		}
		var btnLock [5]int64
		var dirLock int64
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
						if current-btnLock[btn] < 1000 {
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
}


//读取方向
func GetJoystickDirection(axis []int) int {
	if axis[1] == -32767 {
		return 1
	} else if axis[1] == 32768 {
		return 2
	} else if axis[0] == -32767 {
		return 3
	} else if axis[0] == 32768 {
		return 4
	}
	return 0
}

//读取按钮
func GetJoystickButtons(button uint32) int {
	btn := 0
	switch button {
	case 1:
		btn = 1
		break
	case 2:
		btn = 2
		break
	case 4:
		btn = 3
		break
	case 8:
		btn = 4
		break
	default:
		btn = 0
	}
	return btn
}
