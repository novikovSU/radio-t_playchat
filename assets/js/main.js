Date.prototype.getLocalTime = function() {
    h = (h = this.getHours()) < 10 ? ('0' + h) : h;
    m = (m = this.getMinutes()) < 10 ? ('0' + m) : m;
    s = (s = this.getSeconds()) < 10 ? ('0' + s) : s;
    return h + ":" + m + ":" + s;
}

$.urlParam = function (name) {
    return decodeURI( (RegExp(name + '=' + '(.+?)(&|$)').exec(location.search)||[,null])[1] );
}

$(document).ready(function(){
    var volume = $.urlParam('vol');

    $('#totop').click(function(e){
        e.preventDefault();
        $("#messages").prop({ scrollTop: 0 });
    });

    $('#tobottom').click(function(e){
        e.preventDefault();
        $("#messages").prop({ scrollTop: $("#messages")[0].scrollHeight });
    });

    $("#jquery_jplayer_1").jPlayer({
        solution: "html, flash",
        swfPath: "//cdnjs.cloudflare.com/ajax/libs/jplayer/2.3.0/Jplayer.swf",
        supplied: "mp3",
        preload: "metadata",
        wmode: "window",
        ready: function () {
            $(this).jPlayer("setMedia", {
                mp3: "http://91.213.196.99:8181/stream"
            });
        },
        timeupdate: function (event) {
            var last_offset = parseInt($('#last_offset').text());
            var cur_offset = event.jPlayer.status.currentTime;
            if (last_offset == cur_offset) { 
                //$('#last_offset').text(cur_offset.toString());
                return;
            }
            var start_time = parseInt($("#start_time").text());
            var elem = $('#messages');
            if (cur_offset > last_offset) {
              if (elem[0].scrollHeight - elem.scrollTop() == elem.outerHeight()) {
                for(var i=last_offset; i <= cur_offset; i++) {
                  $(".t"+(start_time+i).toString()).css("display","block");
                }
                $("#messages").prop({ scrollTop: $("#messages").prop("scrollHeight") });
              } else {
                for(var i=last_offset; i <= cur_offset; i++) {
                  $(".t"+(start_time+i).toString()).css("display","block");
                }
              }
            } else {
              for(var i=last_offset; i >= cur_offset; i--) {
                $(".t"+(start_time+i).toString()).css("display","none");
              }
            }
            $('#last_offset').text(cur_offset.toString());
        },
        ended: function () {
            $(".message").css("display","block");
        }
    });
    $.jPlayer.timeFormat.showHour = true;

    $('#picker_select').change(function(){
        var num = $(this).val();
        var chatDbLink = new Firebase('https://radio-t.firebaseio.com/' + num + '/chat');
        var descDbLink = new Firebase('https://radio-t.firebaseio.com/' + num + '/desc');
        $('.message').remove();
        $('#live_chat').remove();
        $('#numbers').css("display","inline-table");
        $('#loaded').text('0');
        $('#total').text('0');
        $('#fork_me').css("display","block");
        if (num == "--") {
            return;
        }
        if (num == "live") {
            $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: "http://91.213.196.99:8181/stream"});
            //$('#messages').append('<iframe src="http://chat.radio-t.com/" id="live_chat"><p>Your browser does not support iframes.</p></iframe>');
            $('#messages').append('<div id="candy"></div>');
            $('#fork_me').show();
            
            Candy.init('http://towee.net:5280/http-bind/', {
                core: { debug: true, autojoin: ['online@conference.radio-t.com'] },
                view: { language: 'en',
                    resources: './assets/candy/',
                    crop: { message: { nickname: 15, body: 1000 }, roster: { nickname: 15 } },
                    messages: {limit: 2000, remove: 500}
                }
            });

            // enable Colors plugin (default: 8 colors)
            CandyShop.Colors.init();

            Candy.Core.connect();

            return;
        }
        descDbLink.once('value', function(snapshot) {
            $('#start_time').text(snapshot.val().start_time);
            $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: snapshot.val().url});
            $('#total').text(snapshot.val().number_chat_messages);
        });
        chatDbLink.on('child_added', function(snapshot) {
            var message = snapshot.val();
            $('#messages').append('<div class="message t' + message.datetime + '"><p><span class="date">' + (new Date(message.datetime*1000).getLocalTime()) + '</span><span class="nickname ' + message.type + '">' + message.author + '</span><span class="text">' + message.text + '</span></p></div>');
            $('#loaded').text(parseInt($('#loaded').text())+1);
        });
        /*var MongoHQAPI = "https://api.mongohq.com/databases/radio-t/collections/meta/documents?_apikey=fi2miyanyecuw8wqxgsh";
        $.getJSON(MongoHQAPI, {'q':{'number':num}}, function(response){
                var meta = response[0];
                console.log(JSON.stringify(response));
                $('#start_time').text(meta.start_time);
                $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: meta.url});
                $('#total').text(meta.chat_messages);
        });
        $.getJSON('/data/' + num.toString() + '.json', function(data) {
            //console.log('loading data from: '+'/data/' + num.toString() + '.json');
            var messages=[];
            $.each(data, function(key, message) {
                messages.push('<div class="message t' + message.datetime + '"><p><span class="date">' + (new Date(message.datetime*1000).getLocalTime()) + '</span><span class="nickname ' + message.type + '">' + message.author + '</span><span class="text">' + message.text + '</span></p></div>');
            });
            $('#messages').append(messages.join(''));
            var start_time = parseInt($('#start_time').text());
            console.log('start_time is ' + start_time.toString());
            for (var i = start_time - 10800; i < start_time; i++ ) {
                $('.t' + i).css("display","block");
            }
            $('#tobottom').click();
        })*/
    });

    $('#loaded').bind('DOMSubtreeModified', function(e) {
        if ((e.target.innerHTML == $('#total').text()) && ($('#total').text() != '0')) {
            window.setTimeout(function(){$('#numbers').css("display","none");},15000);
            var start_time = parseInt($('#start_time').text());
            for (var i = start_time - 10800; i < start_time; i++ ) {
                $('.t' + i).css("display","block");
            }
            $('#tobottom').click();
        }
    });

    if(volume) {
        $("select#picker_select option").each(function() { this.selected = (this.text == volume); });
        $('#picker_select').trigger('change');
    }
});
