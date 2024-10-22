﻿function checkWinActive() {
   return view.windowIsActive;
}

function joystickDirection(dir) {

   var switch_romlist = $(#switch_romlist).attributes["value"];

   if (dir == 1) { //上
      scrollUp(switch_romlist); //div滚动
      if (switch_romlist == "1") { //模块
         keyRomUp();
      } else if (switch_romlist == "2") { //列表
         keyRomLeft();
      }
   } else if (dir == 2) { //下
      scrollDown(switch_romlist); //div滚动
      if (switch_romlist == "1") { //模块
         keyRomDown();
      } else if (switch_romlist == "2") { //列表
         keyRomRight();
      }
   } else if (dir == 3) { //左
      if (switch_romlist == "1") { //模块
         scrollUp(switch_romlist); //div滚动
         keyRomLeft();
      }
   } else if (dir == 4) { //右
      if (switch_romlist == "1") { //模块
         scrollDown(switch_romlist); //div滚动
         keyRomRight();
      }
   }
}

//手柄按键
function joystickButton(btn) {
   var current;
   switch (btn) {
      case "A":
         current = $(#buttons li:current);
         if (current == undefined) {
            current = $(#romlist li:current);
            if (current != undefined) {
               romListRunGame(current);
            }
         } else {
            sidebarRunGame(current);
         }
         break;
      case "B":
         nextMenu();
         break;
      case "Y":
         nextPlatform();
         break;
      case "LB":
         nextSim();
         break;
      case "RB":
         nextSubGame();
         break;
   }

}

//检测外设
function checkJoystick() {

   var status = view.checkJoystick();
   switch (status) {
      case 1:
         alert("设备连接成功");
         break;
      case 0:
         alert("没有找到设备或设备连接失败");
         break;
      case -1:
         alert("设备已连接");
         break;
   }


}

//变更手柄焦点框
function joystickFocus(obj) {
   info(obj.attributes['rid']);
   var(x, y, w, h) = obj.box(#rectw, #margin, #parent);
   moveJoystickFocus(x, y, w, h);
}

//移动焦点框
function moveJoystickFocus(x, y, w, h) {
   $(#joystickFocus).style["width"] = dip(w);
   $(#joystickFocus).style["height"] = dip(h);
   $(#joystickFocus).style["left"] = dip(x + 10);
   $(#joystickFocus).style["top"] = dip(y);
   $(#joystickFocus).style["display"] = "block";
}