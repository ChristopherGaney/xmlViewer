var app = (function() {

   $(document).ready(function() {
        console.log('app js here\ndocument ready!');
    var dt = JSON.stringify({
        "firstName": 'Fred',
        "lastName": 'Flintstone'
      });
    axios.post('/poster/rockOutYo', dt)
      .then(function (response) {
        console.log("here is response:");
        console.log(response);
      })
      .catch(function (error) {
        console.log(error);
      });


    } );
    

})();
