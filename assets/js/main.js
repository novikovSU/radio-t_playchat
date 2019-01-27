Date.prototype.getLocalTime = function() {
    h = (h = this.getHours()) < 10 ? ('0' + h) : h;
    m = (m = this.getMinutes()) < 10 ? ('0' + m) : m;
    s = (s = this.getSeconds()) < 10 ? ('0' + s) : s;
    return h + ":" + m + ":" + s;
}

$.routes.add('/{id:int}/', 'issueNumRoute', function() {
    $('#picker_select').val(this.id)
    $('#picker_select').trigger("change");
});

$(document).ready(function(){

    $('#totop').click(function(e){
        e.preventDefault();
        $("#messages").prop({ scrollTop: 0 });
    });

    $('#tobottom').click(function(e){
        e.preventDefault();
        $("#messages").prop({ scrollTop: $("#messages")[0].scrollHeight });
    });

    $("#cc_panel").click(function(e){
        e.stopPropagation();
        var time = $(e.target).parent().parent().attr('data-time')
        console.log(time);
        $("#jquery_jplayer_1").jPlayer('play',time)
    });

    $.jPlayer.timeFormat.showHour = true;

    $("#jquery_jplayer_1").jPlayer({
        solution: "html, flash",
        swfPath: "//cdnjs.cloudflare.com/ajax/libs/jplayer/2.9.2/jplayer/jquery.jplayer.swf",
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
                $("#cc_panel").prop({ scrollTop: $("#cc_panel").prop("scrollHeight") });
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
            $(".cc").css("display","block");
        }
    });

    $('#picker_select').change(function(){
        var num = $(this).val();

        $('.message').remove();
        $('.cc').remove();
        $('#search_panel').remove();
        $('#live_chat').remove();
        //$('#numbers').css("display","inline-table");
        $('#loaded').text('0');
        $('#total').text('0');
        $('#fork_me').css("display","block");

        if (num == "--") {
            return;
        }

        location.href = "/#/"+num

        if (num == "live") {
            $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: "http://stream.radio-t.com"});
            //$("#jquery_jplayer_1").jPlayer("setMedia",{mp3: "http://pf.volna.top/PilotBy48"});
            //$('#messages').append('<iframe src="http://chat.radio-t.com/" id="live_chat"><p>Your browser does not support iframes.</p></iframe>');
            //$('#messages').append('<div class="message" >Try to implement Giiter integration later.</div>');
            $('#messages').append('<iframe src="https://gitter.im/radio-t/chat/~embed"><p>Your browser does not support iframes.</p></iframe>');
            $('.message').css("display","block");
            $('#fork_me').show();
            $("#jquery_jplayer_1").jPlayer('play');

            return;
        }

        $.getJSON("data/"+num+"/desc.json", function (data) {
            //console.log(data)
            $('#start_time').text(data.start_time);
            $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: data.url});
            $('#total').text(data.chat_n);
        });

        $.getJSON("data/"+num+"/chat.json", function (data) {
            var messages_tmp_panel = $(document.createElement('div'));
            $.each( data.chat, function( key, message ) {
                messages_tmp_panel.append('<div class="message t' + message.datetime + '"><p><span class="date">' + (new Date(message.datetime*1000).getLocalTime()) + '</span><span class="nickname ' + message.author_type + '">' + message.author_name + '</span><span class="text">' + message.text + '</span></p></div>');
                $('#loaded').text(parseInt($('#loaded').text())+1);
            });
            $('#messages').append(messages_tmp_panel)

            var start_time = parseInt($('#start_time').text());
            for (var i = start_time - 10800; i < start_time; i++ ) {
                $('.t' + i).css("display","block");
            }
            $('#tobottom').click();
        });

       $.getJSON("data/"+num+"/cc.json", function (data) {
            console.log('CC load start')
            var start_time = parseInt($('#start_time').text());
            var cc_tmp_panel = $(document.createElement('div'))
            $.each( data.subs, function( key, sub ) {
                var time = Math.trunc(sub.stime)
                cc_tmp_panel.append('<div data-time='+time+' class="cc t' + (start_time+time) + '"><p><span class="time"><!--a href="#" class="cc_time" onClick="return testFunction(this)"-->' + (new Date(time*1000-18000000).getLocalTime()) + '</a></span><span class="text">' + sub.text + '</span></p></div>');
                //$('#loaded').text(parseInt($('#loaded').text())+1);
                //console.log();
            });
            $('#cc_panel').append(cc_tmp_panel)
            console.log('CC load end')
        });
    });

    $('#picker_select').trigger("change");
});