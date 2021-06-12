<div id="addNodesModal" class="modal"><div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">{{LANG.modal.addNodes.title}}</h3> <a href="#close" title="{{LANG.modal.close}}" class="close" onclick="$('#addNodesModal').hide(200)">X</a></div><div class="modal-body">
<div style="text-align: center"><span class="switch active" id="switch-single" onclick="switchTab('single')">{{LANG.modal.addNodes.label_single}}</span><span class="switch" id="switch-multiple" onclick="switchTab('multiple')">{{LANG.modal.addNodes.label_multiple}}</span></div>
<div id="addNodesErr"></div>
<div id="addNodesSingle">
<div style="margin: 20px 0 0 0;"><p>{{LANG.modal.addNodes.label_nodeIP}}</p><input id="nodeIP" type="text" class="inputtext" value="" placeholder="{{LANG.modal.addNodes.label_nodeIP_placeholder}}"></div>
<div style="margin: 0px">
<p>{{LANG.modal.addNodes.label_nodeName}}</p><input id="nodeName" class="inputtext" type="text" value="" placeholder="{{LANG.modal.addNodes.label_nodeName_placeholder}}">
<span style="font-size: 12px; padding-left: 0px; margin-left: 0px;"><br>{{LANG.modal.addNodes.label_recommend}}</span>
</div>
</div>
<div id="addNodesMultiple" style="display: none">
<div style="margin: 20px 0 0 0;">
<span>{{LANG.modal.addNodes.label_multi_nodeIP}}:</span>
<span style="font-size: 12px; padding-left: 0px; margin-left: 0px;"><br>{{LANG.modal.addNodes.label_multi_prompt}}</span>
<textarea style="width: 100%; height: 150px;" id="nodeIPList"></textarea>
</div>
<div>
<p>{{LANG.modal.addNodes.label_multi_nodeName}}:</p><input id="nodeNameList" class="inputtext" type="text" value="" placeholder="{{LANG.modal.addNodes.label_nodeName_placeholder}}">
<span style="font-size: 12px; padding-left: 0px; margin-left: 0px;"><br>{{LANG.modal.addNodes.label_recommend}}</span>
</div>
</div>
<div style="text-align: center">
<input type="button" value="{{LANG.modal.addNodes.add_btn}}" onclick="addNode()" id="addNodeButton" class="sendbutton">
</div>
</div>
</div>
</div>
</div>

