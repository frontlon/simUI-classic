/**
 * 数据初始化
 */
view.root.on("ready", function(){
    try{
        Init();
    }catch(e){
        view.msgbox(#alert,e);
    }
});



