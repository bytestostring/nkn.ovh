<div id="addNodesModal" class="modal"><div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{.LANG.Modal.addNodes.title}}</h3> <a href="#close" title="{{.LANG.Modal.control.close}}" class="close" onclick="closeModal('addNodes')">X</a></div><div class="modal-body">
<div style="text-align: center"><span class="switch active" id="switch-nodes-single" onclick="switchTab('single')">{{.LANG.Modal.addNodes.label_single}}</span><span class="switch" id="switch-nodes-multiple" onclick="switchTab('multiple')">{{.LANG.Modal.addNodes.label_multiple}}</span></div>
<div id="addNodesErr"></div>
<div id="addNodesSingle">
<div style="margin: 20px 0 0 0;"><p>{{.LANG.Modal.addNodes.label_nodeIP}}</p><input id="nodeIP" type="text" class="inputtext" value="" placeholder="{{.LANG.Modal.addNodes.label_nodeIP_placeholder}}"></div>
<div style="margin: 0px">
<p>{{.LANG.Modal.addNodes.label_nodeName}}</p><input id="nodeName" class="inputtext" type="text" value="" placeholder="{{.LANG.Modal.addNodes.label_nodeName_placeholder}}">
<span style="font-size: 12px; padding-left: 0px; margin-left: 0px;"><br>{{.LANG.Modal.addNodes.label_recommend}}</span>
</div>
</div>
<div id="addNodesMultiple" style="display: none">
<div style="margin: 20px 0 0 0;">
<span>{{.LANG.Modal.addNodes.label_multi_nodeIP}}:</span>
<span style="font-size: 12px; padding-left: 0px; margin-left: 0px;"><br>{{.LANG.Modal.addNodes.label_multi_prompt}}</span>
<textarea style="width: 100%; height: 150px;" id="nodeIPList"></textarea>
</div>
<div>
<p>{{.LANG.Modal.addNodes.label_multi_nodeName}}:</p><input id="nodeNameList" class="inputtext" type="text" value="" placeholder="{{.LANG.Modal.addNodes.label_nodeName_placeholder}}">
<span style="font-size: 12px; padding-left: 0px; margin-left: 0px;"><br>{{.LANG.Modal.addNodes.label_recommend}}</span>
</div>
</div>
<div style="text-align: center">
<input type="button" value="{{.LANG.Modal.addNodes.add_btn}}" onclick="addNodes()" id="addNodeButton" class="sendbutton">
</div>
</div>
</div>
</div>
</div>

