var ajax = (function() {


   $(document).ready(function() {
        var params;
        var getRequest = function(ext, cb) {
            //var dt = JSON.stringify(t);
            console.log(ext);
            axios.get(ext)
              .then(cb)
              .catch(function (error) {
                console.log(error);
              });
        };

        var sendRequest = function(ext, t, cb) {
            var dt = JSON.stringify(t);
            axios.post(ext, dt)
              .then(cb)
              .catch(function (error) {
                console.log(error);
              });
        };

       
        var displayXML = function(result) {
            var display = $('#display_tb');
            var items = '';
            
            if(params.method === "flat-xml" || params.method === "deep-xml") {
                items += '<table id="fancytable" class="display"><col width="35%"><col width="65%">' +
                        '<thead><tr><th>Title</th><th>Keywords</th></tr></thead><tbody>';

                $.each(result, function(i,v) {
                        items += '<tr><td><a href="' + v.Location + '" target="_blank">' + v.Title + '</a></td><td>' + v.Keyword + '</td></tr>';
                });

                items += '</tbody></table>';

                
            }
            else if(params.method === "raw-xml") {
                
               items += '<textarea style="width: 100%; min-height: 500px;">' + result[0] + '</textarea>';
            }
            display.html(items);
        }
        var displayHTML = function(result) {

        }
        var displayList = function(result) {
            var display = $('#display_tb');
            var items = '';

            //console.log("here");
            //console.log(result);
            //console.log('not here');

            items += '<div class="list_wrapper">';
                
               $.each(result, function(i,v) {
                    //console.log(v)
                    items += '<div class="row name_row"><div class="six columns">' + v.Name + '</div><div class="six columns"></div></div>';

                    $.each(v.Urls, function(i,m) {
                        //console.log(m)
                       items += '<div class="row url_row"><div class="two columns">' + m.ID + '</div><div class="ten columns"><a href="' + m.Url + '" target="_blank">' + m.Url + '</a></div></div>' +

                                '<div class="row meta_row"><div class="eight columns"><span>' + m.Type + '</span><span>' + m.Method + '</span></div><div class="four columns"><input type="button" id="item_parse" name="item-parse" value="parse" /><input type="button" id="item_edit" name="item-edit" value="edit" /></div></div>';

                    });
                   
                });

                items += '</div>';
                display.html(items);
        }
        var makeRequest = function(u,t,m) {
             params = {
                "url": u,
                "type": t,
                "method": m
              };
            sendRequest('/poster', params, function (response) {
                var res = response.data;
                console.log('returned from /poster');

                if(params.type === 'xml') {
                    displayXML(res);
                }
                else {
                    displayHTML(res);
                }
            });
        };
        var makeListRequest = function() {
           /* var params = {
                "list": "true"
              };*/
            var params = "?list=bigList";
            getRequest('/lister' + params, function (response) {
                var res = response.data;
                
               displayList(res.Items);
            
            });
        };
            // ,u,t,m
        var addItemRequest = function(params) {
           
            sendRequest('/items', params, function (response) {
                var display = $('#display_tb');
                var res = response.data;
                
                console.dir(res);
                
                display.text("Go says: status ok");
            
            });
        };
        var delItemRequest = function(params) {
            sendRequest('/items', params, function (response) {
                var display = $('#display_tb');
                var res = response.data;
                
                console.dir(res);
                
                display.text("Go says: status ok");
            
            });
        };
        var modifyItemRequest = function(r,n,u,t,m) {
            var params = {
                "req": r,
                "name": n,
                "url": u,
                "type": t,
                "method": m
              };
              
            sendRequest('/items', params, function (response) {
                var display = $('#display_tb');
                var res = response.data;
                
                console.dir(res);
                
                display.text("Go says: status ok");
            
            });
        };
        var handleRadios = (function() {
            var type = '';
             var isSet = false;
            $('#resource-method').submit(function(e) {
                var url,type, method, name = '';
                url = $('#inp_url').val();
                type = $('input[name=rsource_typ]:checked').val();
                e.preventDefault();
                e.stopImmediatePropagation();

                if(type === 'xml') {
                    name = 'xml_typ';
                }
                else {
                    name = 'http_typ';
                }
                method = $('input[name=' + name + ']:checked').val();
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
       var getList = (function() {
            $('#show-list').on('click', function() {
                makeListRequest();
            });
       })();
       $('#new-item-form').submit(function(e) {
                var name, url, type, method, req = '';
                name = $('#i_name').val();
                url = $('#i_url').val();
                url_name = $('#i_url_name').val();
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

                //console.log("name: " + name + "url: " + url + "type: " + type + " method: " + method + " req: " + req);
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
                    delItemRequest({"req": req,
                                    "name": name});
                }
                 else if(req === 'del-url' && url !== '') {
                    console.log(req + ' ' + url);
                    $.modal.close();
                    delItemRequest({"req": req,
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

            $('a[data-modal]').click(function(event) {
                var url,type, method, name = '';
                url = $('#inp_url').val();
                type = $('input[name=rsource_typ]:checked').val();
                if(type === 'xml') {
                    name = 'xml_typ';
                }
                else {
                    name = 'http_typ';
                }
                method = $('input[name=' + name + ']:checked').val();
                $('#i_url').val(url);
                $('#i_type').val(type);
                $('#i_method').val(method);

              $(this).modal();
              return false;
            });
       })();
       var editListItem = (function() {

            $('a[data-modal]').click(function(event) {
                var url,type, method, name = '';
                url = $('#inp_url').val();
                type = $('input[name=rsource_typ]:checked').val();
                if(type === 'xml') {
                    name = 'xml_typ';
                }
                else {
                    name = 'http_typ';
                }
                method = $('input[name=' + name + ']:checked').val();
                $('#i_url').val(url);
                $('#i_type').val(type);
                $('#i_method').val(method);

              $(this).modal();
              return false;
            });
       })();
    });
})();
module.exports = ajax;
