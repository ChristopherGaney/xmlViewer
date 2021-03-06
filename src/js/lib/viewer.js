var ajax = require('./ajax.js');

var viewer = (function() {
    
    var getData = function() {
            var saveText = $('#save_text');
            return saveText.val();
        
        };
    var makeRequest = function(params) {
        ajax.sendRequest('/items', params, function (response) {
                    var disPlay = $('#display_tb');
                    var res = response.data;
                    console.dir(res.code);
                    console.log(res.message)
                    if(res.code == 500) {
                        disPlay.html("message: " + res.message + "<br>code: " + res.code);
                    }
                    else {

                        disPlay.text("Go says: status ok");
                    }
                });
    };

       return {
        
         setClear: function() {
                var el = $('#save_text');
                console.log('textarea cleared');
                el.val('');
             },
        setSave: function() {
                var stuff = getData();
                var url = $('#inp_url').val();
                var req = 'save-xml';
                console.log(stuff);
                var params = {"req": req,
                                "url": url,
                                "data": stuff};

                console.log('editing display');
                makeRequest(params);
             },
        setDelete: function() {
                var url = $('#inp_url').val();
                var req = 'del-xml-cache';
                console.log('deleting cache');
                var params = {"req": req,
                                "url": url};

                console.log('editing cache');
                makeRequest(params);
             },
        displayXML: function(result, params) {
            var display = $('#display_tb');
            var items = '';
            var check = 0;
            if(params.method === "deep-xml") {
                items += '<table id="fancytable" class="display"><col width="100%">' +
                        '<thead><tr><th>Location</th></tr></thead><tbody>';

                $.each(result, function(i,v) {
                        items += '<tr><td><a href="' + v.Location + '" target="_blank">' + v.Location + '</a></td></tr>';
                });

                items += '</tbody></table>';

                
            }
            else if(params.method === "raw-xml") {
                console.log(result);
                console.log(result[0]);
               items += '<textarea id="save_text" style="width: 100%; min-height: 500px;">' + result[0] + '</textarea>';
                check = 1;
            }
            else if(result[0].Keyword) {
                console.log('Keyword Cometh: ' + result[0].Keyword);

                     items += '<table id="fancytable" class="display"><col width="35%"><col width="65%">' +
                        '<thead><tr><th>Title</th><th>Keywords</th></tr></thead><tbody>';

                    $.each(result, function(i,v) {
                        var t,ht;
                        if(v.Title) {
                            ht = v.Title.replace("<!--// <![CDATA[", "");
                            ht = ht.replace("// ]]> -->", "");
                            t = ht;
                        }
                        else {
                            t = v.Location;
                        }
                        items += '<tr><td><a href="' + v.Location + '" target="_blank">' + t + '</a></td><td>' + v.Keyword + '</td></tr>';
                    });

                    items += '</tbody></table>';
            }
            else {
                items += '<table id="fancytable" class="display"><col width="25%"><col width="75%">' +
                        '<thead><tr><th>Publish date</th><th>Location</th></tr></thead><tbody>';
                
                $.each(result, function(i,v) {
                    var t;
                    if(v.Title) {
                            console.log('has tiltel');
                            t = v.Title;
                            t = t.replace("<![CDATA[", "").replace("]]>", "");
                        }
                     else {
                         t = "Location";
                    }
                    items += '<tr><td>' + t + '</td><td><a class="small_txt" href="' + v.Location + '" target="_blank">' + v.Location + '</a></td></tr>';
                });

                items += '</tbody></table>';
            }
           
            display.html(items);

            $('#fancytable').DataTable({
                "searching": true
            });
            if(check === 1) {
                $('#save_display').on('click', viewer.setSave);
                $('#clear_display').on('click', viewer.setClear);
            }
            else {
                $('#save_display').off('click', viewer.setSave);
                $('#clear_display').off('click', viewer.setClear);
            }
            $('#delete_display').off('click').on('click', viewer.setDelete);
        },

        displayHTML: function(result) {

        },
        displayList: function(result, cb1, cb2) {
            var display = $('#display_tb');
            var items = '';
            $('#save_display').off('click', viewer.setSave);
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
                       items += '<div class="url_wrapper_wrap' + cls + '"><div class="row url_row"><div class="two columns listid">' + m.ID + '</div><div class="ten columns"><a href="' + m.Url + '" target="_blank" class="listurl">' + m.Url + '</a></div></div>' + '<div class="row meta_row"><div class="eight columns"><span class="listtype">' + m.Type + '</span><span class="listmethod">' + m.Method + '</span></div><div class="four columns"><span class="prs"><a href="javascript:void(0)" class="item_parse" name="item-parse" value="parse">Parse</a></span><span class="ers"><a href="#ex1" class="item_edit" data-modal>Edit</a></span></div></div></div>';

                        count++;

                    });
                   items += '</div>';
                });

                items += '</div>';
                display.html(items);
                cb1();
                cb2();
                $('#delete_display').off('click').on('click', viewer.setDelete);
        }

    };

})();
module.exports = viewer;
