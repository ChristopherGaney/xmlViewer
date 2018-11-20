var ajax = (function() {

   $(document).ready(function() {
        console.log('app js here\ndocument ready!');
        console.log('starting output');
    axios.post('/poster/rockOutYo', {
        firstName: 'Fred',
        lastName: 'Flintstone'
      })
      .then(function (response) {
        console.log("here is response:");
        console.log(response);
      })
      .catch(function (error) {
        console.log(error);
      });

    });
    
})();
module.exports = ajax;
