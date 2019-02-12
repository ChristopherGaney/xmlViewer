
var ajax = (function() {

    return {
        getRequest: function(ext, cb) {
            console.log(ext);
            axios.get(ext)
              .then(cb)
              .catch(function (error) {
                console.log(error);
              });
        },

        sendRequest: function(ext, t, cb) {
            var dt = JSON.stringify(t);
            axios.post(ext, dt)
              .then(cb)
              .catch(function (error) {
                console.log(error);
              });
        }

       };
})();
module.exports = ajax;
