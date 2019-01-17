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

            items += '<table id="fancytable" class="display"><col width="25%"><col width="45%"><col width="15%"><col width="15%">' +
                        '<thead><tr><th>Name</th><th>Url</th><th>Type</th><th>Method</th></tr></thead><tbody>';

               $.each(result, function(i,v) {
                        items += '<tr><td>' + v.Name + '</td><td><a href="' + v.Url + '" target="_blank">' + v.Url + '</a></td><td>' + v.Type + '</td><td>' + v.Method + '</td></tr>';
                });

                items += '</tbody></table>';
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
                
               displayList(res.Outlets);
            
            });
        };
        var addItemRequest = function(n,u,t,m) {
            var params = {
                "name": n,
                "url": u,
                "type": t,
                "method": m
              };
              
            sendRequest('/adder', params, function (response) {
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
                var name, url, type, method = '';
                name = $('#i_name').val();
                url = $('#i_url').val();
                type = $('#i_type').val();
                method = $('#i_method').val();
                e.preventDefault();
                e.stopImmediatePropagation();
                console.log("name: " + name + "url: " + url + "type: " + type + " method: " + method);
                if(name !== '' && url !== '' && type !== '' && method !== '') {
                    $.modal.close();
                    addItemRequest(name, url, type, method);

                }
                else {
                    alert('Please select all fields.');
                    console.log('required fields empty!');
                }
            });
       var addItem = (function() {

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
