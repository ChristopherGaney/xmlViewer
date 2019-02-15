
var ajax = require('./ajax.js');
var viewer = require('./viewer.js');

var parser = (function() {
        var params;

        var makeRequest = function(u,t,m) {
             params = {
                "url": u,
                "type": t,
                "method": m
              };
            ajax.sendRequest('/poster', params, function (response) {
                var res = response.data;
                console.log('returned from /poster');

                if(params.type === 'xml') {
                    viewer.displayXML(res, params);
                }
                else {
                    viewer.displayHTML(res, params);
                }
            });
        };
        
        var makeListRequest = function() {
            var params = "?list=bigList";
            ajax.getRequest('/lister' + params, function (response) {
                var res = response.data;
                
               viewer.displayList(res.Items, editListItem.setButton, parseListItem.setButton);
            
            });
        };
        
        var addItemRequest = function(params) {
           
            ajax.sendRequest('/items', params, function (response) {
                var display = $('#display_tb');
                var res = response.data;
                console.dir(res.code);
                console.log(res.message)
                if(res.code == 500) {
                    display.text("message: " + res.message + " code: " + res.code);
                }
                else {
                    display.text("Go says: status ok");
                }
            });
        };
        

         // Handler for Sidebar Parser Inputs
        var handleRadios = (function() {
            var type = '';
             var isSet = false;
            $('#resource-method').submit(function(e) {
                var url,type, method, tname = '';
                url = $('#inp_url').val();
                type = $('input[name=rsource_typ]:checked').val();
                e.preventDefault();
                e.stopImmediatePropagation();

                if(type === 'xml') {
                    tname = 'xml_typ';
                }
                else {
                    tname = 'http_typ';
                }
                method = $('input[name=' + tname + ']:checked').val();
                if(url !== '' && type !== '' && method !== '') {
                    makeRequest(url, type, method);
                }
                else {
                    alert('Please select all fields.');
                    console.log('required fields empty!');
                }
            });
            $('input[name=rsource_typ]').on('click', function() {
                var id = '';
                type = $('input[name=rsource_typ]:checked').val();
                 if(type === 'xml') {
                    
                    $('#http-method-type').fadeOut('slow', function() {
                         $('#xml-method-type').fadeIn('slow');
                    });
                   
                }
                else {
                    
                    $('#xml-method-type').fadeOut('slow', function() {
                         $('#http-method-type').fadeIn('slow');
                    });
                    
                }
                
            });
            $('input[name=xml_typ], input[name=http_typ]').on('click', function() {
                var formName = $('#resource-method');
                var quitListening = function() {
                    $('body').off('keypress');
                    $('#filter_btn').off('click');
                    isSet = false;
                };
                var setter = function(event) {
                    if(event.keyCode == 13 || event.which == 13) {
                        formName.submit();
                        quitListening();
                    }
                    else {
                        quitListening();
                    }
                };
                if(!isSet) {
                    $('body').on('keypress', function(event) {
                        event.preventDefault();
                        setter(event);  
                    });

                isSet = true;
                }
            });
            
       })();

       // Get the List of Content Providers
       var getList = (function() {
            $('#show-list').on('click', function() {
                makeListRequest();
            });
       })();


      // Handler for Editor Modal Inputs
       $('#new-item-form').submit(function(e) {
                var name, url, type, method, req = '';
                name = $('#i_name').val();
                url = $('#i_url').val();
                type = $('#i_type').val();
                method = $('#i_method').val();
                req = $('input[name=protor]:checked').val();
                var obj = {"req": req,
                        "name": name,
                        "url": url, 
                        "type": type, 
                        "method": method};
                e.preventDefault();
                e.stopImmediatePropagation();

                console.log("name: " + name + "url: " + url + "type: " + type + " method: " + method + " req: " + req);
                // && url !== '' && type !== '' && method !== ''
                // , url, type, method
                if(req === 'add' && name !== '') {
                    $.modal.close();
                    if(url !== '') {
                        addItemRequest(obj);
                    }
                    else {
                        addItemRequest({"req": req,
                                        "name": name});
                    }
                }

                else if(req === 'del-cp' && name !== '') {
                    console.log(req + ' ' + name);
                    $.modal.close();
                    addItemRequest({"req": req,
                                    "name": name});
                }
                 else if(req === 'del-url' && url !== '') {
                    console.log(req + ' ' + url);
                    $.modal.close();
                    addItemRequest({"req": req,
                                    "url": url});
                }
                else if(req === 'modify' && name !== '' && url !== '') {
                    $.modal.close();
                    addItemRequest(obj);
                }
                else {
                    alert('Please select all necessary fields.');
                    console.log('required fields empty!');
                }
            });


       var editParserItem = (function() {
            $('.add-btn a[data-modal]').click(function(event) {
                var name, url, type, method, tname = '';
                name = $('#hidden_name').val();
                url = $('#inp_url').val();
                type = $('input[name=rsource_typ]:checked').val();
                if(type === 'xml') {
                    tname = 'xml_typ';
                }
                else {
                    tname = 'http_typ';
                }
                method = $('input[name=' + tname + ']:checked').val();
                $('#i_name').val(name);
                $('#i_url').val(url);
                $('#i_type').val(type);
                $('#i_method').val(method);

              $(this).modal();
              return false;
            });
       })();


       var editListItem = (function() {
       		var setEditButton = function() {
       			var wasClicked = function(event) {
       				console.log("the box");
	                var name, url,type, method = '';
	                //console.log('item edit clicked');
	                var box = $(this).closest('.list_tables');
	                
	                name = box.find('.listname').text();
	                url = box.find('.listurl').text();
	                type = box.find('.listtype').text();
	                method = box.find('.listmethod').text();
	                console.log(url);
	                console.log("not box");
	                $('#i_name').val(name);
	                $('#i_url').val(url);
	                $('#i_type').val(type);
	                $('#i_method').val(method);

	              $(this).modal();
	              return false;
	            };
	            
	            $('.ers a[data-modal]').off('click').on('click', wasClicked);
        };

        return {
        	setButton: function() {
        		setEditButton();
        	}

        };
       })();


       var parseListItem = (function() {
       		var setParseButton = function() {
       			var wasClicked = function(event) {
       				console.log("the box");
	                var name, url, type, method, ptype, pmethod = '';
	                //console.log('item edit clicked');
	                var box = $(this).closest('.list_tables');
	                id = box.find('.listid').text();
	                name = box.find('.listname').text();
	                url = box.find('.listurl').text();
	                type = box.find('.listtype').text();
	                method = box.find('.listmethod').text();
	                console.log(name + url +type + method);
	                console.log("not box");
	                $('#hidden_id').val(id);
	                $('#hidden_name').val(name);
	                $('#inp_url').val(url);
	                if(type !== "" && type === "xml") {
	                	console.log('is xml');
	                	ptype = "xml_typ";
	                	if(method !== '' && method === "deep-xml") {
	                			
	                		pmethod = "deep_xml";
	                	} else if(method !== '' && method === "flat-xml") {
	                		pmethod = "flat_xml";
	                	} else {
	                		pmethod = "raw_xml";
	                	}
	                }
	                if(type !== "" && type === "http") {
	                	console.log('is http');
	                	ptype = "http_typ";
	                	if(method !== '' && method === "deep-http") {

	                		pmethod = "deep_http";
	                	} else if(method !== '' && method === "flat-http") {
	                		console.log('is flat http');
	                		pmethod = "flat_http";
	                	} else {
	                		pmethod = "raw_http";
	                	}
	                }
	                if(type === 'xml') {
		                    $('#http-method-type').fadeOut('slow', function() {
		                         $('#xml-method-type').fadeIn('slow');
		                    });
		                   
		                }
		                else {
		                    
		                    $('#xml-method-type').fadeOut('slow', function() {
		                         $('#http-method-type').fadeIn('slow');
		                    });
		                    
		                }
	                $('#' + ptype).prop("checked", true);
	                $('#' + pmethod).prop("checked", true);

	              //$(this).modal();
	              //return false;
	            };
	            
	            $('.prs a.item_parse').off('click').on('click', wasClicked);
        };

        return {
        	setButton: function() {
        		setParseButton();
        	}
        };
       })();

    
})();
module.exports = parser;
