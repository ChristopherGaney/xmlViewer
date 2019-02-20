var viewer = (function() {

       return {
       
        displayXML: function(result, params) {
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
        },
        displayHTML: function(result) {

        },
        displayList: function(result, cb1, cb2) {
            var display = $('#display_tb');
            var items = '';

            items += '<div class="list_wrapper">';
                
               $.each(result, function(i,v) {
                    var len = v.Urls.length;
                    var count = 0;
                     console.log(v.Urls.length);
                    items += '<div class="list_tables"><div class="row name_row"><div class="six columns listname">' + v.Name + '</div><div class="six columns"></div></div>';

                    $.each(v.Urls, function(i,m) {
                        var cls = '';
                        if(len > 1) {
                            if(count % 2 === 0) { 
                                cls = ' even';
                            } else { 
                                cls = ' odd';
                            }
                        }
                       items += '<div class="url_wrapper_wrap' + cls + '"><div class="row url_row"><div class="two columns listid">' + m.ID + '</div><div class="ten columns"><a href="' + m.Url + '" target="_blank" class="listurl">' + m.Url + '</a></div></div>' + '<div class="row meta_row"><div class="eight columns"><span class="listtype">' + m.Type + '</span><span class="listmethod">' + m.Method + '</span></div><div class="four columns"><span class="prs"><a href="#" class="item_parse" name="item-parse" value="parse">Parse</a></span><span class="ers"><a href="#ex1" class="item_edit" data-modal>Edit</a></span></div></div></div>';

                        count++;

                    });
                   items += '</div>';
                });

                items += '</div>';
                display.html(items);
                cb1();
                cb2();
        }

    };

})();
module.exports = viewer;
