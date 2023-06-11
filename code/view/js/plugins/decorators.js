function @when(func, evtType)
{
  function t(evt)
  {
    var r = false;
    if( evt.type == evtType) 
      r = func.call(this,evt);
    if(t.next) return t.next.call(this,evt) || r;
    return r;  
  }
  var principal = this instanceof Behavior? this : self;
  t.next = principal.onControlEvent; 
  return principal.onControlEvent = t;
}
function @click(func)
{
  function t(evt)
  {
    var r = false;
    if(  evt.type == Event.BUTTON_CLICK
      || evt.type == Event.HYPERLINK_CLICK
      || evt.type == Event.MENU_ITEM_CLICK ) 
      r = func.call(this,evt);
    if(t.next) return t.next.call(this,evt) || r;
    return r;  
  }
  var principal = this instanceof Behavior? this : self;
  t.next = principal.onControlEvent; 
  return principal.onControlEvent = t;
}
function @change(func)
{
  function t(evt)
  {
    var r = false;
    if(  evt.type == Event.BUTTON_STATE_CHANGED 
      || evt.type == Event.SELECT_SELECTION_CHANGED 
      || evt.type == Event.EDIT_VALUE_CHANGED ) 
      r = func.call(this,evt);
    if(t.next) return t.next.call(this,evt) || r;
    return r;  
  }
  var principal = this instanceof Behavior? this : self;
  t.next = principal.onControlEvent; 
  return principal.onControlEvent = t;
}

function @on(func, selector)
{
  return function(evt)
  {
    if( evt.target.match(selector) )
      return func.call(this,evt);
  }
}
function @key(func, keyCode = undefined, modifiers..)
{
  function t(evt)
  {
    var r = false;
    if( evt.type == Event.KEY_DOWN && 
        (keyCode === undefined || (keyCode == evt.keyCode)) ) 
          r = func.call(this,evt);
    if(t.next) return t.next.call(this,evt) || r;
    return r;  
  }
  var principal = this instanceof Behavior ? this : self;
  t.next = principal.onKey; 
  principal.onKey = t;

}
function @CTRL(func) { return function(evt) { if( evt.ctrlKey === true ) return func.call(this,evt); } }
function @NOCTRL(func) { return function(evt) { if( evt.ctrlKey === false ) return func.call(this,evt); } }
function @SHIFT(func) { return function(evt) { if( evt.shiftKey === true ) return func.call(this,evt); } }
function @NOSHIFT(func) { return function(evt) { if( evt.shiftKey === false ) return func.call(this,evt); } }
function @ALT(func) { return function(evt) { if( evt.altKey === true ) return func.call(this,evt); } }
function @NOALT(func) { return function(evt) { if( evt.altKey === false ) return func.call(this,evt); } }
function @mouse(func, mouseEvtType)
{
  function t(evt)
  {
    var r = false;
    if( evt.type == mouseEvtType) 
      r = func.call(this,evt);
    if(t.next) return t.next.call(this,evt) || r;
    return r;  
  }
  var principal = this instanceof Behavior? this : self;
  t.next = principal.onMouse; 
  return principal.onMouse = t;
}