class FileDropZone : Element 
{
    this var filter;  
    this var files = []; // filtered files

    function checkFiles(list) {
    if(typeof list != #array )
        list = [list];
    const patterns = this.filter;
    function flt(fn) {
        for(var x in patterns)
        if( fn like x ) return true;
        return false;
    }
    this.files = list.filter(flt);
    return this.files.length > 0;
    }

    function attached() {
    this.filter = (this.attributes["accept-drop"] || "*").split(";");
    debug filter: this.filter;
    }

    event dragaccept (evt) {
    if(evt.draggingDataType == #file && this.checkFiles(evt.dragging))
        return true; // accept only files
    return false;
    }

    event dragenter (evt) 
    {
    this.attributes.addClass("active-target");
    return true;
    }  

    event dragleave (evt) 
    {
    this.attributes.removeClass("active-target");
    return true;
    }

    event drag (evt) 
    {
    return true;
    }  

    event drop (evt) 
    {

        var id = ACTIVE_ROM_ID;
        var platform = ACTIVE_PLATFORM;

        if(platform == 0){
            var detailObj = JSON.parse(view.GetGameById(id));
            platform = detailObj.Platform;
        }
   
        this.attributes.removeClass("active-target");
        var opt = this.attributes["opt"];
        var sid = this.attributes["sid"];

        //创建一个模块
        createFileDropZone(platform,opt,sid,this.files.toString());

        return true;
    }       
}