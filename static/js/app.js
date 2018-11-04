var app = (function() {

   $(document).ready(function() {
        console.log('app js here\ndocument ready!');
    
    axios.post('/poster/testing', {
        firstName: 'Fred',
        lastName: 'Flintstone'
      })
      .then(function (response) {
        console.log("here is response:");
        console.log(response);
        alert(response.data.data);
      })
      .catch(function (error) {
        console.log(error);
      });


    } );
    

})();
