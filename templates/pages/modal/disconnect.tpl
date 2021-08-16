
<div id="disconnected" style="display: none">
<div class="modal-dialog">
<div class="modal-content-disconnected">
<div class="modal-dialog">
<div style="text-align: center">
  <video controls="" width="430" height="238" id="disconnect-video" muted="muted" loop="loop" autoplay>
   <source src="/static/videos/monkey.mp4" type='video/mp4; codecs="avc1.42E01E, mp4a.40.2"'>
  </video>
  </div>
  <div style="margin: 20px">
  <h2 style="text-align: center">{{ .LANG.Modal.wsdisconnect.title }}</h2>
  <h2 style="text-align: center">{{ .LANG.Modal.wsdisconnect.info }}</h2>
  <h3 style="text-align: center">{{ .LANG.Modal.wsdisconnect.trying }} <span id="disconnected-seconds"></span> {{ .LANG.SenseSeconds }}</h3>
  <div id="disconnected-success" style="margin: 10px; display: none; text-align: center; color: lime">
  {{ .LANG.Modal.wsdisconnect.connected }}
  </div>
  <div id="disconnected-failed" style="display: none; text-align: center; color: red; margin: 10px">
  {{ .LANG.Modal.wsdisconnect.failed }}
  </div>
  <div id="disconnected-process" style="margin: 10px; display: none; text-align: center; color: orange">
  {{ .LANG.Modal.wsdisconnect.connecting }}
  </div>
  </div>
  </div>
  </div>
 </div>
</div>