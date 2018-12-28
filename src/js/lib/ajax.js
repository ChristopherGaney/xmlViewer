var ajax = (function() {

   $(document).ready(function() {
        var params;
        
        var sendParseRequest = function(t, cb) {
             var dt = JSON.stringify(t);
            axios.post('/poster/dats', dt)
              .then(cb)
              .catch(function (error) {
                console.log(error);
              });
        };
        var makeRequest = function(u,t,m) {
             params = {
                "url": u,
                "type": t,
                "method": m
              };
            sendParseRequest(params, function (response) {
                var display = $('#display_tb');
                 var items = '';
                var res;
                console.log("here is response from cb:");
               
                res = response.data;
                items += 'url: ' + res.url + ' <br>type: ' + res.type + '<br>method: ' + res.method;
                 
                display.html(items);
            });
        };

        console.log('app js here\ndocument ready!');

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

    });
})();
module.exports = ajax;
