var app = (function() {

   $(document).ready(function() {
        console.log('app js here\ndocument ready!');

    var sendParseRequest = function(r, cb) {
        var dt = JSON.stringify(r);
    
       axios.post('/poster/dats', dt)
          .then(function (response) {
            console.log("here is response:");
            console.log(response);
            cb(r);
          })
          .catch(function (error) {
            console.log(error);
          });
        });
    };
   
    sendParseRequest({
            "firstName": 'Fred',
            "lastName": 'Flintstone'
        }, function(ret) {
           var display = $('#display_tb');
            display.html(ret);
        });
    

})();
